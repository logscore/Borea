// We either are using protbuff or just regular models.
// This is used to define tables in the db, data to insert into the db, ingoing or outgoing data for API requests, and responses.

package models

type Auth_item struct {
	ID           int    `json:"ID"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

type Request_body struct {
	Query  string        `json:"query"`
	Params []interface{} `json:"params"`
}
