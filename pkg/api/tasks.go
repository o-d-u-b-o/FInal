package api

import (
	"go1f/pkg/db"
	"net/http"
)

// Изменяем структуру ответа
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

	// Преобразуем задачи в нужный формат ответа
	taskResponses := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = TaskResponse{
			ID:      task.ID,
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		}
	}

	writeJSON(w, TasksResponse{Tasks: taskResponses}, http.StatusOK)
}
