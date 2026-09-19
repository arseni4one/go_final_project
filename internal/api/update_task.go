package api

import (
	"encoding/json"
	"go_final_project/internal/database"
	"log"
	"net/http"
	"time"
)

const dateFormat = "20060102"

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task database.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, `{"error":"Неверный формат JSON"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок"}`, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if task.Date == "" {
		task.Date = nowDate.Format(dateFormat)
	}

	date, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
		return
	}

	if date.Before(nowDate) && task.Repeat != "" {
		nextDate, err := NextDate(nowDate, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
		task.Date = nextDate
	} else if date.Before(nowDate) && task.Repeat == "" {
		task.Date = nowDate.Format(dateFormat)
	}

	err = database.UpdateTask(task)
	if err != nil {
		log.Printf("Ошибка обновления задачи: %v", err)
		http.Error(w, `{"error":"Ошибка обновления задачи"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
