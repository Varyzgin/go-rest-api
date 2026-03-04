package main

import (
	"fmt"
	"net/http"
)

/*
Реализовать методы, использовать PG для хранения
GET /user?status=true&name=alex
GET /user/{id}
DELETE /user/{id}
PUT /user/{id}
POST /user/
*/

// urlExample := "postgres://postgres:0000@localhost:5432/rest_api?sslmode=disable"

func main() {
	store := NewMemStore()
	usersHandler := NewUsersHandler(store)
	mux := http.NewServeMux()

	mux.Handle("/user", usersHandler)
	mux.Handle("/user/", usersHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Server initialization error:", err)
	}
}
