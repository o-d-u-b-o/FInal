package api

import (
	"net/http"
	"strconv"
	"time"

	"go1f/pkg/db"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, ErrorResponse{Error: "Method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, ErrorResponse{Error: "ID is required"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, ErrorResponse{Error: "Invalid ID format"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, ErrorResponse{Error: "task not found"}, http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		nextDateStr, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
			return
		}

		if err := db.UpdateTask(&db.Task{
			ID:      id,
			Date:    nextDateStr,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		}); err != nil {
			writeJSON(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, struct{}{}, http.StatusOK)
}
