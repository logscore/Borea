// This is the file to put our crud and query functions into.

package handlers

import (
	"encoding/json"
	"fmt"
	"io"

	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"Borea/backend/db"
	"Borea/backend/helper"
	"Borea/backend/models"

	"github.com/joho/godotenv"
)

// TODO: change this to GET and find a way to send the query & param data without a POST or URL params
// TODO: change this to only run SELECT statements
func GetItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if db.DB == nil {
		log.Println("Database connection not initialized")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var requestBody models.Request_body

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if !helper.IsAllowedQuery(requestBody.Query, "SELECT") {
		http.Error(w, "Invalid query: only SELECT queries allowed", http.StatusBadRequest)
		return
	}

	// The two below methods prevent SQL injection
	stmt, err := db.DB.Prepare(requestBody.Query)
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(requestBody.Params...)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return
	}

	items := make([]map[string]interface{}, 0)

	values := make([]interface{}, len(columns))
	for i := range values {
		var value interface{}
		values[i] = &value
	}

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return
		}

		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			rowMap[colName] = *(values[i].(*interface{}))
		}

		items = append(items, rowMap)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// TODO: change this to GET and find a way to send the query & param data without a POST or URL splice
// TODO: change this functio nto only run SELECT sql queries
// Note that this returns an interface type, while GetItems returns an array of interface types
func GetItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if db.DB == nil {
		log.Println("Database connection not initialized")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var requestBody models.Request_body

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if !helper.IsAllowedQuery(requestBody.Query, "SELECT") {
		http.Error(w, "Invalid query: only SELECT queries allowed", http.StatusBadRequest)
		return
	}

	// The two below methods prevent SQL injection
	stmt, err := db.DB.Prepare(requestBody.Query)
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(requestBody.Params...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Get column count and names dynamically
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for i := range values {
		valuePtrs[i] = &values[i]
	}

	results := make(map[string]interface{})

	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal(err)
		}

		for i, col := range columns {
			results[col] = values[i]
		}
	}

	if err := rows.Err(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// TODO: adjust this function to only run INSERT sql queries only
// Create a new item
func CreateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if db.DB == nil {
		log.Println("Database connection not initialized")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var requestBody models.Request_body

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if !helper.IsAllowedQuery(requestBody.Query, "INSERT") {
		http.Error(w, "Invalid query: only INSERT queries allowed", http.StatusBadRequest)
		return
	}

	// The two below methods prevent SQL injection
	stmt, err := db.DB.Prepare(requestBody.Query)
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(requestBody.Params...)
	if err != nil {
		log.Fatalf("SQL execution error: %v", err)
	}

	id, _ := result.LastInsertId()
	item := int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Update an existing item
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid query: only PUT queries allowed", http.StatusMethodNotAllowed)
		return
	}

	if db.DB == nil {
		log.Println("Database connection not initialized")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var requestBody models.Request_body

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if !helper.IsAllowedQuery(requestBody.Query, "UPDATE") {
		http.Error(w, "Invalid query: only PUT queries allowed", http.StatusBadRequest)
		return
	}

	// The two below methods prevent SQL injection
	stmt, err := db.DB.Prepare(requestBody.Query)
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(requestBody.Params...)
	if err != nil {
		log.Fatalf("SQL execution error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// // Delete an item
// func DeleteItem(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodDelete {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	id, err := extractIDFromURL(r.URL.Path)
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	_, err = db.DB.Exec("DELETE FROM test_table WHERE id = ?", id)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted"})
// }

func HandleScriptRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("No caller information")
	}

	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

	err := godotenv.Load(filepath.Join(projectRoot, ".env"))
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Token check
	token := r.URL.Query().Get("token")

	TOKEN := os.Getenv("API_TOKEN")

	validToken := false
	if TOKEN == token {
		validToken = true
	}

	if !validToken {
		http.Error(w, "Invalid token in request", http.StatusForbidden)
		return
	}

	// Domain enforcement
	domain := helper.ParseDomainRequest(r)
	if domain == "" {
		http.Error(w, "Error getting domain", http.StatusForbidden)
	}

	DOMAIN := os.Getenv("DOMAIN")

	domainAllowed := false
	if DOMAIN == domain {
		domainAllowed = true
	}

	if !domainAllowed {
		http.Error(w, "Domain not allowed for this token", http.StatusForbidden)
		return
	}

	jsContent, err := os.ReadFile(filepath.Join(projectRoot, "/tracker/Borea.js"))
	if err != nil {
		http.Error(w, "Error reading script file", http.StatusInternalServerError)
		log.Printf("Error reading script file: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/javascript")
	w.Write(jsContent)
}

func PostSessionData(w http.ResponseWriter, r *http.Request) {
	DOMAIN := os.Getenv("DOMAIN")

	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	var sessionData map[string]interface{}
	err = json.Unmarshal(body, &sessionData)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Extract sessionId
	sessionId, ok := sessionData["sessionId"].(string)
	if !ok {
		http.Error(w, "sessionId not found in session data", http.StatusBadRequest)
		return
	}

	// Prepare the SELECT statement
	stmt, err := db.DB.Prepare("SELECT id FROM sessions WHERE sessionId = ?")
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(sessionId).Scan(&id)
	if err != nil {
		if id == 0 {
			_, err = db.DB.Exec(`
            INSERT INTO sessions (
                sessionId, lastActivityTime, userId, userPath,
                sessionDuration, userAgent, referrer, token,
                startTime, screenResolution, language
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
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
				sessionData["language"])
			if err != nil {
				http.Error(w, "Error inserting new session", http.StatusInternalServerError)
				log.Printf("Error inserting new session: %v", err)
				return
			}
		} else {
			log.Printf("Error querying session: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		// Session found, update it
		_, err := db.DB.Exec(`
		UPDATE sessions
		SET
			lastActivityTime = ?,
			userId = ?,
			userPath = ?,
			sessionDuration = ?,
			userAgent = ?,
			referrer = ?,
			token = ?,
			startTime = ?,
			screenResolution = ?,
			language = ?
		WHERE sessionId = ?`,
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
			sessionData["sessionId"],
		)
		if err != nil {
			http.Error(w, "Error updating session", http.StatusInternalServerError)
			log.Printf("Error updating session: %v", err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}
