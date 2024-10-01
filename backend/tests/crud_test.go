package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"Borea/backend/db"
	"Borea/backend/handlers"
	"Borea/backend/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"

	_ "github.com/mattn/go-sqlite3"
)

type TestTable struct {
	ID             int     `json:"id"`              // INTEGER type (auto-incremented)
	Name           string  `json:"name"`            // TEXT type
	Age            int     `json:"age"`             // INTEGER type
	Weight         float64 `json:"weight"`          // REAL type (floating-point)
	ProfilePicture []byte  `json:"profile_picture"` // BLOB type (binary data)
	IsActive       int     `json:"is_active"`       // INTEGER used as boolean (1=true, 0=false)
	Notes          string  `json:"notes"`           // TEXT type for additional notes
	CreatedAt      string  `json:"created_at"`      // TEXT type to store date/time as string
}

func CreateTestTable() error {
	_, err := db.DB.Exec(`
	CREATE TABLE IF NOT EXISTS test_table (
		id INTEGER PRIMARY KEY AUTOINCREMENT,  -- INTEGER type (auto-incremented)
		name TEXT,                             -- TEXT type
		age INTEGER,                           -- INTEGER type
		weight REAL,                           -- REAL type (floating-point)
		profile_picture BLOB,                  -- BLOB type (binary data)
		is_active INTEGER,                     -- INTEGER used as boolean (1=true, 0=false)
		notes TEXT,                            -- TEXT type for additional notes
		created_at TEXT                        -- TEXT type to store date/time as string
	);`)
	if err != nil {
		log.Printf("Error creating test_table: %v", err)
		return err
	}

	log.Println("test_table created successfully")
	return nil
}

func TearDownTestTable() error {
	_, err := db.DB.Exec("DROP TABLE IF EXISTS test_table")
	if err != nil {
		log.Printf("Error dropping test_table: %v", err)
		return err
	}

	log.Println("test_table dropped successfully")
	return nil
}

func PopulateTestData() error {
	_, err := db.DB.Exec("DELETE FROM test_table;")
	if err != nil {
		log.Printf("Error clearing test_table: %v", err)
		return err
	}
	stmt, err := db.DB.Prepare(`
	INSERT INTO test_table 
	(name, age, weight, profile_picture, is_active, notes, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	for i := 0; i < 20; i++ {
		name := fmt.Sprintf("Name%d", i+1)
		age := rand.Intn(1000)
		weight := rand.Float64() * 100
		profile_picture := []byte(fmt.Sprintf("Blob data %d", i+1))
		is_active := rand.Intn(2) == 1
		notes := "This is a note"
		created_at := time.Now().Add(-time.Duration(rand.Intn(10000)) * time.Minute)

		_, err := stmt.Exec(name, age, weight, profile_picture, is_active, notes, created_at)
		if err != nil {
			return fmt.Errorf("error inserting fake data: %v", err)
		}
	}

	log.Println("test_table populated successfully")
	return nil
}

func TestCreateItem(t *testing.T) {
	// Initialize the database
	err := db.InitDB("../../test_db.sqlite")
	require.NoError(t, err, "Database initialization error")
	defer db.DB.Close()

	err = CreateTestTable()
	require.NoError(t, err, "Failed to create test table")
	defer TearDownTestTable()

	t.Run("Successful item creation", func(t *testing.T) {
		exampleBlob := []byte{0x89, 0x50, 0x4E, 0x47}
		requestBody := models.Request_body{
			Query: `INSERT INTO test_table 
				(name, age, weight, profile_picture, is_active, notes, created_at)
				VALUES (?, ?, ?, ?, ?, ?, ?);`,
			Params: []interface{}{
				"John Doe", 25, 72.5, exampleBlob, true, "This is a note",
				time.Now().Format("2006-01-02 15:04:05"),
			},
		}

		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err, "Error marshaling request body")

		req, err := http.NewRequest("POST", "/createItem", bytes.NewBuffer(jsonBody))
		require.NoError(t, err, "Error creating request")
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.CreateItem)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

		var responseID int
		err = json.Unmarshal(rr.Body.Bytes(), &responseID)
		assert.NoError(t, err, "Error parsing response body")
		assert.Greater(t, responseID, 0, "Returned ID should be greater than 0")

		// Verify the item was actually created in the database
		var count int
		err = db.DB.QueryRow("SELECT COUNT(*) FROM test_table WHERE id = ?", responseID).Scan(&count)
		assert.NoError(t, err, "Error querying database")
		assert.Equal(t, 1, count, "Item should exist in the database")
	})

	t.Run("Invalid SQL query - non-INSERT query", func(t *testing.T) {
		requestBody := models.Request_body{
			Query:  "DELETE FROM test_table WHERE id = ?",
			Params: []interface{}{1},
		}
	
		jsonBody, _ := json.Marshal(requestBody)
		req, err := http.NewRequest("POST", "/createItem", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
	
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
	
		handler := http.HandlerFunc(handlers.CreateItem)
		handler.ServeHTTP(rr, req)
	
		// Expecting a bad request status because the query is not a INSERT
		assert.Equal(t, http.StatusBadRequest, rr.Code, "Should return bad request for non-INSERT query")
	
		// Optionally, you can check the response body for a more specific error message
		expectedErrorMessage := "Invalid query: only INSERT queries allowed"
		assert.Contains(t, rr.Body.String(), expectedErrorMessage, "Response should contain the correct error message")
	})

	t.Run("Invalid request method", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/createItem", nil)
		require.NoError(t, err, "Error creating request")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.CreateItem)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Should return method not allowed for GET request")
	})

	t.Run("Invalid JSON in request body", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/createItem", bytes.NewBufferString("invalid json"))
		require.NoError(t, err, "Error creating request")
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.CreateItem)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Should return internal server error for invalid JSON")
	})

	t.Run("SQL injection attempt", func(t *testing.T) {
		requestBody := models.Request_body{
			Query:  `INSERT INTO test_table (name) VALUES (?); DROP TABLE test_table; --`,
			Params: []interface{}{"Malicious User"},
		}

		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err, "Error marshaling request body")

		req, err := http.NewRequest("POST", "/createItem", bytes.NewBuffer(jsonBody))
		require.NoError(t, err, "Error creating request")
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.CreateItem)
		handler.ServeHTTP(rr, req)

		// The handler should either return an error or process only the first statement
		if rr.Code == http.StatusOK {
			var responseID int
			err = json.Unmarshal(rr.Body.Bytes(), &responseID)
			assert.NoError(t, err, "Error parsing response body")
			assert.Greater(t, responseID, 0, "Returned ID should be greater than 0")
		} else {
			assert.Equal(t, http.StatusInternalServerError, rr.Code, "Should return internal server error for SQL injection attempt")
		}

		// Verify that the table still exists
		var count int
		err = db.DB.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&count)
		assert.NoError(t, err, "Table should still exist")
	})

	t.Run("Create item with null values", func(t *testing.T) {
		requestBody := models.Request_body{
			Query: `INSERT INTO test_table 
				(name, age, weight, profile_picture, is_active, notes, created_at)
				VALUES (?, ?, ?, ?, ?, ?, ?);`,
			Params: []interface{}{
				"Jane Doe", nil, nil, nil, false, nil,
				time.Now().Format("2006-01-02 15:04:05"),
			},
		}

		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err, "Error marshaling request body")

		req, err := http.NewRequest("POST", "/createItem", bytes.NewBuffer(jsonBody))
		require.NoError(t, err, "Error creating request")
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.CreateItem)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

		var responseID int
		err = json.Unmarshal(rr.Body.Bytes(), &responseID)
		assert.NoError(t, err, "Error parsing response body")
		assert.Greater(t, responseID, 0, "Returned ID should be greater than 0")

		// Verify the item was created with null values
		var name string
		var age, weight sql.NullFloat64
		var isActive bool
		var notes sql.NullString
		err = db.DB.QueryRow("SELECT name, age, weight, is_active, notes FROM test_table WHERE id = ?", responseID).
			Scan(&name, &age, &weight, &isActive, &notes)
		assert.NoError(t, err, "Error querying database")
		assert.Equal(t, "Jane Doe", name)
		assert.False(t, age.Valid, "Age should be null")
		assert.False(t, weight.Valid, "Weight should be null")
		assert.False(t, isActive, "Is_active should be false")
		assert.False(t, notes.Valid, "Notes should be null")
	})
}

func TestGetItems(t *testing.T) {
	// Initialize the database
	err := db.InitDB("../../test_db.sqlite")
	require.NoError(t, err, "Database initialization error")
	defer db.DB.Close()

	err = CreateTestTable()
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}
	defer TearDownTestTable()

	err = PopulateTestData()
	if err != nil {
		t.Fatalf("Failed to populate test table: %v", err)
	}

	if db.DB == nil {
		log.Println("Database connection not initialized")
		return
	}

	requestBody := models.Request_body{
		Query: `
			SELECT id, name, age, weight, profile_picture, is_active, notes, created_at
			FROM test_table
			WHERE id = ?;`, // Add your filtering condition here
		Params: []interface{}{"1"}, // Example parameter to filter by active status
	}

	t.Run("Successful GET request", func(t *testing.T) {

		jsonBody, _ := json.Marshal(requestBody)
		req, err := http.NewRequest("POST", "/getItems", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(handlers.GetItems)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

		var responseItems []TestTable
		err = json.Unmarshal(rr.Body.Bytes(), &responseItems)
		assert.NoError(t, err, "Error parsing response body")
		assert.NotEmpty(t, responseItems, "Response should not be empty")

		// Check the structure of the first item
		firstItem := responseItems[0]
		assert.NotZero(t, firstItem.ID, "ID should not be zero")
		assert.NotEmpty(t, firstItem.Name, "Name should not be empty")
		assert.GreaterOrEqual(t, firstItem.Age, 0, "Age should be non-negative")
		assert.GreaterOrEqual(t, firstItem.Weight, 0.0, "Weight should be non-negative")
		assert.NotNil(t, firstItem.ProfilePicture, "Profile picture should not be nil")
		assert.IsType(t, []byte{}, firstItem.ProfilePicture, "Profile picture should be of type []byte")
		assert.Contains(t, []int{0, 1}, firstItem.IsActive, "IsActive should not be empty")
		assert.NotEmpty(t, firstItem.Notes, "Notes should not be empty")
		assert.NotZero(t, firstItem.CreatedAt, "CreatedAt should not be zero")
	})

	t.Run("Invalid SQL query - non-SELECT query", func(t *testing.T) {
		requestBody := models.Request_body{
			Query:  "DELETE FROM test_table WHERE id = ?", // Invalid, since only SELECT is allowed
			Params: []interface{}{1},
		}
	
		jsonBody, _ := json.Marshal(requestBody)
		req, err := http.NewRequest("POST", "/getItems", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
	
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
	
		handler := http.HandlerFunc(handlers.GetItems)
		handler.ServeHTTP(rr, req)
	
		// Expecting a bad request status because the query is not a SELECT
		assert.Equal(t, http.StatusBadRequest, rr.Code, "Should return bad request for non-SELECT query")
	
		// Optionally, you can check the response body for a more specific error message
		expectedErrorMessage := "Invalid query: only SELECT queries allowed"
		assert.Contains(t, rr.Body.String(), expectedErrorMessage, "Response should contain the correct error message")
	})

	t.Run("Invalid request method", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/getItems", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.GetItems)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Should return method not allowed for GET request")
	})

	t.Run("Invalid JSON in request body", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/getItems", bytes.NewBufferString("invalid json"))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(handlers.GetItems)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Should return internal server error for invalid JSON")
	})

	t.Run("SQL injection attempt", func(t *testing.T) {
		requestBody := models.Request_body{
			Query:  "SELECT * FROM test_table WHERE name = ?",
			Params: []interface{}{"'; DROP TABLE test_table; --"},
		}

		jsonBody, _ := json.Marshal(requestBody)
		req, err := http.NewRequest("POST", "/getItems", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(handlers.GetItems)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Handler should process the request normally")

		// Verify that the table still exists
		var count int
		err = db.DB.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&count)
		assert.NoError(t, err, "Table should still exist")
		assert.Greater(t, count, 0, "Table should not be empty")
	})

	t.Run("Empty result set", func(t *testing.T) {
		requestBody := models.Request_body{
			Query:  "SELECT * FROM test_table WHERE id = ?",
			Params: []interface{}{9999}, // Assuming this ID doesn't exist
		}

		jsonBody, _ := json.Marshal(requestBody)
		req, err := http.NewRequest("POST", "/getItems", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(handlers.GetItems)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Should return OK even for empty result set")

		var responseItems []TestTable
		err = json.Unmarshal(rr.Body.Bytes(), &responseItems)
		assert.NoError(t, err, "Error parsing response body")
		assert.Empty(t, responseItems, "Response should be an empty array")
	})
}

func TestGetItem(t *testing.T) {
	err := db.InitDB("../../test_db.sqlite")
	require.NoError(t, err, "Database initialization error")
	defer db.DB.Close()

	err = CreateTestTable()
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}
	defer TearDownTestTable()

	err = PopulateTestData()
	if err != nil {
		t.Fatalf("Failed to populate test table: %v", err)
	}

	if db.DB == nil {
		log.Println("Database connection not initialized")
		return
	}

	requestBody := models.Request_body{
		Query: `
			SELECT id, name, age, weight, profile_picture, is_active, notes, created_at
			FROM test_table
			WHERE id = ?;`,
		Params: []interface{}{"1"},
	}

	t.Run("Successful GET request", func(t *testing.T) {

		jsonBody, _ := json.Marshal(requestBody)
		req, err := http.NewRequest("POST", "/getItem", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(handlers.GetItem)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

		var responseItem TestTable
		err = json.Unmarshal(rr.Body.Bytes(), &responseItem)
		assert.NoError(t, err, "Error parsing response body")
		assert.NotEmpty(t, responseItem, "Response should not be empty")

		// Check the structure of the first item
		firstItem := responseItem
		assert.NotZero(t, firstItem.ID, "ID should not be zero")
		assert.NotEmpty(t, firstItem.Name, "Name should not be empty")
		assert.GreaterOrEqual(t, firstItem.Age, 0, "Age should be non-negative")
		assert.GreaterOrEqual(t, firstItem.Weight, 0.0, "Weight should be non-negative")
		assert.NotNil(t, firstItem.ProfilePicture, "Profile picture should not be nil")
		assert.IsType(t, []byte{}, firstItem.ProfilePicture, "Profile picture should be of type []byte")
		assert.Contains(t, []int{0, 1}, firstItem.IsActive, "IsActive should be a boolean")
		assert.NotEmpty(t, firstItem.Notes, "Notes should not be empty")
		assert.NotZero(t, firstItem.CreatedAt, "CreatedAt should not be zero")
	})

	t.Run("Invalid SQL query - non-SELECT query", func(t *testing.T) {
		requestBody := models.Request_body{
			Query:  "DELETE FROM test_table WHERE id = ?", // Invalid, since only SELECT is allowed
			Params: []interface{}{1},
		}
	
		jsonBody, _ := json.Marshal(requestBody)
		req, err := http.NewRequest("POST", "/getItem", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
	
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
	
		handler := http.HandlerFunc(handlers.GetItem)
		handler.ServeHTTP(rr, req)
	
		// Expecting a bad request status because the query is not a SELECT
		assert.Equal(t, http.StatusBadRequest, rr.Code, "Should return bad request for non-SELECT query")
	
		// Optionally, you can check the response body for a more specific error message
		expectedErrorMessage := "Invalid query: only SELECT queries allowed"
		assert.Contains(t, rr.Body.String(), expectedErrorMessage, "Response should contain the correct error message")
	})

	t.Run("Invalid request method", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/getItem", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.GetItem)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Should return method not allowed for GET request")
	})

	t.Run("Invalid JSON in request body", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/getItem", bytes.NewBufferString("invalid json"))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(handlers.GetItem)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Should return internal server error for invalid JSON")
	})

	t.Run("SQL injection attempt", func(t *testing.T) {
		requestBody := models.Request_body{
			Query:  "SELECT * FROM test_table WHERE name = ?",
			Params: []interface{}{"Name1; DROP TABLE test_table; --"},
		}

		jsonBody, _ := json.Marshal(requestBody)
		req, err := http.NewRequest("POST", "/getItem", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(handlers.GetItem)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Handler should process the request normally")

		// Verify that the table still exists
		var count int
		err = db.DB.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&count)
		assert.NoError(t, err, "Table should still exist")
		assert.Greater(t, count, 0, "Table should not be empty")
	})

	t.Run("Empty result set", func(t *testing.T) {
		requestBody := models.Request_body{
			Query: `
			SELECT id, name, age, weight, profile_picture, is_active, notes, created_at
			FROM test_table
			WHERE id = ?;`,
			Params: []interface{}{9999}, // Assuming this ID doesn't exist
		}

		jsonBody, _ := json.Marshal(requestBody)
		req, err := http.NewRequest("POST", "/getItem", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(handlers.GetItem)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Should return OK even for empty result set")

		var responseItem TestTable
		err = json.Unmarshal(rr.Body.Bytes(), &responseItem)
		assert.NoError(t, err, "Error parsing response body")
		assert.Empty(t, responseItem, "Response should be an empty array")
	})
}

func TestUpdateItem(t *testing.T) {
	// Initialize the database
	err := db.InitDB("../../test_db.sqlite")
	require.NoError(t, err, "Database initialization error")
	defer db.DB.Close()

	err = CreateTestTable()
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}
	defer TearDownTestTable()

	err = PopulateTestData()
	if err != nil {
		t.Fatalf("Failed to populate test table: %v", err)
	}

	if db.DB == nil {
		log.Println("Database connection not initialized")
		return
	}

	t.Run("Successful item update", func(t *testing.T) {
		requestBody := models.Request_body{
			Query: `SELECT test_table 
				SET name = ?, age = ?, weight = ?, is_active = ?, notes = ?
				WHERE id = ?;`,
			Params: []interface{}{
				"Updated John Doe", 26, 73.5, false, "Updated note",
				1,
			},
		}

		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err, "Error marshaling request body")

		req, err := http.NewRequest("PUT", "/updateItem", bytes.NewBuffer(jsonBody))
		require.NoError(t, err, "Error creating request")
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.UpdateItem)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

		// Verify the item was actually updated in the database
		var name string
		var age int
		var weight float64
		var isActive bool
		var notes string
		err = db.DB.QueryRow("SELECT name, age, weight, is_active, notes FROM test_table WHERE id = 1").
			Scan(&name, &age, &weight, &isActive, &notes)
		assert.NoError(t, err, "Error querying database")
		assert.Equal(t, "Updated John Doe", name)
		assert.Equal(t, 26, age)
		assert.Equal(t, 73.5, weight)
		assert.Equal(t, false, isActive)
		assert.Equal(t, "Updated note", notes)
	})

	t.Run("Invalid SQL query - non-PUT query", func(t *testing.T) {
		requestBody := models.Request_body{
			Query:  "DELETE FROM test_table WHERE id = ?",
			Params: []interface{}{1},
		}
	
		jsonBody, _ := json.Marshal(requestBody)
		req, err := http.NewRequest("POST", "/updateItem", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
	
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
	
		handler := http.HandlerFunc(handlers.UpdateItem)
		handler.ServeHTTP(rr, req)
	
		// Expecting a bad request status because the query is not a PUT
		assert.Equal(t, http.StatusBadRequest, rr.Code, "Should return bad request for non-PUT query")
	
		// Optionally, you can check the response body for a more specific error message
		expectedErrorMessage := "Invalid query: only PUT queries allowed"
		assert.Contains(t, rr.Body.String(), expectedErrorMessage, "Response should contain the correct error message")
	})

	t.Run("Invalid request method", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/updateItem", nil)
		require.NoError(t, err, "Error creating request")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.UpdateItem)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Should return method not allowed for POST request")
	})

	t.Run("Invalid JSON in request body", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/updateItem", bytes.NewBufferString("invalid json"))
		require.NoError(t, err, "Error creating request")
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.UpdateItem)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Should return internal server error for invalid JSON")
	})

	t.Run("SQL injection attempt", func(t *testing.T) {
		requestBody := models.Request_body{
			Query:  `UPDATE test_table SET name = ? WHERE id = 1; DROP TABLE test_table; --`,
			Params: []interface{}{"Malicious User"},
		}

		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err, "Error marshaling request body")

		req, err := http.NewRequest("PUT", "/updateItem", bytes.NewBuffer(jsonBody))
		require.NoError(t, err, "Error creating request")
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.UpdateItem)
		handler.ServeHTTP(rr, req)

		// The handler should either return an error or process only the first statement
		if rr.Code == http.StatusOK {
			// Verify that only the name was updated and the table still exists
			var name string
			err = db.DB.QueryRow("SELECT name FROM test_table WHERE id = 1").Scan(&name)
			assert.NoError(t, err, "Error querying database")
			assert.Equal(t, "Malicious User", name)
		} else {
			assert.Equal(t, http.StatusInternalServerError, rr.Code, "Should return internal server error for SQL injection attempt")
		}

		// Verify that the table still exists
		var count int
		err = db.DB.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&count)
		assert.NoError(t, err, "Table should still exist")
	})

	t.Run("Update item with null values", func(t *testing.T) {
		requestBody := models.Request_body{
			Query: `UPDATE test_table 
				SET name = ?, age = ?, weight = ?, is_active = ?, notes = ?
				WHERE id = ?;`,
			Params: []interface{}{
				"Null Value Test", nil, nil, false, nil,
				2, // Assuming the second item has ID 2
			},
		}

		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err, "Error marshaling request body")

		req, err := http.NewRequest("PUT", "/updateItem", bytes.NewBuffer(jsonBody))
		require.NoError(t, err, "Error creating request")
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.UpdateItem)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

		// Verify the item was updated with null values
		var name string
		var age, weight sql.NullFloat64
		var isActive bool
		var notes sql.NullString
		err = db.DB.QueryRow("SELECT name, age, weight, is_active, notes FROM test_table WHERE id = 2").
			Scan(&name, &age, &weight, &isActive, &notes)
		assert.NoError(t, err, "Error querying database")
		assert.Equal(t, "Null Value Test", name)
		assert.False(t, age.Valid, "Age should be null")
		assert.False(t, weight.Valid, "Weight should be null")
		assert.False(t, isActive, "Is_active should be false")
		assert.False(t, notes.Valid, "Notes should be null")
	})
}
