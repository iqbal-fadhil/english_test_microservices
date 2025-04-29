package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
)

type UserProfile struct {
    ID            int    `json:"id"`
    Username      string `json:"username"`
    Score         int    `json:"score"`
    TestAttempted bool   `json:"test_attempted"`
}

var users = map[int]UserProfile{
    1: {ID: 1, Username: "alice", Score: 0, TestAttempted: false},
}

func getUserProfile(w http.ResponseWriter, r *http.Request) {
    enableCORS(&w)
    idStr := r.URL.Query().Get("user_id")
    id, err := strconv.Atoi(idStr)
    if err != nil || idStr == "" {
        http.Error(w, "Invalid user_id", http.StatusBadRequest)
        return
    }

    user, exists := users[id]
    if !exists {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func updateScore(w http.ResponseWriter, r *http.Request) {
    enableCORS(&w)
    var payload struct {
        UserID int `json:"user_id"`
        Score  int `json:"score"`
    }
    err := json.NewDecoder(r.Body).Decode(&payload)
    if err != nil || payload.UserID == 0 {
        http.Error(w, "Invalid payload", http.StatusBadRequest)
        return
    }

    user, ok := users[payload.UserID]
    if !ok {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    user.Score = payload.Score
    user.TestAttempted = true
    users[payload.UserID] = user
    w.WriteHeader(http.StatusOK)
}

func enableCORS(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    (*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
    http.HandleFunc("/api/user/profile", getUserProfile)
    http.HandleFunc("/api/user/update_score", updateScore)
    log.Println("User service running on :8001")
    log.Fatal(http.ListenAndServe(":8001", nil))
}
