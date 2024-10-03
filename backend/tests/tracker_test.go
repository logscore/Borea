package main

import (
	"Borea/backend/db"
	"Borea/backend/handlers"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateSessionTestTable() error {
	_, err := db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sessionId TEXT,
			lastActivityTime TEXT,
			userId TEXT,
			userPath TEXT,
			sessionDuration INTEGER,
			userAgent TEXT,
			referrer TEXT,
			token TEXT,
			startTime TEXT,
			screenResolution TEXT,
			language TEXT
		)`)
	if err != nil {
		log.Printf("Error creating sessions table: %v", err)
		return err
	}

	log.Println("sessions created successfully")
	return nil
}

func TearDownSessionTestTable() error {
	_, err := db.DB.Exec(`DROP TABLE IF EXISTS sessions`)
	if err != nil {
		log.Printf("Error dropping session table: %v", err)
		return err
	}

	log.Println("sessions table dropped successfully")
	return nil
}

func TestHandleScriptRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		os.Setenv("DOMAIN", "http://borea.dev")

		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/script?token=5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Referer", "http://borea.dev")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.HandleScriptRequest)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", status)
		}

		contentType := rr.Header().Get("Content-Type")
		if contentType != "text/javascript" {
			t.Errorf("Expected content type 'text/javascript', got %s", contentType)
		}

		expectedResponse := "// sessionTrack.js"
		if rr.Body.String()[0:18] != expectedResponse {
			t.Errorf("Expected response body '%s', got '%s'", expectedResponse, rr.Body.String())
		}
	})

	t.Run("MethodNotAllowed", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "http://example.com", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.HandleScriptRequest)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("Expected status code 405, got %d", status)
		}
	})

	t.Run("DomainNotAllowed", func(t *testing.T) {
		os.Setenv("DOMAIN", "http://example.com")

		req, err := http.NewRequest(http.MethodGet, "http://other.com/oigjgjgvk/ll", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Referer", "http://other.com/edovinw/wgewfv")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.HandleScriptRequest)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusForbidden {
			t.Errorf("Expected status code 403, got %d", status)
		}
	})

	t.Run("InvalidToken", func(t *testing.T) {
		os.Setenv("DOMAIN", "http://example.com")
		os.Setenv("API_TOKEN", "123456")

		req, err := http.NewRequest(http.MethodGet, "http://example.com/script?token=123467", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Referer", "http://example.com")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.HandleScriptRequest)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusForbidden {
			t.Errorf("Expected status code 403, got %d", status)
		}
	})
}

func TestPostSessionData(t *testing.T) {
	// Set up the test database connection (use your actual DB for real tests)
	err := db.InitDB("../../test_db.sqlite")
	require.NoError(t, err, "Database initialization error")
	defer db.DB.Close()

	err = CreateSessionTestTable()
	require.NoError(t, err, "Failed to create test table")
	defer TearDownSessionTestTable()

	// Mock environment variable for DOMAIN
	os.Setenv("DOMAIN", "http://example.com")

	// Prepare common session data for tests
	sessionData := map[string]interface{}{
		"sessionId":        "12345",
		"lastActivityTime": "2024-10-02T15:04:05Z",
		"userId":           "user123",
		"userPath":         "/dashboard",
		"sessionDuration":  120,
		"userAgent":        "Mozilla/5.0",
		"referrer":         "http://google.com",
		"token":            "abcdefg",
		"startTime":        "2024-10-02T14:04:05Z",
		"screenResolution": "1920x1080",
		"language":         "en",
	}

	// Preflight request (OPTIONS method)
	t.Run("PreflightRequest", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/postSession", nil)
		w := httptest.NewRecorder()

		handlers.PostSessionData(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %v", resp.StatusCode)
		}
	})

	// Method Not Allowed (non-POST)
	t.Run("MethodNotAllowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/postSession", nil)
		w := httptest.NewRecorder()

		handlers.PostSessionData(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405 Method Not Allowed, got %v", resp.StatusCode)
		}
	})

	// Bad Request for invalid JSON
	t.Run("BadRequestInvalidJSON", func(t *testing.T) {
		invalidJSON := []byte(`{invalid json}`)
		req := httptest.NewRequest(http.MethodPost, "/postSession", bytes.NewBuffer(invalidJSON))
		w := httptest.NewRecorder()

		handlers.PostSessionData(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 Bad Request, got %v", resp.StatusCode)
		}
	})

	// Missing sessionId in session data
	t.Run("MissingSessionID", func(t *testing.T) {
		sessionDataWithoutID := sessionData
		delete(sessionDataWithoutID, "sessionId")

		body, _ := json.Marshal(sessionDataWithoutID)
		req := httptest.NewRequest(http.MethodPost, "/postSession", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handlers.PostSessionData(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400 Bad Request, got %v", resp.StatusCode)
		}
	})

	// Insert new session
	t.Run("InsertNewSession", func(t *testing.T) {
		sessionData := map[string]interface{}{
			"sessionId":        "12345",
			"lastActivityTime": "2024-10-02T15:04:05Z",
			"userId":           "user123",
			"userPath":         "/dashboard",
			"sessionDuration":  120,
			"userAgent":        "Mozilla/5.0",
			"referrer":         "http://google.com",
			"token":            "abcdefg",
			"startTime":        "2024-10-02T14:04:05Z",
			"screenResolution": "1920x1080",
			"language":         "en",
		}

		body, _ := json.Marshal(sessionData)
		req := httptest.NewRequest(http.MethodPost, "/postSession", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handlers.PostSessionData(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %v", resp.StatusCode)
		}

		// Verify session is inserted
		var count int
		row := db.DB.QueryRow("SELECT COUNT(*) FROM sessions WHERE sessionId = ?", sessionData["sessionId"])
		err := row.Scan(&count)
		if err != nil || count != 1 {
			t.Errorf("Expected 1 session to be inserted, got %v", count)
		}
	})

	// Update existing session
	t.Run("UpdateExistingSession", func(t *testing.T) {
		sessionData := map[string]interface{}{
			"sessionId":        "123456",
			"lastActivityTime": "2024-10-02T15:04:05Z",
			"userId":           "user123",
			"userPath":         "/dashboard",
			"sessionDuration":  120,
			"userAgent":        "Mozilla/5.0",
			"referrer":         "http://google.com",
			"token":            "abcdefg",
			"startTime":        "2024-10-02T14:04:05Z",
			"screenResolution": "1920x1080",
			"language":         "en",
		}
		// Insert new mock session into DB first (previous test will have inserted a session with sesswionId 12345, this makes a new one with a different sessionId to test with)
		_, err = db.DB.Exec(`
			INSERT INTO sessions (sessionId, lastActivityTime, userId, userPath, sessionDuration, userAgent, referrer, token, startTime, screenResolution, language)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			sessionData["sessionId"],
			sessionData["lastActivityTime"],
			sessionData["userId"],
			sessionData["userPath"],
			sessionData["sessionDuration"],
			sessionData["userAgent"],
			sessionData["referrer"],
			sessionData["token"],
			sessionData["startTime"],
			sessionData["screenResolution"],
			sessionData["language"],
		)
		if err != nil {
			log.Printf("Error inserting test data in to sessions table: %v", err)
			return
		}

		UpdatedSessionData := map[string]interface{}{
			"sessionId":        "123456",
			"lastActivityTime": "2024-10-02T22:44:05Z",
			"userId":           "user123",
			"userPath":         "/dashboard",
			"sessionDuration":  120,
			"userAgent":        "Mozilla/5.0",
			"referrer":         "http://google.com",
			"token":            "abcdefg",
			"startTime":        "2024-10-02T14:04:05Z",
			"screenResolution": "1920x1080",
			"language":         "en",
		}

		body, _ := json.Marshal(UpdatedSessionData)
		req := httptest.NewRequest(http.MethodPost, "/postSession", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handlers.PostSessionData(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 OK, got %v", resp.StatusCode)
		}

		// Verify session is updated
		var updatedActivityTime string
		row := db.DB.QueryRow("SELECT lastActivityTime FROM sessions WHERE sessionId = ?", sessionData["sessionId"])
		err := row.Scan(&updatedActivityTime)
		if err != nil || updatedActivityTime == sessionData["lastActivityTime"] {
			t.Errorf("Expected lastActivityTime to change from '%v', to '%v'", sessionData["lastActivityTime"], updatedActivityTime)
		}
	})
}
