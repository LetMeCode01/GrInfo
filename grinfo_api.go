package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type GrInfoQuestion struct {
	ID                int      `json:"id"`
	Category          string   `json:"category"`
	Dificultate       string   `json:"dificultate"`
	EloRating         int      `json:"eloRating"`
	Enunt             string   `json:"enunt"`
	ExplicatieRaspuns string   `json:"explicatieRaspuns"`
	GraphData         string   `json:"graphData"`
	Optiuni           []string `json:"optiuni"`
	RaspunsCorect     int      `json:"raspunsCorect"`
}

type GrInfoSessionStartRequest struct {
	InitialElo     float64 `json:"initialElo"`
	TotalQuestions int     `json:"totalQuestions"`
	Category       string  `json:"category"`
	Difficulty     string  `json:"difficulty"`
}

type GrInfoSessionAnswerRequest struct {
	SessionID  int     `json:"sessionId"`
	QuestionID int     `json:"questionId"`
	IsCorrect  bool    `json:"isCorrect"`
	EloBefore  float64 `json:"eloBefore"`
	EloAfter   float64 `json:"eloAfter"`
}

type GrInfoSessionFinishRequest struct {
	SessionID      int     `json:"sessionId"`
	FinalElo       float64 `json:"finalElo"`
	CorrectAnswers int     `json:"correctAnswers"`
	TotalQuestions int     `json:"totalQuestions"`
}

type GrInfoSecurityIncidentRequest struct {
	SessionID    int    `json:"sessionId"`
	IncidentType string `json:"incidentType"`
	Description  string `json:"description"`
	EloPenalty   int    `json:"eloPenalty"`
}

type GrInfoCompatIncidentRequest struct {
	Reason     string `json:"reason"`
	Category   string `json:"category"`
	CurrentElo int    `json:"currentElo"`
	EloPenalty int    `json:"eloPenalty"`
}

type GrInfoCompatSessionRequest struct {
	SessionID      int     `json:"sessionId"`
	Category       string  `json:"category"`
	Difficulty     string  `json:"difficulty"`
	InitialElo     float64 `json:"initialElo"`
	FinalElo       float64 `json:"finalElo"`
	CorrectAnswers int     `json:"correctAnswers"`
	TotalQuestions int     `json:"totalQuestions"`
}

func initGrInfoTablesAndSeed(db *sql.DB) error {
	createSQL := `
	CREATE TABLE IF NOT EXISTS grinfo_categories (
		id BIGSERIAL PRIMARY KEY,
		slug TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS grinfo_questions (
		id BIGSERIAL PRIMARY KEY,
		category_id BIGINT NOT NULL REFERENCES grinfo_categories(id) ON DELETE RESTRICT,
		difficulty TEXT NOT NULL CHECK (difficulty IN ('usoara', 'medie', 'grea')),
		elo_rating INTEGER NOT NULL,
		enunt TEXT NOT NULL,
		explicatie_raspuns TEXT NOT NULL,
		graph_data JSONB NOT NULL,
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS grinfo_question_options (
		id BIGSERIAL PRIMARY KEY,
		question_id BIGINT NOT NULL REFERENCES grinfo_questions(id) ON DELETE CASCADE,
		option_index INTEGER NOT NULL CHECK (option_index BETWEEN 0 AND 3),
		option_text TEXT NOT NULL,
		is_correct BOOLEAN NOT NULL DEFAULT FALSE,
		UNIQUE(question_id, option_index)
	);

	CREATE TABLE IF NOT EXISTS grinfo_sessions (
		id BIGSERIAL PRIMARY KEY,
		user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
		category TEXT NOT NULL DEFAULT 'all',
		difficulty TEXT NOT NULL DEFAULT 'all',
		initial_elo NUMERIC(6,2) NOT NULL DEFAULT 1000,
		final_elo NUMERIC(6,2),
		total_questions INTEGER NOT NULL DEFAULT 10,
		correct_answers INTEGER NOT NULL DEFAULT 0,
		started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		ended_at TIMESTAMPTZ
	);

	CREATE TABLE IF NOT EXISTS grinfo_session_answers (
		id BIGSERIAL PRIMARY KEY,
		session_id BIGINT NOT NULL REFERENCES grinfo_sessions(id) ON DELETE CASCADE,
		question_id BIGINT NOT NULL REFERENCES grinfo_questions(id) ON DELETE RESTRICT,
		is_correct BOOLEAN NOT NULL,
		elo_before NUMERIC(6,2) NOT NULL,
		elo_after NUMERIC(6,2) NOT NULL,
		answered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS grinfo_security_logs (
		id BIGSERIAL PRIMARY KEY,
		session_id BIGINT REFERENCES grinfo_sessions(id) ON DELETE CASCADE,
		user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
		description TEXT,
		elo_penalty INTEGER NOT NULL DEFAULT 50,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_grinfo_questions_category ON grinfo_questions(category_id);
	CREATE INDEX IF NOT EXISTS idx_grinfo_questions_elo ON grinfo_questions(elo_rating);
	CREATE INDEX IF NOT EXISTS idx_grinfo_questions_active ON grinfo_questions(is_active);
	CREATE INDEX IF NOT EXISTS idx_grinfo_sessions_user ON grinfo_sessions(user_id);
	CREATE INDEX IF NOT EXISTS idx_grinfo_answers_session ON grinfo_session_answers(session_id);
	CREATE INDEX IF NOT EXISTS idx_grinfo_security_session ON grinfo_security_logs(session_id);
	`

	if _, err := db.Exec(createSQL); err != nil {
		return err
	}

	if _, err := db.Exec(`ALTER TABLE grinfo_sessions ADD COLUMN IF NOT EXISTS difficulty TEXT NOT NULL DEFAULT 'all'`); err != nil {
		return err
	}

	if _, err := db.Exec(`ALTER TABLE grinfo_sessions ADD COLUMN IF NOT EXISTS started_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`); err != nil {
		return err
	}

	if _, err := db.Exec(`
		INSERT INTO grinfo_categories (slug, name)
		VALUES
		('orientate', 'Grafuri orientate'),
		('neorientate', 'Grafuri neorientate')
		ON CONFLICT (slug) DO NOTHING
	`); err != nil {
		return err
	}

	if err := seedGrInfoQuestions(db); err != nil {
		return err
	}

	return cleanupInvalidGrInfoSessions(db)
}

func seedGrInfoQuestions(db *sql.DB) error {
	// Preserve manually curated datasets. Seed only on first boot when no questions exist.
	var existingCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM grinfo_questions`).Scan(&existingCount); err != nil {
		return err
	}
	if existingCount > 0 {
		return nil
	}

	var orientateID, neorientateID int
	if err := db.QueryRow(`SELECT id FROM grinfo_categories WHERE slug='orientate'`).Scan(&orientateID); err != nil {
		return err
	}
	if err := db.QueryRow(`SELECT id FROM grinfo_categories WHERE slug='neorientate'`).Scan(&neorientateID); err != nil {
		return err
	}

	var orientateCount, neorientateCount int
	if err := db.QueryRow(`
		SELECT COUNT(*)
		FROM grinfo_questions q
		JOIN grinfo_categories c ON c.id = q.category_id
		WHERE c.slug='orientate' AND q.is_active=TRUE
	`).Scan(&orientateCount); err != nil {
		return err
	}
	if err := db.QueryRow(`
		SELECT COUNT(*)
		FROM grinfo_questions q
		JOIN grinfo_categories c ON c.id = q.category_id
		WHERE c.slug='neorientate' AND q.is_active=TRUE
	`).Scan(&neorientateCount); err != nil {
		return err
	}

	if orientateCount < 45 {
		if err := insertGeneratedQuestions(db, orientateID, "orientate", orientateCount+1, 45); err != nil {
			return err
		}
	}

	if neorientateCount < 45 {
		if err := insertGeneratedQuestions(db, neorientateID, "neorientate", neorientateCount+1, 45); err != nil {
			return err
		}
	}

	return nil
}

func buildGrInfoQuestionContent(category string, index int) (string, int, string, string, [4]string, int, string) {
	difficulty := "usoara"
	elo := 900 + (index * 12)
	if index > 15 {
		difficulty = "medie"
		elo = 1050 + (index * 14)
	}
	if index > 30 {
		difficulty = "grea"
		elo = 1250 + (index * 16)
	}

	graphData := `{"nodes":[{"id":"n1"},{"id":"n2"},{"id":"n3"}],"edges":[{"source":"n1","target":"n2"},{"source":"n2","target":"n3"}]}`
	var questionText string
	var explanation string
	var options [4]string
	correctIndex := 0

	pattern := (index - 1) % 6
	if category == "orientate" {
		switch pattern {
		case 0:
			questionText = fmt.Sprintf("[Orientate %d] Ce descrie corect un arc într-un graf orientat?", index)
			explanation = "Un arc are sens și merge dintr-un nod sursă către un nod destinație."
			options = [4]string{
				"Leagă două noduri cu sens de la sursă la destinație",
				"Leagă două noduri fără direcție",
				"Poate exista doar în grafuri complete",
				"Nu influențează gradul nodurilor",
			}
		case 1:
			questionText = fmt.Sprintf("[Orientate %d] Ce relație este corectă pentru gradele unui nod?", index)
			explanation = "În grafurile orientate, gradul intern și extern descriu intrările și ieșirile unui nod."
			options = [4]string{
				"Suma gradului intern și extern este gradul total al nodului",
				"Gradul intern este mereu egal cu numărul de noduri",
				"Gradul extern nu poate fi 0",
				"Fiecare nod are exact două arce incidente",
			}
		case 2:
			questionText = fmt.Sprintf("[Orientate %d] Când poate fi aplicată sortarea topologică?", index)
			explanation = "Sortarea topologică se aplică doar pe grafuri orientate fără cicluri."
			options = [4]string{
				"Pe un DAG",
				"Pe orice graf orientat, inclusiv cu cicluri",
				"Doar pe grafuri neorientate",
				"Numai dacă există un singur nod",
			}
		case 3:
			questionText = fmt.Sprintf("[Orientate %d] Ce afirmă corect despre un graf orientat aciclic?", index)
			explanation = "Un DAG nu conține cicluri orientate și este des folosit în dependențe și planificare."
			options = [4]string{
				"Nu conține cicluri orientate",
				"Are obligatoriu un ciclu cu toate nodurile",
				"Este același lucru cu un graf complet",
				"Nu poate avea drumuri",
			}
		case 4:
			questionText = fmt.Sprintf("[Orientate %d] Cum se interpretează matricea de adiacență?", index)
			explanation = "În matricea de adiacență pentru grafuri orientate, poziția aij indică existența unui arc i→j."
			options = [4]string{
				"aij = 1 dacă există un arc de la i la j",
				"aij = 1 doar dacă i și j sunt noduri izolate",
				"aij reprezintă numărul de componente conexe",
				"aij este mereu 0 pe diagonală, indiferent de graf",
			}
		default:
			questionText = fmt.Sprintf("[Orientate %d] Ce este adevărat despre un drum orientat?", index)
			explanation = "Un drum orientat urmează sensul arcelor, dintr-un nod inițial către unul final."
			options = [4]string{
				"Parcurge arcele în sensul lor",
				"Ignoră complet direcția arcelor",
				"Poate folosi aceeași muchie de două ori fără restricții",
				"Este identic cu un ciclu obligatoriu",
			}
		}
	} else {
		switch pattern {
		case 0:
			questionText = fmt.Sprintf("[Neorientate %d] Ce descrie corect o muchie într-un graf neorientat?", index)
			explanation = "În grafurile neorientate, muchiile nu au sens și conectează două noduri în mod simetric."
			options = [4]string{
				"Leagă două noduri fără direcție",
				"Are sens doar într-un singur sens",
				"Poate lega un nod de el însuși în orice caz",
				"Este folosită doar în DAG-uri",
			}
		case 1:
			questionText = fmt.Sprintf("[Neorientate %d] Care este relația corectă pentru suma gradelor?", index)
			explanation = "Într-un graf neorientat, suma gradelor tuturor nodurilor este de două ori numărul muchiilor."
			options = [4]string{
				"Este egală cu 2m",
				"Este egală cu n",
				"Este mereu 0",
				"Este egală cu numărul de componente conexe",
			}
		case 2:
			questionText = fmt.Sprintf("[Neorientate %d] Câte muchii are un arbore cu n noduri?", index)
			explanation = "Un arbore cu n noduri are întotdeauna n - 1 muchii."
			options = [4]string{
				"n - 1 muchii",
				"n muchii",
				"n + 1 muchii",
				"2n muchii",
			}
		case 3:
			questionText = fmt.Sprintf("[Neorientate %d] Câte muchii are un graf complet cu n noduri?", index)
			explanation = "Un graf complet neorientat cu n noduri are n(n - 1)/2 muchii."
			options = [4]string{
				"n(n - 1)/2",
				"n - 1",
				"n",
				"n^2",
			}
		case 4:
			questionText = fmt.Sprintf("[Neorientate %d] Ce afirmă corect despre un graf conex?", index)
			explanation = "Conexitatea spune că între orice două noduri există cel puțin un drum."
			options = [4]string{
				"Între orice două noduri există cel puțin un drum",
				"Nu poate conține cicluri",
				"Are obligatoriu exact o muchie",
				"Este mereu complet",
			}
		default:
			questionText = fmt.Sprintf("[Neorientate %d] Ce este adevărat despre un graf simplu?", index)
			explanation = "Un graf simplu nu are bucle și nu are muchii multiple între aceleași noduri."
			options = [4]string{
				"Nu are bucle și nici muchii multiple",
				"Trebuie să fie complet",
				"Are obligatoriu bucle pe fiecare nod",
				"Poate avea mai multe muchii paralele între aceleași noduri",
			}
		}
	}

	shuffledOptions, shuffledCorrectIndex := shuffleGrInfoOptions(category, index, options, correctIndex)

	return difficulty, elo, questionText, explanation, shuffledOptions, shuffledCorrectIndex, graphData
}

func shuffleGrInfoOptions(category string, index int, options [4]string, correctIndex int) ([4]string, int) {
	type optionItem struct {
		text      string
		isCorrect bool
	}

	items := []optionItem{
		{text: options[0], isCorrect: correctIndex == 0},
		{text: options[1], isCorrect: correctIndex == 1},
		{text: options[2], isCorrect: correctIndex == 2},
		{text: options[3], isCorrect: correctIndex == 3},
	}

	seedInput := fmt.Sprintf("%s:%d", category, index)
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(seedInput))
	rng := rand.New(rand.NewSource(int64(hasher.Sum32())))
	rng.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	var shuffled [4]string
	shuffledCorrectIndex := 0
	for i, item := range items {
		shuffled[i] = item.text
		if item.isCorrect {
			shuffledCorrectIndex = i
		}
	}

	return shuffled, shuffledCorrectIndex
}

func insertGeneratedQuestions(db *sql.DB, categoryID int, category string, fromIndex int, toIndex int) error {
	for i := fromIndex; i <= toIndex; i++ {
		difficulty, elo, questionText, explanation, options, correctIndex, graphData := buildGrInfoQuestionContent(category, i)

		var qid int
		if err := db.QueryRow(`
			INSERT INTO grinfo_questions (category_id, difficulty, elo_rating, enunt, explicatie_raspuns, graph_data, is_active)
			VALUES ($1, $2, $3, $4, $5, $6::jsonb, TRUE)
			RETURNING id
		`, categoryID, difficulty, elo, questionText, explanation, graphData).Scan(&qid); err != nil {
			return err
		}

		for idx := 0; idx < 4; idx++ {
			if _, err := db.Exec(`
				INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
				VALUES ($1, $2, $3, $4)
			`, qid, idx, options[idx], idx == correctIndex); err != nil {
				return err
			}
		}
	}

	return nil
}

func repairGrInfoQuestionOptions(db *sql.DB) error {
	rows, err := db.Query(`
		SELECT q.id, c.slug
		FROM grinfo_questions q
		JOIN grinfo_categories c ON c.id = q.category_id
		ORDER BY q.id
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	categoryCounters := map[string]int{}
	for rows.Next() {
		var questionID int
		var categorySlug string
		if err := rows.Scan(&questionID, &categorySlug); err != nil {
			continue
		}

		categoryCounters[categorySlug]++
		index := categoryCounters[categorySlug]
		difficulty, _, questionText, explanation, options, correctIndex, graphData := buildGrInfoQuestionContent(categorySlug, index)

		if _, err := db.Exec(`
			UPDATE grinfo_questions
			SET difficulty = $1,
				enunt = $2,
				explicatie_raspuns = $3,
				graph_data = $4::jsonb
			WHERE id = $5
		`, difficulty, questionText, explanation, graphData, questionID); err != nil {
			return err
		}

		if _, err := db.Exec(`DELETE FROM grinfo_question_options WHERE question_id = $1`, questionID); err != nil {
			return err
		}

		for idx := 0; idx < 4; idx++ {
			if _, err := db.Exec(`
				INSERT INTO grinfo_question_options (question_id, option_index, option_text, is_correct)
				VALUES ($1, $2, $3, $4)
			`, questionID, idx, options[idx], idx == correctIndex); err != nil {
				return err
			}
		}
	}

	return rows.Err()
}

func cleanupInvalidGrInfoSessions(db *sql.DB) error {
	_, err := db.Exec(`
		DELETE FROM grinfo_sessions
		WHERE final_elo IS NOT NULL
		  AND final_elo <= 0
	`)
	return err
}

func apiGrInfoCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT slug, name FROM grinfo_categories ORDER BY name`)
	if err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := make([]map[string]string, 0)
	for rows.Next() {
		var slug, name string
		if err := rows.Scan(&slug, &name); err == nil {
			items = append(items, map[string]string{"slug": slug, "name": name})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"categories": items})
}

func apiGrInfoQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		category = "all"
	}
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	query := `
		SELECT q.id, c.slug, q.difficulty, q.elo_rating, q.enunt, q.explicatie_raspuns, q.graph_data
		FROM grinfo_questions q
		JOIN grinfo_categories c ON c.id = q.category_id
		WHERE q.is_active = TRUE
	`
	args := []interface{}{}
	if category != "all" {
		query += ` AND c.slug = $1`
		args = append(args, category)
		query += ` ORDER BY q.elo_rating ASC, q.id ASC LIMIT $2`
		args = append(args, limit)
	} else {
		query += ` ORDER BY q.elo_rating ASC, q.id ASC LIMIT $1`
		args = append(args, limit)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	questions := make([]GrInfoQuestion, 0)
	for rows.Next() {
		var q GrInfoQuestion
		if err := rows.Scan(&q.ID, &q.Category, &q.Dificultate, &q.EloRating, &q.Enunt, &q.ExplicatieRaspuns, &q.GraphData); err != nil {
			continue
		}

		optionRows, err := db.Query(`
			SELECT option_index, option_text, is_correct
			FROM grinfo_question_options
			WHERE question_id = $1
			ORDER BY option_index ASC
		`, q.ID)
		if err != nil {
			continue
		}

		q.Optiuni = make([]string, 0, 4)
		q.RaspunsCorect = 0
		for optionRows.Next() {
			var idx int
			var text string
			var isCorrect bool
			if err := optionRows.Scan(&idx, &text, &isCorrect); err == nil {
				q.Optiuni = append(q.Optiuni, text)
				if isCorrect {
					q.RaspunsCorect = idx
				}
			}
		}
		optionRows.Close()
		questions = append(questions, q)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"questions": questions,
		"count":     len(questions),
	})
}

func apiGrInfoSessionStartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GrInfoSessionStartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req = GrInfoSessionStartRequest{}
	}

	if req.InitialElo <= 0 {
		req.InitialElo = 1000
	}
	if req.TotalQuestions <= 0 {
		req.TotalQuestions = 10
	}
	if req.Category == "" {
		req.Category = "all"
	}
	if req.Difficulty == "" {
		req.Difficulty = "all"
	}

	var userID sql.NullInt64
	if uidStr := r.Header.Get("X-User-ID"); uidStr != "" {
		if uid, err := strconv.Atoi(uidStr); err == nil {
			userID = sql.NullInt64{Int64: int64(uid), Valid: true}
		}
	}

	var sessionID int
	err := db.QueryRow(`
		INSERT INTO grinfo_sessions (user_id, category, difficulty, initial_elo, total_questions)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, userID, req.Category, req.Difficulty, req.InitialElo, req.TotalQuestions).Scan(&sessionID)
	if err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"sessionId":      sessionID,
		"initialElo":     req.InitialElo,
		"totalQuestions": req.TotalQuestions,
	})
}

func apiGrInfoSessionAnswerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GrInfoSessionAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SessionID == 0 || req.QuestionID == 0 {
		jsonError(w, "Missing session/question id", http.StatusBadRequest)
		return
	}

	if _, err := db.Exec(`
		INSERT INTO grinfo_session_answers (session_id, question_id, is_correct, elo_before, elo_after)
		VALUES ($1, $2, $3, $4, $5)
	`, req.SessionID, req.QuestionID, req.IsCorrect, req.EloBefore, req.EloAfter); err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
}

func apiGrInfoSessionFinishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GrInfoSessionFinishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SessionID == 0 {
		jsonError(w, "Missing sessionId", http.StatusBadRequest)
		return
	}
	if req.TotalQuestions <= 0 {
		req.TotalQuestions = 10
	}

	if _, err := db.Exec(`
		UPDATE grinfo_sessions
		SET final_elo = $1,
			correct_answers = $2,
			total_questions = $3,
			ended_at = NOW()
		WHERE id = $4
	`, req.FinalElo, req.CorrectAnswers, req.TotalQuestions, req.SessionID); err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
}

func apiGrInfoSecurityIncidentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GrInfoSecurityIncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.IncidentType == "" {
		req.IncidentType = "UNKNOWN"
	}
	if req.EloPenalty <= 0 {
		req.EloPenalty = 50
	}

	var userID sql.NullInt64
	if uidStr := r.Header.Get("X-User-ID"); uidStr != "" {
		if uid, err := strconv.Atoi(uidStr); err == nil {
			userID = sql.NullInt64{Int64: int64(uid), Valid: true}
		}
	}

	if _, err := db.Exec(`
		INSERT INTO grinfo_security_logs (session_id, user_id, incident_type, description, elo_penalty)
		VALUES ($1, $2, $3, $4, $5)
	`, req.SessionID, userID, req.IncidentType, req.Description, req.EloPenalty); err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
}

func apiGrInfoProfileHandler(w http.ResponseWriter, r *http.Request) {
	uidStr := r.Header.Get("X-User-ID")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		jsonError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var currentElo float64 = 1000
	_ = db.QueryRow(`
		SELECT COALESCE(final_elo, initial_elo, 1000)
		FROM grinfo_sessions
		WHERE user_id = $1
		ORDER BY id DESC
		LIMIT 1
	`, uid).Scan(&currentElo)

	var totalSessions, totalQuestionsAnswered, totalCorrectAnswers int
	_ = db.QueryRow(`
		SELECT COUNT(*), COALESCE(SUM(total_questions), 0), COALESCE(SUM(correct_answers), 0)
		FROM grinfo_sessions gs
		WHERE user_id = $1
	`, uid).Scan(&totalSessions, &totalQuestionsAnswered, &totalCorrectAnswers)

	accuracy := 0.0
	if totalQuestionsAnswered > 0 {
		accuracy = (float64(totalCorrectAnswers) / float64(totalQuestionsAnswered)) * 100
	}

	rows, err := db.Query(`
		SELECT
			id,
			category,
			COALESCE(NULLIF(difficulty, ''), 'all') AS difficulty,
			COALESCE(final_elo, initial_elo, 1000) AS display_elo,
			correct_answers,
			total_questions,
			COALESCE(started_at, ended_at, NOW()) AS started_at,
			ended_at,
			(ended_at IS NULL) AS is_in_progress
		FROM grinfo_sessions
		WHERE user_id = $1
		ORDER BY id DESC
		LIMIT 10
	`, uid)
	if err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	history := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int
		var category string
		var difficulty string
		var finalElo float64
		var correct, total int
		var started time.Time
		var ended sql.NullTime
		var inProgress bool
		if err := rows.Scan(&id, &category, &difficulty, &finalElo, &correct, &total, &started, &ended, &inProgress); err == nil {
			endedStr := ""
			if ended.Valid {
				endedStr = ended.Time.Format(time.RFC3339)
			}
			history = append(history, map[string]interface{}{
				"sessionId":      id,
				"category":       category,
				"difficulty":     difficulty,
				"finalElo":       finalElo,
				"correctAnswers": correct,
				"totalQuestions": total,
				"startedAt":      started.Format(time.RFC3339),
				"endedAt":        endedStr,
				"isInProgress":   inProgress,
			})
		}
	}

	eloRows, err := db.Query(`
		SELECT
			COALESCE(final_elo, initial_elo, 1000) AS display_elo,
			COALESCE(ended_at, started_at, NOW()) AS at_time
		FROM grinfo_sessions
		WHERE user_id = $1
		ORDER BY id DESC
	`, uid)
	if err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer eloRows.Close()

	eloHistory := make([]map[string]interface{}, 0)
	for eloRows.Next() {
		var elo float64
		var atTime time.Time
		if err := eloRows.Scan(&elo, &atTime); err == nil {
			eloHistory = append(eloHistory, map[string]interface{}{
				"elo": elo,
				"at":  atTime.Format(time.RFC3339),
			})
		}
	}

	var securityEvents int
	_ = db.QueryRow(`SELECT COUNT(*) FROM grinfo_security_logs WHERE user_id = $1`, uid).Scan(&securityEvents)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"currentElo":             currentElo,
		"securityEvents":         securityEvents,
		"history":                history,
		"eloHistory":             eloHistory,
		"totalSessions":          totalSessions,
		"totalCorrectAnswers":    totalCorrectAnswers,
		"totalQuestionsAnswered": totalQuestionsAnswered,
		"accuracy":               accuracy,
	})
}

func apiGrInfoIncidentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GrInfoCompatIncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Reason == "" {
		req.Reason = "ANTI_CHEAT_EVENT"
	}
	if req.EloPenalty <= 0 {
		req.EloPenalty = 50
	}

	var userID sql.NullInt64
	if uidStr := r.Header.Get("X-User-ID"); uidStr != "" {
		if uid, err := strconv.Atoi(uidStr); err == nil {
			userID = sql.NullInt64{Int64: int64(uid), Valid: true}
		}
	}

	description := fmt.Sprintf("%s | category=%s | elo=%d", req.Reason, req.Category, req.CurrentElo)
	if _, err := db.Exec(`
		INSERT INTO grinfo_security_logs (session_id, user_id, incident_type, description, elo_penalty)
		VALUES (NULL, $1, $2, $3, $4)
	`, userID, req.Reason, description, req.EloPenalty); err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
}

func apiGrInfoSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uidStr := r.Header.Get("X-User-ID")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		jsonError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req GrInfoCompatSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Category == "" {
		req.Category = "all"
	}
	if req.TotalQuestions <= 0 {
		req.TotalQuestions = 10
	}
	if req.InitialElo <= 0 {
		req.InitialElo = 1000
	}

	if req.SessionID > 0 {
		if _, err := db.Exec(`
			UPDATE grinfo_sessions
			SET final_elo = $1,
				correct_answers = $2,
				total_questions = $3,
				difficulty = CASE
					WHEN $4 IN ('usoara', 'medie', 'grea') THEN $4
					ELSE difficulty
				END,
				ended_at = NOW()
			WHERE id = $5 AND user_id = $6
		`, req.FinalElo, req.CorrectAnswers, req.TotalQuestions, req.Difficulty, req.SessionID, uid); err != nil {
			jsonError(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "sessionId": req.SessionID})
		return
	}

	if req.Difficulty == "" {
		req.Difficulty = "all"
	}

	var sessionID int
	if err := db.QueryRow(`
		INSERT INTO grinfo_sessions (user_id, category, difficulty, initial_elo, final_elo, total_questions, correct_answers)
		VALUES ($1, $2, $3, $4, $4, $5, 0)
		RETURNING id
	`, uid, req.Category, req.Difficulty, req.InitialElo, req.TotalQuestions).Scan(&sessionID); err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "sessionId": sessionID, "initialElo": req.InitialElo})
}

func apiGrInfoSessionProgressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uidStr := r.Header.Get("X-User-ID")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		jsonError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req GrInfoCompatSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SessionID == 0 {
		jsonError(w, "Missing sessionId", http.StatusBadRequest)
		return
	}
	if req.TotalQuestions < 0 {
		req.TotalQuestions = 0
	}
	if req.InitialElo <= 0 {
		req.InitialElo = 1000
	}

	if _, err := db.Exec(`
		UPDATE grinfo_sessions
		SET final_elo = $1,
			correct_answers = $2,
			total_questions = $3
		WHERE id = $4 AND user_id = $5 AND ended_at IS NULL
	`, req.FinalElo, req.CorrectAnswers, req.TotalQuestions, req.SessionID, uid); err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
}
