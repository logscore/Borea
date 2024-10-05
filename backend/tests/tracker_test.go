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
	"time"

	"github.com/stretchr/testify/require"
)

func CreateSessionTestTable() error {
	_, err := db.DB.Exec(`
	CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		last_activity_time TIMESTAMP,
		user_id UUID,
		session_id UUID NOT NULL,
		token TEXT,
		start_time TIMESTAMP,
		session_duration INTEGER,
		user_agent TEXT,
		referrer TEXT,
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
	// Set up the test database connection
	err := db.InitDB()
	require.NoError(t, err, "Database initialization error")
	defer db.DB.Close()

	err = CreateSessionTestTable()
	require.NoError(t, err, "Failed to create test table")
	defer TearDownSessionTestTable()

	// Mock environment variable for DOMAIN
	os.Setenv("DOMAIN", "http://example.com")

	// Prepare common session data for tests
	sessionData := map[string]interface{}{
		"sessionId":        "a415c043-3570-4fab-9db0-f040925321be", // Generate a new UUID for the sessionId
		"lastActivityTime": time.Now().Format(time.RFC3339),
		"userId":           "adc0d882-329f-4f83-88b4-38fc593ad217", // Use a valid UUID here
		"sessionDuration":  120,
		"userAgent":        "Mozilla/5.0",
		"referrer":         "http://google.com",
		"token":            "abcdefg",
		"startTime":        time.Now().Format(time.RFC3339),
		"language":         "en",
	}
	// Insert new session
	t.Run("InsertNewSession", func(t *testing.T) {
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
		row := db.DB.QueryRow("SELECT COUNT(*) FROM sessions WHERE session_id = $1", sessionData["sessionId"])
		err := row.Scan(&count)
		if err != nil || count != 1 {
			t.Errorf("Expected 1 session to be inserted, got %v", count)
		}
	})

	// Update existing session
	t.Run("UpdateExistingSession", func(t *testing.T) {
		sessionData := map[string]interface{}{
			"sessionId":        "a415c043-3570-4fab-9db0-f040925321be", // Generate a new UUID for the sessionId
			"lastActivityTime": time.Now().Format(time.RFC3339),
			"userId":           "adc0d882-329f-4f83-88b4-38fc593ad217", // Use a valid UUID here
			"sessionDuration":  120,
			"userAgent":        "Mozilla/5.0",
			"referrer":         "http://google.com",
			"token":            "abcdefg",
			"startTime":        time.Now().Format(time.RFC3339),
			"language":         "en",
		}

		// Insert new mock session into DB first
		_, err = db.DB.Exec(`
			INSERT INTO sessions (session_id, last_activity_time, user_id, session_duration, user_agent, referrer, token, start_time, language)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			sessionData["sessionId"],
			sessionData["lastActivityTime"],
			sessionData["userId"],
			sessionData["sessionDuration"],
			sessionData["userAgent"],
			sessionData["referrer"],
			sessionData["token"],
			sessionData["startTime"],
			sessionData["language"],
		)
		if err != nil {
			log.Printf("Error inserting test data into sessions table: %v", err)
			return
		}

		UpdatedSessionData := map[string]interface{}{
			"sessionId":        sessionData["sessionId"], // Use the same sessionId for update
			"lastActivityTime": "2024-10-02T22:44:05Z",
			"userId":           sessionData["userId"],
			"sessionDuration":  120,
			"userAgent":        "Mozilla/5.0",
			"referrer":         "http://google.com",
			"token":            "abcdefg",
			"startTime":        sessionData["startTime"],
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
		row := db.DB.QueryRow("SELECT last_activity_time FROM sessions WHERE session_id = $1", sessionData["sessionId"])
		err = row.Scan(&updatedActivityTime)
		if err != nil || updatedActivityTime == sessionData["lastActivityTime"] {
			t.Errorf("Expected last_activity_time to change from '%v', to '%v'", sessionData["lastActivityTime"], updatedActivityTime)
		}
	})

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
}
