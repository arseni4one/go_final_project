package api

import (
	"encoding/json"
	"go_final_project/internal/database"
	"net/http"
	"time"
)

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	idStr := r.FormValue("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	// Получаем задачу
	task, err := database.GetTask(idStr)
	if err != nil {
		writeError(w, http.StatusNotFound, "Задача не найдена")
		return
	}

	// Если нет правила повторения - удаляем задачу
	if task.Repeat == "" {
		err = database.DeleteTask(idStr)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Ошибка удаления задачи")
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
		writeError(w, http.StatusBadRequest, `{"error":"`+err.Error()+`"}`)
		return
	}

	// Обновляем дату задачи
	err = database.UpdateTaskDate(idStr, nextDate)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка обновления даты задачи")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
