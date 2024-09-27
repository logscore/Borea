// This is the file to put our crud and query functions into.

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"Borea/backend/db"
)

type Auth_item struct {
	ID           int    `json:"ID"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

// Helper function to extract the ID from the URL in GetItem
// TODO: throw all the helper functions into a separate file
//
//	func extractIDFromURL(path string) (int, error) {
//		parts := strings.Split(path, "/")
//		idStr := parts[len(parts)-1]
//		return strconv.Atoi(idStr)
//	}
func GetItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Make this check with every function to avoid null pointers.
	if db.DB == nil {
		log.Println("Database connection not initialized")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Define a struct to hold the incoming request data
	var requestBody struct {
		Query  string        `json:"query"`
		Params []interface{} `json:"params"`
	}

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Println("Error reading request body:", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Use a parameterized query to avoid SQL injection
	// Make this work by passing in a SQL query from the reqeust body. How can I prevent injections?
	query := "SELECT ID, username, passwordHash FROM admin_users WHERE username = ?"
	rows, err := db.DB.Query(query, requestBody.Params)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close() // Ensure the rows are closed after we're done

	var items []Auth_item
	for rows.Next() {
		var item Auth_item
		if err := rows.Scan(&item.ID, &item.Username, &item.PasswordHash); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set response header and encode items as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// // Get a single item by ID
// func GetItem(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	id, err := extractIDFromURL(r.URL.Path)
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	var item Auth_item
// 	err = db.DB.QueryRow("SELECT id, name, price FROM items WHERE id = ?", id).Scan(&item.ID, &item.username, &item.password)
// 	if err == sql.ErrNoRows {
// 		http.Error(w, "Item not found", http.StatusNotFound)
// 		return
// 	} else if err != nil {
// 		log.Fatal(err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(item)
// }

// // Create a new item
// func CreateItem(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var item Auth_item
// 	err := json.NewDecoder(r.Body).Decode(&item)
// 	if err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}

// 	result, err := db.DB.Exec("INSERT INTO items (name, price) VALUES (?, ?)", item.username, item.password)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	id, _ := result.LastInsertId()
// 	item.ID = int(id)

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(item)
// }

// // Update an existing item
// func UpdateItem(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPut {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	id, err := extractIDFromURL(r.URL.Path)
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	var item Auth_item
// 	err = json.NewDecoder(r.Body).Decode(&item)
// 	if err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}

// 	_, err = db.DB.Exec("UPDATE items SET name = ?, price = ? WHERE id = ?", item.username, item.password, id)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	item.ID = id
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(item)
// }

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
