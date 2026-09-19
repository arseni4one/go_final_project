package api

import (
	"encoding/json"
	"go_final_project/internal/database"
	"net/http"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	err := database.DeleteTask(idStr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка удаления задачи")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
