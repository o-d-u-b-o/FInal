package api

import (
	"log"
	"net/http"
	"time"
)

func Init() {
	corsHandler := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next(w, r)
		}
	}

	http.HandleFunc("/api/task", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addTaskHandler(w, r)
		default:
			taskHandler(w, r)
		}
	}))
	http.HandleFunc("/api/task/done", corsHandler(taskDoneHandler))
	http.HandleFunc("/api/tasks", corsHandler(getTaskListHandler))
	http.HandleFunc("/api/nextdate", corsHandler(nextDateHandler))
	http.HandleFunc("/api/signin", corsHandler(signinHandler))

	log.Println("API handlers initialized")
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, ErrorResponse{Error: "Method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			writeJSON(w, ErrorResponse{Error: "Invalid now parameter"}, http.StatusBadRequest)
			return
		}
	}

	result, err := NextDate(now, dateStr, repeat)
	if err != nil {
		writeJSON(w, ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}
