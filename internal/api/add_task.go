package api

import (
	"encoding/json"
	"go_final_project/internal/database"
	"log"
	"net/http"
	"time"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task database.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, `{"error":"Неверный формат JSON"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Получена задача: Date=%s, Title=%s, Repeat=%s", task.Date, task.Title, task.Repeat)

	if task.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
		return
	}

	if date.Before(nowDate) && task.Repeat != "" {
		log.Printf("Вызов NextDate: now=%s, date=%s, repeat=%s", nowDate.Format("20060102"), task.Date, task.Repeat)
		nextDate, err := NextDate(nowDate, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
		task.Date = nextDate
	} else if date.Before(nowDate) && task.Repeat == "" {
		task.Date = nowDate.Format("20060102")
	}

	log.Printf("После обработки: Date=%s", task.Date)

	id, err := database.AddTask(task)
	if err != nil {
		http.Error(w, `{"error":"Ошибка добавления задачи"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"id": id}); err != nil {
		log.Printf("Ошибка кодирования ответа: %v", err)
		http.Error(w, `{"error":"Ошибка формирования ответа"}`, http.StatusInternalServerError)
		return
	}
}
