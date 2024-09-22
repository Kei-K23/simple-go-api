package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

// In-memory database
type ClientProfile struct {
	Id    string
	Name  string
	Gmail string
	Token string
}

var database = map[string]ClientProfile{
	"user1": {
		Id:    "user1",
		Name:  "John Doe",
		Gmail: "johndoe@gmail.com",
		Token: "john123",
	},
	"user2": {
		Id:    "user2",
		Name:  "Michael",
		Gmail: "michael@gmail.com",
		Token: "michael123",
	},
}

// Middleware functions
type Middleware func(http.HandlerFunc) http.HandlerFunc

func tokenAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Profile create POST method is public
		if r.Method == http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}

		clientId := r.URL.Query().Get("clientId")

		clientProfile, ok := database[clientId]

		if !ok || clientId == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		token := r.Header.Get("Authorization")
		if !isValidToken(&clientProfile, token) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "clientProfile", clientProfile)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

func isValidToken(clientProfile *ClientProfile, token string) bool {
	if strings.HasPrefix(token, "Bearer ") {
		return strings.TrimPrefix(token, "Bearer ") == clientProfile.Token
	}
	// Invalid token format
	return false
}

var middlewares = []Middleware{
	tokenAuthMiddleware,
}

// Main function
func main() {

	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createClientProfile(w, r)
		case http.MethodGet:
			getClientProfile(w, r)
		case http.MethodPatch:
			updateClientProfile(w, r)
		case http.MethodDelete:
			deleteClientProfile(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}

	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	http.HandleFunc("/user/profile", handler)

	log.Println("Server is running on port 8080...")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler methods Get, PATCH
func createClientProfile(w http.ResponseWriter, r *http.Request) {
	var payloadData ClientProfile

	if err := json.NewDecoder(r.Body).Decode(&payloadData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	userId := fmt.Sprintf("user%d", len(database)+1)
	response := ClientProfile{
		Id:    userId,
		Name:  payloadData.Name,
		Gmail: payloadData.Gmail,
		Token: fmt.Sprintf("%s%d", payloadData.Name, rand.Int()),
	}

	database[userId] = response
	json.NewEncoder(w).Encode(response)
}

func getClientProfile(w http.ResponseWriter, r *http.Request) {
	clientProfile := r.Context().Value("clientProfile").(ClientProfile)

	w.Header().Set("Content-Type", "application/json")

	response := ClientProfile{
		Id:    clientProfile.Id,
		Name:  clientProfile.Name,
		Gmail: clientProfile.Gmail,
		Token: clientProfile.Token,
	}

	json.NewEncoder(w).Encode(response)
}

func updateClientProfile(w http.ResponseWriter, r *http.Request) {
	clientProfile := r.Context().Value("clientProfile").(ClientProfile)

	var payloadData ClientProfile
	if err := json.NewDecoder(r.Body).Decode(&payloadData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if payloadData.Gmail != "" {
		clientProfile.Gmail = payloadData.Gmail
	}

	if payloadData.Name != "" {
		clientProfile.Name = payloadData.Name
	}

	w.Header().Set("Content-Type", "application/json")

	response := ClientProfile{
		Id:    clientProfile.Id,
		Name:  clientProfile.Name,
		Gmail: clientProfile.Gmail,
		Token: clientProfile.Token,
	}

	json.NewEncoder(w).Encode(response)
}

func deleteClientProfile(w http.ResponseWriter, r *http.Request) {
	clientProfile := r.Context().Value("clientProfile").(ClientProfile)

	delete(database, clientProfile.Id)

	w.WriteHeader(http.StatusNoContent)
}
