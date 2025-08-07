package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"
    _ "github.com/lib/pq"
)

type Question struct {
    ID           int    `json:"id"`
    QuestionText string `json:"question_text"`
    OptionA      string `json:"option_a"`
    OptionB      string `json:"option_b"`
    OptionC      string `json:"option_c"`
    OptionD      string `json:"option_d"`
    CorrectOption string `json:"-"` // hidden from client
}

type NewQuestionRequest struct {
    QuestionText string `json:"question_text"`
    OptionA      string `json:"option_a"`
    OptionB      string `json:"option_b"`
    OptionC      string `json:"option_c"`
    OptionD      string `json:"option_d"`
    CorrectOption string `json:"correct_option"`
}

type AnswerSubmission struct {
    QuestionID     int    `json:"question_id"`
    SelectedOption string `json:"selected_option"`
}

type SubmitRequest struct {
    Answers []AnswerSubmission `json:"answers"`
}

type SubmitResponse struct {
    Score  int    `json:"score"`
    Total  int    `json:"total"`
    Message string `json:"message"`
}

var db *sql.DB

func validateToken(token string) (string, bool, error) {
    url := fmt.Sprintf("http://localhost:8003/api/auth/me?token=%s", token)
    resp, err := http.Get(url)
    if err != nil {
        return "", false, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return "", false, fmt.Errorf("unauthorized")
    }

    var result struct {
        Username string `json:"username"`
        IsStaff  bool   `json:"is_staff"`
    }
    err = json.NewDecoder(resp.Body).Decode(&result)
    return result.Username, result.IsStaff, err
}

func createQuestionHandler(w http.ResponseWriter, r *http.Request) {
    token := r.URL.Query().Get("token")
    username, isStaff, err := validateToken(token)
    if err != nil || !isStaff {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var q NewQuestionRequest
    err = json.NewDecoder(r.Body).Decode(&q)
    if err != nil || q.QuestionText == "" || q.CorrectOption == "" {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    _, err = db.Exec(`
        INSERT INTO questions (question_text, option_a, option_b, option_c, option_d, correct_option)
        VALUES ($1, $2, $3, $4, $5, $6)`,
        q.QuestionText, q.OptionA, q.OptionB, q.OptionC, q.OptionD, strings.ToUpper(q.CorrectOption))

    if err != nil {
        log.Println("Insert error:", err)
        http.Error(w, "Failed to insert question", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Question added"})
}

func listQuestionsHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query(`
        SELECT id, question_text, option_a, option_b, option_c, option_d FROM questions
    `)
    if err != nil {
        http.Error(w, "DB error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var questions []Question
    for rows.Next() {
        var q Question
        err := rows.Scan(&q.ID, &q.QuestionText, &q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD)
        if err != nil {
            http.Error(w, "Scan error", http.StatusInternalServerError)
            return
        }
        questions = append(questions, q)
    }

    json.NewEncoder(w).Encode(questions)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
    token := r.URL.Query().Get("token")
    username, _, err := validateToken(token)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var req SubmitRequest
    err = json.NewDecoder(r.Body).Decode(&req)
    if err != nil || len(req.Answers) == 0 {
        http.Error(w, "Invalid submission", http.StatusBadRequest)
        return
    }

    tx, err := db.Begin()
    if err != nil {
        http.Error(w, "Transaction error", http.StatusInternalServerError)
        return
    }

    var score int
    var submissionID int
    err = tx.QueryRow(`INSERT INTO submissions (username, score) VALUES ($1, 0) RETURNING id`, username).Scan(&submissionID)
    if err != nil {
        tx.Rollback()
        http.Error(w, "Insert submission failed", http.StatusInternalServerError)
        return
    }

    for _, ans := range req.Answers {
        var correct string
        err := tx.QueryRow(`SELECT correct_option FROM questions WHERE id = $1`, ans.QuestionID).Scan(&correct)
        if err != nil {
            tx.Rollback()
            http.Error(w, "Question not found", http.StatusBadRequest)
            return
        }

        isCorrect := strings.ToUpper(ans.SelectedOption) == correct
        if isCorrect {
            score++
        }

        _, err = tx.Exec(`
            INSERT INTO answers (submission_id, question_id, selected_option, is_correct)
            VALUES ($1, $2, $3, $4)`, submissionID, ans.QuestionID, strings.ToUpper(ans.SelectedOption), isCorrect)
        if err != nil {
            tx.Rollback()
            http.Error(w, "Answer insert failed", http.StatusInternalServerError)
            return
        }
    }

    _, err = tx.Exec(`UPDATE submissions SET score = $1 WHERE id = $2`, score, submissionID)
    if err != nil {
        tx.Rollback()
        http.Error(w, "Score update failed", http.StatusInternalServerError)
        return
    }

    err = tx.Commit()
    if err != nil {
        http.Error(w, "Commit failed", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(SubmitResponse{
        Score:  score,
        Total:  len(req.Answers),
        Message: "Submission saved",
    })
}

func main() {
    connStr := "host=localhost port=5432 user=auth_service_user password=yourpassword dbname=auth_service_db sslmode=disable"
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("DB connection failed:", err)
    }
    if err := db.Ping(); err != nil {
        log.Fatal("DB unreachable:", err)
    }

    log.Println("Test service connected to PostgreSQL")

    http.HandleFunc("/questions", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            listQuestionsHandler(w, r)
        } else if r.Method == http.MethodPost {
            createQuestionHandler(w, r)
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    http.HandleFunc("/submit", submitHandler)

    log.Println("Test service running on :8005")
    log.Fatal(http.ListenAndServe(":8005", nil))
}
