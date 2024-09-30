// This is the file to put our crud and query functions into.

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"Borea/backend/db"
	"Borea/backend/models"
)

// Helper function to extract the ID from the URL. Do we want to do the queries via the URL?
// TODO: throw all the helper functions into a separate file
//
//	func extractIDFromURL(path string) (int, error) {
//		parts := strings.Split(path, "/")
//		idStr := parts[len(parts)-1]
//		return strconv.Atoi(idStr)
//	}

// TODO: change this to GET and find a way to send the query & param data without a POST or URL splice
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

// 	_, err = db.DB.Exec("DELETE FROM items WHERE id = ?", id)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted"})
// }
