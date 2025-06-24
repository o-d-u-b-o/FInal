package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go1f/pkg/db"
)

type taskResponse struct {
	ID    string `json:"id"`
	Error string `json:"error,omitempty"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, taskResponse{Error: "Method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	var response taskResponse

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response.Error = "Invalid JSON format"
		writeJSON(w, response, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		response.Error = "Task title is required"
		writeJSON(w, response, http.StatusBadRequest)
		return
	}

	// Валидация правила повторения
	if task.Repeat != "" {
		if err := validateRepeatRule(task.Repeat); err != nil {
			response.Error = err.Error()
			writeJSON(w, response, http.StatusBadRequest)
			return
		}
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	} else {
		if _, err := time.Parse("20060102", task.Date); err != nil {
			response.Error = "Invalid date format, expected YYYYMMDD"
			writeJSON(w, response, http.StatusBadRequest)
			return
		}
	}

	// Проверка что дата не в прошлом
	date, _ := time.Parse("20060102", task.Date)
	if Before(date, now) {
		if task.Repeat == "" {
			// Для неповторяющихся задач - установить сегодняшнюю дату
			task.Date = now.Format("20060102")
		} else {
			// Для повторяющихся задач - вычислить следующую валидную дату
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				response.Error = err.Error()
				writeJSON(w, response, http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		response.Error = "Failed to add task to database"
		writeJSON(w, response, http.StatusInternalServerError)
		return
	}

	response.ID = strconv.FormatInt(id, 10)
	writeJSON(w, response, http.StatusOK)
}

func validateRepeatRule(repeat string) error {
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return nil
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return errors.New("invalid daily repeat format")
		}
		if _, err := strconv.Atoi(parts[1]); err != nil {
			return errors.New("invalid days count")
		}
	case "w":
		if len(parts) != 2 {
			return errors.New("invalid weekly repeat format")
		}
		for _, day := range strings.Split(parts[1], ",") {
			if _, err := strconv.Atoi(day); err != nil || day < "1" || day > "7" {
				return errors.New("invalid weekday")
			}
		}
	case "m":
		if len(parts) < 2 {
			return errors.New("invalid monthly repeat format")
		}
		// Дополнительная валидация для monthly
	case "y":

	default:
		return errors.New("unsupported repeat rule")
	}
	return nil
}
