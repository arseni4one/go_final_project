package api

import (
	"encoding/json"
	"go_final_project/internal/database"
	"net/http"
	"time"
)

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	// Получаем задачу
	task, err := database.GetTask(idStr)
	if err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	// Если нет правила повторения - удаляем задачу
	if task.Repeat == "" {
		err = database.DeleteTask(idStr)
		if err != nil {
			http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}

	// Если есть правило повторения - вычисляем следующую дату
	now := time.Now()
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	nextDate, err := NextDate(nowDate, task.Date, task.Repeat)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Обновляем дату задачи
	err = database.UpdateTaskDate(idStr, nextDate)
	if err != nil {
		http.Error(w, `{"error":"Ошибка обновления даты задачи"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
