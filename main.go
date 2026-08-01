package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// JWT Secret - in production keep this in an env var
var jwtSecret = []byte("your-secret-key-change-this-in-production")

// Auth structures
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	UserID   int    `json:"userId"`
	IsAdmin  bool   `json:"isAdmin"`
}

type Claims struct {
	UserID   int    `json:"userId"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ReviewRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Rating  int    `json:"rating"`
}

type Review struct {
	ID        int    `json:"id"`
	UserID    *int   `json:"userId,omitempty"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	Rating    int    `json:"rating"`
	CreatedAt string `json:"createdAt"`
}

// Structură pentru Leaderboard
type LeaderboardEntry struct {
	Username               string  `json:"username"`
	Accuracy               float64 `json:"accuracy"`
	TotalCorrectAnswers    int     `json:"totalCorrectAnswers"`
	TotalQuestionsAnswered int     `json:"totalQuestionsAnswered"`
	CurrentElo             float64 `json:"currentElo"`
}

// Helper function to send JSON error responses
func jsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

const defaultPostgresDSN = "postgres://dev:1234@localhost:5432/local_db?sslmode=disable"

func initDB() (*sql.DB, error) { // Initialize PostgreSQL database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = defaultPostgresDSN
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Verify connection early
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Initialize tables
	if err := initTables(db); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return db, nil
}

func initTables(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id BIGSERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		is_admin BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS reviews (
		id BIGSERIAL PRIMARY KEY,
		user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
		username TEXT,
		content TEXT NOT NULL,
		rating INTEGER CHECK (rating >= 1 AND rating <= 5),
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	CREATE INDEX IF NOT EXISTS idx_reviews_created_at ON reviews(created_at DESC);
	`

	if _, err := db.Exec(schema); err != nil {
		return err
	}

	if _, err := db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS is_admin BOOLEAN NOT NULL DEFAULT FALSE`); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_is_admin ON users(is_admin)`); err != nil {
		return err
	}

	if _, err := promoteAdminsByEmail(db); err != nil {
		return err
	}

	return nil
}

func parseAdminEmailsFromEnv() []string {
	raw := strings.TrimSpace(os.Getenv("ADMIN_EMAILS"))
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	seen := map[string]bool{}
	for _, p := range parts {
		email := strings.ToLower(strings.TrimSpace(p))
		if email == "" || seen[email] {
			continue
		}
		seen[email] = true
		out = append(out, email)
	}
	return out
}

func shouldEmailBeAdmin(email string) bool {
	emails := parseAdminEmailsFromEnv()
	target := strings.ToLower(strings.TrimSpace(email))
	if target == "" {
		return false
	}
	for _, e := range emails {
		if e == target {
			return true
		}
	}
	return false
}

func promoteAdminsByEmail(db *sql.DB) (int64, error) {
	emails := parseAdminEmailsFromEnv()
	if len(emails) == 0 {
		return 0, nil
	}

	var total int64
	for _, email := range emails {
		res, err := db.Exec(`UPDATE users SET is_admin = TRUE WHERE LOWER(email) = $1`, email)
		if err != nil {
			return total, err
		}
		affected, _ := res.RowsAffected()
		total += affected
	}

	return total, nil
}

func isAdminUser(userID int) (bool, error) {
	var isAdmin bool
	err := db.QueryRow(`SELECT is_admin FROM users WHERE id = $1`, userID).Scan(&isAdmin)
	if err != nil {
		return false, err
	}
	return isAdmin, nil
}

func insertTestUserIfNotExists(db *sql.DB, username, email, password string) (int, error) {
	// Hash password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return 0, err
	}

	var userID int
	err = db.QueryRow(`
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO NOTHING
		RETURNING id
	`, username, email, hashedPassword).Scan(&userID)

	// If user already exists, we select it
	if err == sql.ErrNoRows {
		err = db.QueryRow(`
			SELECT id FROM users WHERE email = $1
		`, email).Scan(&userID)
	}

	return userID, err
}

func apiLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	// Leaderboard GrInfo ordonat după acuratețe, apoi după nr. de răspunsuri corecte.
	rows, err := db.Query(`
		SELECT
			u.username,
			COALESCE(stats.total_correct_answers, 0) AS total_correct_answers,
			COALESCE(stats.total_questions_answered, 0) AS total_questions_answered,
			COALESCE(latest.current_elo, 1000) AS current_elo
		FROM users u
		LEFT JOIN LATERAL (
			SELECT
				SUM(correct_answers) AS total_correct_answers,
				SUM(total_questions) AS total_questions_answered
			FROM grinfo_sessions gs
			WHERE gs.user_id = u.id
		) stats ON TRUE
		LEFT JOIN LATERAL (
			SELECT COALESCE(final_elo, initial_elo) AS current_elo
			FROM grinfo_sessions gs
			WHERE gs.user_id = u.id
			ORDER BY gs.id DESC
			LIMIT 1
		) latest ON TRUE
		ORDER BY
			CASE
				WHEN COALESCE(stats.total_questions_answered, 0) > 0
				THEN (COALESCE(stats.total_correct_answers, 0)::float / COALESCE(stats.total_questions_answered, 0)::float)
				ELSE 0
			END DESC,
			COALESCE(stats.total_correct_answers, 0) DESC,
			u.username ASC
		LIMIT 100
	`)
	if err != nil {
		jsonError(w, "Eroare la preluarea clasamentului", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	leaderboard := make([]LeaderboardEntry, 0)
	for rows.Next() {
		var entry LeaderboardEntry
		if err := rows.Scan(&entry.Username, &entry.TotalCorrectAnswers, &entry.TotalQuestionsAnswered, &entry.CurrentElo); err != nil {
			continue
		}
		if entry.TotalQuestionsAnswered > 0 {
			entry.Accuracy = (float64(entry.TotalCorrectAnswers) / float64(entry.TotalQuestionsAnswered)) * 100
		}
		leaderboard = append(leaderboard, entry)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaderboard)
}

// Hash password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Check password hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate JWT token
func generateToken(userID int, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		IsAdmin:  false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Validate JWT token
func validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// Middleware to protect routes
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := validateToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", strconv.Itoa(claims.UserID))
		r.Header.Set("X-Username", claims.Username)
		r.Header.Set("X-Is-Admin", strconv.FormatBool(claims.IsAdmin))

		next(w, r)
	}
}

func adminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		uidStr := r.Header.Get("X-User-ID")
		uid, err := strconv.Atoi(uidStr)
		if err != nil || uid <= 0 {
			jsonError(w, "Autentificare invalidă", http.StatusUnauthorized)
			return
		}

		isAdmin, err := isAdminUser(uid)
		if err == sql.ErrNoRows {
			jsonError(w, "Utilizator inexistent", http.StatusUnauthorized)
			return
		}
		if err != nil {
			jsonError(w, "Eroare la verificarea permisiunilor", http.StatusInternalServerError)
			return
		}
		if !isAdmin {
			jsonError(w, "Acces interzis: contul nu are rol de admin", http.StatusForbidden)
			return
		}

		next(w, r)
	})
}

// API models and handlers
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func apiUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		id = strings.TrimPrefix(r.URL.Path, "/api/user/")
	}
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	row := db.QueryRow("SELECT id, username FROM users WHERE id = $1", userID)
	var u User
	if err := row.Scan(&u.ID, &u.Username); err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(u)
}

// POST /api/register - Register new user
func apiRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		jsonError(w, "Username, email and password are required", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		jsonError(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		jsonError(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	isAdmin := shouldEmailBeAdmin(req.Email)

	var userID int
	err = db.QueryRow(
		"INSERT INTO users (username, email, password_hash, is_admin) VALUES ($1, $2, $3, $4) RETURNING id",
		req.Username, req.Email, hashedPassword, isAdmin,
	).Scan(&userID)
	if err != nil {
		// Postgres unique violation details vary; keep generic message
		jsonError(w, "Username or email already exists", http.StatusConflict)
		return
	}

	token, err := generateToken(userID, req.Username)
	if err != nil {
		jsonError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token:    token,
		Username: req.Username,
		UserID:   userID,
		IsAdmin:  isAdmin,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

// POST /api/login - Login user
func apiLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var userID int
	var username, passwordHash string
	var isAdmin bool
	err := db.QueryRow(
		"SELECT id, username, password_hash, is_admin FROM users WHERE email = $1",
		req.Email,
	).Scan(&userID, &username, &passwordHash, &isAdmin)

	if err != nil {
		if err == sql.ErrNoRows {
			jsonError(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !checkPasswordHash(req.Password, passwordHash) {
		jsonError(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := generateToken(userID, username)
	if err != nil {
		jsonError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token:    token,
		Username: username,
		UserID:   userID,
		IsAdmin:  isAdmin,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// GET /api/profile - Get current user profile (protected route)
func apiProfileHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, _ := strconv.Atoi(userIDStr)
	username := r.Header.Get("X-Username")
	isAdmin, err := isAdminUser(userID)
	if err != nil {
		isAdmin = false
	}

	response := map[string]interface{}{
		"userId":   userID,
		"username": username,
		"isAdmin":  isAdmin,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// POST /api/reviews - Create a new review
func apiCreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		jsonError(w, "Content is required", http.StatusBadRequest)
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		jsonError(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	// Prefer the name entered in the form; authenticated users still get linked by user_id,
	// but the visible review name comes from the form.
	var userID *int
	var username string = "Anonim"

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			if claims, err := validateToken(parts[1]); err == nil {
				userID = &claims.UserID
			}
		}
	}

	providedName := strings.TrimSpace(req.Name)
	if providedName != "" {
		username = providedName
	}

	var reviewID int
	err := db.QueryRow(`
		INSERT INTO reviews (user_id, username, content, rating)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, userID, username, req.Content, req.Rating).Scan(&reviewID)

	if err != nil {
		fmt.Printf("Failed to insert review: %v\n", err)
		jsonError(w, "Failed to create review", http.StatusInternalServerError)
		return
	}

	var createdAt time.Time
	err = db.QueryRow(`
		SELECT created_at FROM reviews WHERE id = $1
	`, reviewID).Scan(&createdAt)

	if err != nil {
		jsonError(w, "Failed to fetch review", http.StatusInternalServerError)
		return
	}

	review := Review{
		ID:        reviewID,
		UserID:    userID,
		Username:  username,
		Content:   req.Content,
		Rating:    req.Rating,
		CreatedAt: createdAt.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(review)
}

// GET /api/reviews - Get all reviews
func apiGetReviewsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query(`
		SELECT id, user_id, username, content, rating, created_at
		FROM reviews
		ORDER BY created_at DESC
	`)
	if err != nil {
		fmt.Printf("Failed to query reviews: %v\n", err)
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		var userID sql.NullInt64
		var createdAt time.Time

		err := rows.Scan(&review.ID, &userID, &review.Username, &review.Content, &review.Rating, &createdAt)
		if err != nil {
			fmt.Printf("Error scanning review: %v\n", err)
			continue
		}

		if userID.Valid {
			uid := int(userID.Int64)
			review.UserID = &uid
		}
		review.CreatedAt = createdAt.Format(time.RFC3339)
		reviews = append(reviews, review)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(reviews)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ Could not load .env file, using system environment variables")
	}

	var dbErr error
	db, dbErr = initDB()
	if dbErr != nil {
		panic(dbErr)
	}

	if err := initGrInfoTablesAndSeed(db); err != nil {
		panic(err)
	}

	insertTestUserIfNotExists(db, "testuser", "test@test.com", "password123")

	mux := http.NewServeMux()

	mux.HandleFunc("/api/register", apiRegisterHandler)
	mux.HandleFunc("/api/login", apiLoginHandler)
	mux.HandleFunc("/api/profile", authMiddleware(apiProfileHandler))

	mux.HandleFunc("/api/user", apiUserHandler)
	mux.HandleFunc("/api/leaderboard", apiLeaderboardHandler)
	mux.HandleFunc("/api/grinfo/questions", apiGrInfoQuestionsHandler)
	mux.HandleFunc("/api/grinfo/admin/questions", adminMiddleware(apiGrInfoAdminQuestionsHandler))
	mux.HandleFunc("/api/grinfo/categories", apiGrInfoCategoriesHandler)
	mux.HandleFunc("/api/grinfo/quiz-start", authMiddleware(apiGrInfoSessionStartHandler))
	mux.HandleFunc("/api/grinfo/quiz-finish", authMiddleware(apiGrInfoSessionFinishHandler))
	mux.HandleFunc("/api/grinfo/incident", apiGrInfoIncidentHandler)
	mux.HandleFunc("/api/grinfo/session", authMiddleware(apiGrInfoSessionHandler))
	mux.HandleFunc("/api/grinfo/profile", authMiddleware(apiGrInfoProfileHandler))
	mux.HandleFunc("/api/grinfo/ai-review", authMiddleware(apiGrInfoAIReviewHandler))
	mux.HandleFunc("/api/reviews", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			apiCreateReviewHandler(w, r)
		} else if r.Method == http.MethodGet {
			apiGetReviewsHandler(w, r)
		} else {
			jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	distPath := "./frontend/dist"
	if os.Getenv("SERVE_FRONTEND") == "1" {
		if fi, err := os.Stat(distPath); err == nil && fi.IsDir() {
			fs := http.FileServer(http.Dir(distPath))
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				if strings.HasPrefix(r.URL.Path, "/api/") {
					mux.ServeHTTP(w, r)
					return
				}
				path := distPath + r.URL.Path
				if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
					fs.ServeHTTP(w, r)
					return
				}
				http.ServeFile(w, r, distPath+"/index.html")
			})
		}
	} else {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = w.Write([]byte("GRInfo API is running. Use /api/* endpoints or open the frontend on http://localhost:5173."))
		})
	}

	handler := withCORS(mux)
	fmt.Println("Listening on :8000")
	_ = http.ListenAndServe(":8000", handler)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
