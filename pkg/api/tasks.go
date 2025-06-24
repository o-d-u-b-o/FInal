package api

import (
	"go1f/pkg/db"
	"net/http"
	"strconv"
)

type TaskResponse struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

func getTaskListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, ErrorResponse{Error: "Method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSON(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	// Исправленный формат ответа
	taskResponses := make([]map[string]string, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = map[string]string{
			"id":      strconv.FormatInt(task.ID, 10),
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		}
	}

	writeJSON(w, map[string]interface{}{"tasks": taskResponses}, http.StatusOK)
}
