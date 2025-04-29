package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
)

type Question struct {
    ID        int      `json:"id"`
    Text      string   `json:"text"`
    Choices   []Choice `json:"choices"`
    CorrectID int      `json:"correct_id"`
}

type Choice struct {
    ID   int    `json:"id"`
    Text string `json:"text"`
}

type AnswerSubmission struct {
    UserID  int         `json:"user_id"`
    Answers map[int]int `json:"answers"` // question_id -> choice_id
}

var questions = []Question{
    {
        ID:        1,
        Text:      "What's the capital of France?",
        CorrectID: 2,
        Choices: []Choice{
            {ID: 1, Text: "Berlin"},
            {ID: 2, Text: "Paris"},
        },
    },
}

func getQuestions(w http.ResponseWriter, r *http.Request) {
    enableCORS(&w)
    json.NewEncoder(w).Encode(questions)
}

func submitAnswers(w http.ResponseWriter, r *http.Request) {
    enableCORS(&w)
    var submission AnswerSubmission
    err := json.NewDecoder(r.Body).Decode(&submission)
    if err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    score := 0
    for _, q := range questions {
        if submission.Answers[q.ID] == q.CorrectID {
            score++
        }
    }

    // Send score to user service
    body, _ := json.Marshal(map[string]int{
        "user_id": submission.UserID,
        "score":   score,
    })

    http.Post("http://localhost:8001/api/user/update_score", "application/json", bytes.NewBuffer(body))

    w.WriteHeader(http.StatusOK)
}

func enableCORS(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    (*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
    http.HandleFunc("/api/test/questions", getQuestions)
    http.HandleFunc("/api/test/submit", submitAnswers)
    log.Println("Test service running on :8002")
    log.Fatal(http.ListenAndServe(":8002", nil))
}
