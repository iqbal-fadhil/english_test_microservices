// package main

// import (
//     "encoding/json"
//     "github.com/google/uuid"
//     "log"
//     "net/http"
// )

// type LoginRequest struct {
//     Username string `json:"username"`
// }

// type LoginResponse struct {
//     Token string `json:"token"`
// }

// var tokens = map[string]string{} // token -> username

// func loginHandler(w http.ResponseWriter, r *http.Request) {
//     enableCORS(&w)

//     var req LoginRequest
//     err := json.NewDecoder(r.Body).Decode(&req)
//     if err != nil || req.Username == "" {
//         http.Error(w, "Invalid request", http.StatusBadRequest)
//         return
//     }

//     // In real apps: check password and user validity from DB or user_service

//     token := uuid.New().String()
//     tokens[token] = req.Username

//     json.NewEncoder(w).Encode(LoginResponse{Token: token})
// }

// func validateHandler(w http.ResponseWriter, r *http.Request) {
//     enableCORS(&w)

//     token := r.URL.Query().Get("token")
//     username, exists := tokens[token]
//     if !exists {
//         http.Error(w, "Invalid token", http.StatusUnauthorized)
//         return
//     }

//     json.NewEncoder(w).Encode(map[string]string{
//         "username": username,
//     })
// }

// func enableCORS(w *http.ResponseWriter) {
//     (*w).Header().Set("Access-Control-Allow-Origin", "*")
//     (*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
//     (*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
// }

// func main() {
//     http.HandleFunc("/api/auth/login", loginHandler)
//     http.HandleFunc("/api/auth/validate", validateHandler)
//     log.Println("Auth service running on :8003")
//     log.Fatal(http.ListenAndServe(":8003", nil))
// }


package main

import (
    "encoding/json"
    "github.com/google/uuid"
    "log"
    "net/http"
)

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token string `json:"token"`
}

var tokens = map[string]string{}      // token -> username
var users = map[string]string{        // username -> password (dummy user DB)
    "admin": "admin123",
    "user1": "pass123",
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    enableCORS(&w)

    if r.Method == http.MethodOptions {
        return // handle CORS preflight
    }

    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req LoginRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil || req.Username == "" || req.Password == "" {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Check username & password
    storedPassword, userExists := users[req.Username]
    if !userExists || storedPassword != req.Password {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    // Generate token
    token := uuid.New().String()
    tokens[token] = req.Username

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
    enableCORS(&w)

    if r.Method == http.MethodOptions {
        return // handle CORS preflight
    }

    token := r.URL.Query().Get("token")
    username, exists := tokens[token]
    if !exists {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "username": username,
    })
}

func enableCORS(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    (*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func main() {
    http.HandleFunc("/api/auth/login", loginHandler)
    http.HandleFunc("/api/auth/validate", validateHandler)
    log.Println("Auth service running on :8003")
    log.Fatal(http.ListenAndServe(":8003", nil))
}
