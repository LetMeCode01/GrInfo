package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/smtp"
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
}

type Claims struct {
	UserID   int    `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type ProgressRequest struct {
	CourseID       string `json:"courseId"`
	QuizScore      int    `json:"quizScore"`
	TotalQuestions int    `json:"totalQuestions"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SubscribeRequest struct {
	Email string `json:"email"`
}

type SubscribeResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
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
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS user_progress (
		id BIGSERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		course_id TEXT NOT NULL,
		quiz_score INTEGER NOT NULL DEFAULT 0,
		total_questions INTEGER NOT NULL DEFAULT 0,
		completed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS courses (
		id BIGSERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		type TEXT NOT NULL,
		description TEXT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS user_stats (
		user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
		total_xp INTEGER NOT NULL DEFAULT 0,
		current_streak INTEGER NOT NULL DEFAULT 0,
		longest_streak INTEGER NOT NULL DEFAULT 0,
		last_activity_date DATE,
		lessons_completed INTEGER NOT NULL DEFAULT 0,
		quizzes_completed INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS lessons (
		id TEXT PRIMARY KEY,
		course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		"order" INTEGER,
		content TEXT
	);

	CREATE TABLE IF NOT EXISTS user_lessons (
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		lesson_id TEXT NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
		completed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		PRIMARY KEY (user_id, lesson_id)
	);

	CREATE TABLE IF NOT EXISTS newsletter_subscriptions (
		id BIGSERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		subscribed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
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
	CREATE INDEX IF NOT EXISTS idx_user_progress_user_id ON user_progress(user_id);
	CREATE INDEX IF NOT EXISTS idx_courses_type ON courses(type);
	CREATE INDEX IF NOT EXISTS idx_lessons_course_id ON lessons(course_id);
	CREATE INDEX IF NOT EXISTS idx_newsletter_email ON newsletter_subscriptions(email);
	CREATE INDEX IF NOT EXISTS idx_reviews_created_at ON reviews(created_at DESC);
	`

	_, err := db.Exec(schema)
	return err
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

		next(w, r)
	}
}

func listCoursesByType(w http.ResponseWriter, courseType string) { // List courses of a specific type
	rows, err := db.Query("SELECT id, title FROM courses WHERE type = $1", courseType)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	defer rows.Close()

	fmt.Fprintf(w, "(%s) Cursuri:<br>", courseType)

	for rows.Next() {
		var id int
		var title string
		_ = rows.Scan(&id, &title)
		fmt.Fprintf(w, "ID: %d — %s<br>", id, title)
	}
}

// API models and handlers
type Course struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Lesson struct {
	ID       string `json:"id"`
	CourseID string `json:"courseId"`
	Title    string `json:"title"`
	Order    int    `json:"order"`
	Content  string `json:"content"`
}

type CompletedLesson struct {
	ID          string    `json:"id"`
	CourseID    string    `json:"courseId"`
	Title       string    `json:"title"`
	CompletedAt time.Time `json:"completedAt"`
}

type CompletedLessonsByCourse struct {
	CourseID int               `json:"courseId"`
	Title    string            `json:"title"`
	Lessons  []CompletedLesson `json:"lessons"`
}

func apiCoursesHandler(w http.ResponseWriter, r *http.Request) {
	courseType := r.URL.Query().Get("type")
	rows, err := db.Query("SELECT id, title, type FROM courses WHERE type = $1", courseType)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Title, &c.Type); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		courses = append(courses, c)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(courses)
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

// Funcția pentru trimitere email
func sendNewsletterEmail(toEmail string, subject string, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SENDER_EMAIL")
	senderPassword := os.Getenv("SENDER_PASSWORD")

	fmt.Printf("\nAttempting to send email...\n")
	fmt.Printf("  To: %s\n", toEmail)
	fmt.Printf("  From: %s\n", senderEmail)
	fmt.Printf("  Host: %s:%s\n", smtpHost, smtpPort)

	if smtpHost == "" || senderEmail == "" || senderPassword == "" {
		fmt.Println("Email configuration incomplete:")
		fmt.Printf("   - SMTP_HOST: %v\n", smtpHost != "")
		fmt.Printf("   - SENDER_EMAIL: %v\n", senderEmail != "")
		fmt.Printf("   - SENDER_PASSWORD: %v\n", senderPassword != "")
		return fmt.Errorf("email not configured")
	}

	port := 587
	if smtpPort != "" {
		portNum, err := strconv.Atoi(smtpPort)
		if err == nil {
			port = portNum
		}
	}

	addr := net.JoinHostPort(smtpHost, fmt.Sprintf("%d", port))
	fmt.Printf("Connecting to %s...\n", addr)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("TCP connection failed: %v\n", err)
		return err
	}
	fmt.Println("TCP connection established")

	// Create SMTP client
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		fmt.Printf("SMTP client creation failed: %v\n", err)
		conn.Close()
		return err
	}
	defer client.Close()

	// Start TLS
	fmt.Println("Starting TLS...")
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpHost,
	}
	if err := client.StartTLS(tlsconfig); err != nil {
		fmt.Printf("StartTLS failed: %v\n", err)
		return err
	}
	fmt.Println("TLS started")

	// Authenticate
	fmt.Println("Authenticating...")
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)
	if err := client.Auth(auth); err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		return err
	}
	fmt.Println("Authentication successful")

	// Set sender
	if err := client.Mail(senderEmail); err != nil {
		fmt.Printf("Failed to set sender: %v\n", err)
		return err
	}

	// Set recipient
	if err := client.Rcpt(toEmail); err != nil {
		fmt.Printf("Failed to set recipient: %v\n", err)
		return err
	}

	// Send message
	wc, err := client.Data()
	if err != nil {
		fmt.Printf("Failed to get write closer: %v\n", err)
		return err
	}

	// Build email message
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		senderEmail, toEmail, subject, body)

	if _, err := wc.Write([]byte(message)); err != nil {
		fmt.Printf("Failed to write message: %v\n", err)
		return err
	}

	if err := wc.Close(); err != nil {
		fmt.Printf("Failed to close writer: %v\n", err)
		return err
	}

	client.Quit()

	fmt.Printf("✅ Email sent successfully to %s\n\n", toEmail)
	return nil
}

// POST /api/subscribe - Subscribe to newsletter
func apiSubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		jsonError(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Basic email validation
	if !strings.Contains(req.Email, "@") {
		jsonError(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		"INSERT INTO newsletter_subscriptions (email) VALUES ($1)",
		req.Email,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			jsonError(w, "Email already subscribed", http.StatusConflict)
			return
		}
		jsonError(w, "Failed to subscribe", http.StatusInternalServerError)
		return
	}

	// Trimite email de bun venit
	welcomeHTML := `
	<h2>Bine ai venit pe GrInfo! 🎓</h2>
	<p>Mulțumim pentru abonare! Vei primi actualizări cu quiz-uri si teste noi</p>
	`
	sendNewsletterEmail(req.Email, "Bine ai venit pe GrInfo!", welcomeHTML)

	response := SubscribeResponse{
		Message: "Te-ai abonat cu succes!",
		Email:   req.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
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

	var userID int
	err = db.QueryRow(
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
		req.Username, req.Email, hashedPassword,
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
	err := db.QueryRow(
		"SELECT id, username, password_hash FROM users WHERE email = $1",
		req.Email,
	).Scan(&userID, &username, &passwordHash)

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
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// GET /api/profile - Get current user profile (protected route)
func apiProfileHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, _ := strconv.Atoi(userIDStr)
	username := r.Header.Get("X-Username")

	rows, err := db.Query(`
		SELECT course_id, quiz_score, total_questions, completed_at 
		FROM user_progress 
		WHERE user_id = $1 
		ORDER BY completed_at DESC
	`, userID)
	if err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Progress struct {
		CourseID       string `json:"courseId"`
		QuizScore      int    `json:"quizScore"`
		TotalQuestions int    `json:"totalQuestions"`
		CompletedAt    string `json:"completedAt"`
	}

	var progress []Progress
	for rows.Next() {
		var p Progress
		if err := rows.Scan(&p.CourseID, &p.QuizScore, &p.TotalQuestions, &p.CompletedAt); err != nil {
			continue
		}
		progress = append(progress, p)
	}

	response := map[string]interface{}{
		"userId":   userID,
		"username": username,
		"progress": progress,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// POST /api/progress - Save user quiz progress (protected route)
func apiSaveProgressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	userID, _ := strconv.Atoi(userIDStr)

	var req ProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received progress: UserID=%d, CourseID=%s, Score=%d/%d\n",
		userID, req.CourseID, req.QuizScore, req.TotalQuestions)

	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("Database transaction error: %v\n", err)
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO user_progress (user_id, course_id, quiz_score, total_questions)
		VALUES ($1, $2, $3, $4)
	`, userID, req.CourseID, req.QuizScore, req.TotalQuestions)
	if err != nil {
		fmt.Printf("Failed to insert progress: %v\n", err)
		jsonError(w, "Failed to save progress", http.StatusInternalServerError)
		return
	}

	xpEarned := req.QuizScore * 10

	err = updateUserStats(tx, userID, xpEarned)
	if err != nil {
		fmt.Printf("Failed to update stats: %v\n", err)
		jsonError(w, "Failed to update stats", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		fmt.Printf("Failed to commit transaction: %v\n", err)
		jsonError(w, "Failed to save progress", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Progress saved successfully! User %d earned %d XP\n", userID, xpEarned)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"xpEarned": xpEarned,
	})
}

// updateUserStats updates the user's XP and streak
func updateUserStats(tx *sql.Tx, userID int, xpEarned int) error {
	var totalXP, currentStreak, longestStreak, quizzesCompleted int
	var lastActivityDate sql.NullString

	err := tx.QueryRow(`
		SELECT total_xp, current_streak, longest_streak, last_activity_date, quizzes_completed
		FROM user_stats WHERE user_id = $1
	`, userID).Scan(&totalXP, &currentStreak, &longestStreak, &lastActivityDate, &quizzesCompleted)

	today := time.Now().Format("2006-01-02")
	newStreak := currentStreak

	if err == sql.ErrNoRows {
		_, err = tx.Exec(`
			INSERT INTO user_stats (user_id, total_xp, current_streak, longest_streak, last_activity_date, quizzes_completed)
			VALUES ($1, $2, 1, 1, $3, 1)
		`, userID, xpEarned, today)
		return err
	} else if err != nil {
		return err
	}

	if lastActivityDate.Valid {
		lastDate, _ := time.Parse("2006-01-02", lastActivityDate.String)
		todayDate, _ := time.Parse("2006-01-02", today)
		daysDiff := int(todayDate.Sub(lastDate).Hours() / 24)

		if daysDiff == 0 {
			newStreak = currentStreak
		} else if daysDiff == 1 {
			newStreak = currentStreak + 1
		} else {
			newStreak = 1
		}
	} else {
		newStreak = 1
	}

	newLongestStreak := longestStreak
	if newStreak > longestStreak {
		newLongestStreak = newStreak
	}

	_, err = tx.Exec(`
		UPDATE user_stats 
		SET total_xp = $1, 
		    current_streak = $2, 
		    longest_streak = $3,
		    last_activity_date = $4,
		    quizzes_completed = $5,
		    updated_at = NOW()
		WHERE user_id = $6
	`, totalXP+xpEarned, newStreak, newLongestStreak, today, quizzesCompleted+1, userID)

	return err
}

// GET /api/stats - Get user statistics (protected route)
func apiStatsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, _ := strconv.Atoi(userIDStr)

	_, err := db.Exec(`
		INSERT INTO user_stats (user_id, total_xp, current_streak, longest_streak, quizzes_completed)
		VALUES ($1, 0, 0, 0, 0)
		ON CONFLICT (user_id) DO NOTHING
	`, userID)
	if err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	var stats struct {
		TotalXP          int    `json:"totalXp"`
		CurrentStreak    int    `json:"currentStreak"`
		LongestStreak    int    `json:"longestStreak"`
		LastActivityDate string `json:"lastActivityDate"`
		QuizzesCompleted int    `json:"quizzesCompleted"`
		LessonsCompleted int    `json:"lessonsCompleted"`
	}

	var lastActivityDate sql.NullString
	err = db.QueryRow(`
		SELECT total_xp, current_streak, longest_streak, last_activity_date, quizzes_completed, lessons_completed
		FROM user_stats WHERE user_id = $1
	`, userID).Scan(&stats.TotalXP, &stats.CurrentStreak, &stats.LongestStreak, &lastActivityDate, &stats.QuizzesCompleted, &stats.LessonsCompleted)
	if err != nil {
		jsonError(w, "Failed to fetch stats", http.StatusInternalServerError)
		return
	}
	if lastActivityDate.Valid {
		stats.LastActivityDate = lastActivityDate.String
	}

	weeklyXP := []int{0, 0, 0, 0, 0, 0, 0}
	rows, err := db.Query(`
		SELECT DATE(completed_at) as date, SUM(quiz_score * 10) as xp
		FROM user_progress
		WHERE user_id = $1 AND completed_at >= (CURRENT_DATE - INTERVAL '7 days')
		GROUP BY DATE(completed_at)
		ORDER BY date
	`, userID)
	if err == nil {
		defer rows.Close()
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		// Calculate Monday of current week
		// time.Weekday: Sunday=0, Monday=1, ..., Saturday=6
		weekday := int(today.Weekday())
		if weekday == 0 { // Sunday
			weekday = 7 // Treat Sunday as day 7
		}
		daysFromMonday := weekday - 1 // Monday = 0, Tuesday = 1, ..., Sunday = 6
		startOfWeek := today.AddDate(0, 0, -daysFromMonday)

		for rows.Next() {
			var dateStr string
			var xp int
			if err := rows.Scan(&dateStr, &xp); err == nil {
				// Parse date string (format: 2026-01-29 or 2026-01-29T00:00:00Z)
				completedDate, parseErr := time.Parse("2006-01-02", dateStr[:10])
				if parseErr != nil {
					fmt.Printf("Error parsing date %s: %v\n", dateStr, parseErr)
					continue
				}

				// Calculate days difference from Monday
				daysDiff := int(completedDate.Sub(startOfWeek).Hours() / 24)
				fmt.Printf("Date: %s, XP: %d, Weekday: %d, StartOfWeek (Mon): %s, DaysDiff: %d\n",
					dateStr, xp, weekday, startOfWeek.Format("2006-01-02"), daysDiff)

				if daysDiff >= 0 && daysDiff < 7 {
					weeklyXP[daysDiff] = xp
					fmt.Printf("Mapped to weeklyXP[%d] = %d (Day: %s)\n", daysDiff, xp,
						[]string{"Lun", "Mar", "Mie", "Joi", "Vin", "Sâm", "Dum"}[daysDiff])
				}
			}
		}
		fmt.Printf("Final weeklyXP: %v\n", weeklyXP)
	} else {
		fmt.Printf("Error querying weekly XP: %v\n", err)
	}

	response := map[string]interface{}{
		"stats":    stats,
		"weeklyXP": weeklyXP,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// apiCourseHandler returns a single course by id, e.g. /api/course?id=123
func apiCourseHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		idStr = strings.TrimPrefix(r.URL.Path, "/api/course/")
	}
	if idStr == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	row := db.QueryRow("SELECT id, title, type FROM courses WHERE id = $1", id)
	var c Course
	if err := row.Scan(&c.ID, &c.Title, &c.Type); err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(c)
}

// apiLessonsHandler returns a list of lessons for a course (3-5 items)
func apiLessonsHandler(w http.ResponseWriter, r *http.Request) {
	courseID := r.URL.Query().Get("courseId")
	if courseID == "" {
		courseID = strings.TrimPrefix(r.URL.Path, "/api/lessons/")
	}
	if courseID == "" {
		http.Error(w, "Missing courseId", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
        SELECT id, course_id, title
        FROM lessons
        WHERE course_id = $1
        ORDER BY id
    `, courseID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	lessons := make([]Lesson, 0)
	for rows.Next() {
		var l Lesson
		if err := rows.Scan(&l.ID, &l.CourseID, &l.Title); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		lessons = append(lessons, l)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(lessons)
}

// apiLessonHandler returns a single lesson by id and courseId
func apiLessonHandler(w http.ResponseWriter, r *http.Request) {
	courseID := r.URL.Query().Get("courseId")
	lessonID := r.URL.Query().Get("lessonId")

	if courseID == "" || lessonID == "" {
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/lesson/"), "/")
		if len(parts) >= 2 {
			courseID = parts[0]
			lessonID = parts[1]
		}
	}

	if courseID == "" || lessonID == "" {
		http.Error(w, "Missing courseId or lessonId", http.StatusBadRequest)
		return
	}

	row := db.QueryRow(`
        SELECT id, course_id, title
        FROM lessons
        WHERE course_id = $1 AND id = $2
    `, courseID, lessonID)

	var l Lesson
	if err := row.Scan(&l.ID, &l.CourseID, &l.Title); err != nil {
		http.Error(w, "Lesson not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(l)
}

func apiUserLessonsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userId", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
		SELECT
			c.id AS course_id,
			c.title AS course_title,
			l.id AS lesson_id,
			l.title AS lesson_title,
			l.course_id AS lesson_course_id,
			ul.completed_at
		FROM user_lessons ul
		JOIN lessons l ON l.id = ul.lesson_id
		JOIN courses c ON c.id = l.course_id
		WHERE ul.user_id = $1
		ORDER BY c.id, l.id
	`, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	result := []CompletedLessonsByCourse{}
	currentCourseID := -1
	currentCourseTitle := ""
	currentLessons := []CompletedLesson{}

	for rows.Next() {
		var courseID int
		var courseTitle string
		var completed CompletedLesson
		var completedAt time.Time

		if err := rows.Scan(
			&courseID,
			&courseTitle,
			&completed.ID,
			&completed.Title,
			&completed.CourseID,
			&completedAt,
		); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		completed.CompletedAt = completedAt

		if courseID != currentCourseID {
			if currentCourseID != -1 {
				result = append(result, CompletedLessonsByCourse{
					CourseID: currentCourseID,
					Title:    currentCourseTitle,
					Lessons:  currentLessons,
				})
			}

			currentCourseID = courseID
			currentCourseTitle = courseTitle
			currentLessons = []CompletedLesson{}
		}

		currentLessons = append(currentLessons, completed)
	}

	if currentCourseID != -1 {
		result = append(result, CompletedLessonsByCourse{
			CourseID: currentCourseID,
			Title:    currentCourseTitle,
			Lessons:  currentLessons,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
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

	// DEBUG: Verifică dacă variabilele de mediu sunt încărcate
	fmt.Println("\n=== EMAIL CONFIGURATION DEBUG ===")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SENDER_EMAIL")
	senderPassword := os.Getenv("SENDER_PASSWORD")

	fmt.Printf("SMTP_HOST: %s\n", smtpHost)
	fmt.Printf("SMTP_PORT: %s\n", smtpPort)
	fmt.Printf("SENDER_EMAIL: %s\n", senderEmail)
	fmt.Printf("SENDER_PASSWORD: [%d chars]\n", len(senderPassword))
	fmt.Println("================================")

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

	mux.HandleFunc("/api/subscribe", apiSubscribeHandler)
	mux.HandleFunc("/api/register", apiRegisterHandler)
	mux.HandleFunc("/api/login", apiLoginHandler)
	mux.HandleFunc("/api/profile", authMiddleware(apiProfileHandler))
	mux.HandleFunc("/api/progress", authMiddleware(apiSaveProgressHandler))
	mux.HandleFunc("/api/stats", authMiddleware(apiStatsHandler))

	mux.HandleFunc("/api/courses", apiCoursesHandler)
	mux.HandleFunc("/api/course", apiCourseHandler)
	mux.HandleFunc("/api/lessons", apiLessonsHandler)
	mux.HandleFunc("/api/lesson", apiLessonHandler)
	mux.HandleFunc("/api/user", apiUserHandler)
	mux.HandleFunc("/api/userlessons", apiUserLessonsHandler)
	mux.HandleFunc("/api/leaderboard", apiLeaderboardHandler)
	mux.HandleFunc("/api/grinfo/questions", apiGrInfoQuestionsHandler)
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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
