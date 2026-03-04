package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type UsersHandler struct {
	store usersStore
}

func NewUsersHandler(s usersStore) *UsersHandler {
	return &UsersHandler{
		store: s,
	}
}

func (h *UsersHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	name := q.Get("name")
	statusStr := q.Get("status")

	var filterStatus *bool
	if statusStr != "" {
		b, err := strconv.ParseBool(statusStr)
		if err == nil {
			filterStatus = &b
		}
	}

	users, err := h.store.List()
	if err != nil {
		fmt.Println("Can't eject users:", err)
		return
	}

	var result []User
	for _, u := range users {
		if name != "" && u.Name != name {
			continue
		}
		if filterStatus != nil && u.Status != *filterStatus {
			continue
		}
		result = append(result, u)
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Can't jsonify users:", err)
		return
	}
	fmt.Printf("Get users\n")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[2])
	if err != nil {
		fmt.Println("Can't eject id from url path:", err)
		w.Write([]byte("Uncorrect id/url path"))
	}
	user, err := h.store.Get(int64(id))
	if err != nil {
		fmt.Println("Can't get user:", err)
		return
	}
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Can't jsonify user:", err)
		return
	}
	fmt.Printf("Get user %d\n", id)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println("Can't read client's data:", err)
		w.Write([]byte("Uncorrect body"))
		return
	}
	if err := h.store.Add(user.Id, user); err != nil {
		fmt.Println("Can't add user:", err)
		w.Write([]byte("User hasn't added"))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.Trim(r.URL.Path, "/user/"))
	if err != nil {
		fmt.Println("Can't eject id from url path:", err)
		w.Write([]byte("Uncorrect id/url path"))
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Can't read client's data:", err)
		w.Write([]byte("Uncorrect body"))
		return
	}

	err = h.store.Update(int64(id), user)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Server side error"))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[2])
	if err != nil {
		fmt.Println("Can't eject id from url path:", err)
		w.Write([]byte("Uncorrect id/url path"))
	}
	err = h.store.Remove(int64(id))
	if err != nil {
		fmt.Println("Can't delete user:", err)
		w.Write([]byte("Can't delete"))
	}
	w.WriteHeader(http.StatusOK)
}

var (
	UserRe       = regexp.MustCompile(`^/user/*$`)
	UserReWithID = regexp.MustCompile(`^/user/\d+$`)
)

func (h *UsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && UserRe.MatchString(r.URL.Path):
		h.ListUsers(w, r)
	case r.Method == http.MethodGet && UserReWithID.MatchString(r.URL.Path):
		h.GetUser(w, r)
	case r.Method == http.MethodPost && UserRe.MatchString(r.URL.Path):
		h.CreateUser(w, r)
	case r.Method == http.MethodPut && UserReWithID.MatchString(r.URL.Path):
		h.UpdateUser(w, r)
	case r.Method == http.MethodDelete && UserReWithID.MatchString(r.URL.Path):
		h.DeleteUser(w, r)
	default:
		return
	}
}
