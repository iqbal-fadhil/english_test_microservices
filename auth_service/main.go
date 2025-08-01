package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/google/uuid"
    _ "github.com/lib/pq"
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

type User struct {
    Username string
    Password string
}

var tokens = map[string]string{} // token -> username
var db *sql.DB                   // PostgreSQL connection

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

    // Check username and password in PostgreSQL
    var user User
    err = db.QueryRow("SELECT username, password FROM users WHERE username = $1", req.Username).Scan(&user.Username, &user.Password)
    if err == sql.ErrNoRows || user.Password != req.Password {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    } else if err != nil {
        log.Println("DB error:", err)
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    // Generate token
    token := uuid.New().String()
    tokens[token] = user.Username

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
    // Connect to PostgreSQL
    connStr := fmt.Sprintf("host=localhost port=5432 user=auth_service_user password=yourpassword dbname=auth_service_db sslmode=disable")
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Database unreachable:", err)
    }

    log.Println("Connected to PostgreSQL!")

    // Routes
    http.HandleFunc("/api/auth/login", loginHandler)
    http.HandleFunc("/api/auth/validate", validateHandler)

    log.Println("Auth service running on :8003")
    log.Fatal(http.ListenAndServe(":8003", nil))
}
