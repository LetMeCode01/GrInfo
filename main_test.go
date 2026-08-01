package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestHashPasswordAndCheck(t *testing.T) {
	pass := "mysecret"
	hash, err := hashPassword(pass)
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}

	if !checkPasswordHash(pass, hash) {
		t.Fatalf("password should match hash")
	}
}

func TestTokenGenerateValidate(t *testing.T) {
	token, err := generateToken(1, "testuser")
	if err != nil {
		t.Fatalf("token error: %v", err)
	}

	claims, err := validateToken(token)
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}

	if claims.UserID != 1 || claims.Username != "testuser" {
		t.Fatalf("claims mismatch")
	}
}

func TestAuthMiddleware_Valid(t *testing.T) {
	token, _ := generateToken(2, "bob")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler := authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-User-ID") != "2" {
			t.Fatalf("expected user id 2")
		}
	})

	handler(w, req)
}

func TestMain(m *testing.M) {
	// Force correct DSN
	var err error
	db, err = sql.Open("pgx", defaultPostgresDSN)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	if err := initTables(db); err != nil {
		panic(err)
	}

	_, _ = db.Exec("DELETE FROM users WHERE email LIKE 'test-%'")

	// Run tests
	code := m.Run()

	// Cleanup test users after tests
	_, _ = db.Exec("DELETE FROM users WHERE email LIKE 'test-%'")

	_ = db.Close()
	os.Exit(code)
}

func TestApiRegisterHandler_Success(t *testing.T) {
	reqBody := `{"username":"test-user","email":"test-user@example.com","password":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	apiRegisterHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var resp AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.UserID == 0 || resp.Token == "" {
		t.Fatalf("expected valid userId and token, got %+v", resp)
	}
}

func TestApiLoginHandler_Success(t *testing.T) {
	// Ensure user exists
	passwordHash, _ := hashPassword("123456")
	_, _ = db.Exec(
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)",
		"test-login", "test-login@example.com", passwordHash,
	)

	reqBody := `{"email":"test-login@example.com","password":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	apiLoginHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.UserID == 0 || resp.Token == "" {
		t.Fatalf("expected valid userId and token, got %+v", resp)
	}
}

func TestApiLoginHandler_InvalidPassword(t *testing.T) {
	// Ensure user exists
	passwordHash, _ := hashPassword("123456")
	_, _ = db.Exec(
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)",
		"test-login2", "test-login2@example.com", passwordHash,
	)

	reqBody := `{"email":"test-login2@example.com","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	apiLoginHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
