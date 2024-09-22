package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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

func main() {
	http.HandleFunc("/user/profile", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getClientProfile(w, r)
		case http.MethodPatch:
			updateClientProfile(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server is running on port 8080...")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getClientProfile(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("clientId")

	clientProfile, ok := database[clientId]

	if !ok || clientId == "" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := ClientProfile{
		Id:    clientProfile.Id,
		Name:  clientProfile.Name,
		Gmail: clientProfile.Gmail,
	}

	json.NewEncoder(w).Encode(response)
}

func updateClientProfile(w http.ResponseWriter, r *http.Request) {

}
