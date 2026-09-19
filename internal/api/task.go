package api

import (
	"encoding/json"
	"go_final_project/internal/database"
	"log"
	"net/http"
)

const defaultTasksLimit = 50

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}
	tasks, err := database.Tasks(defaultTasksLimit)
	if err != nil {
		http.Error(w, `{"error":"Ошибка получения задач"}`, http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []database.Task{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks}); err != nil {
		log.Printf("error: %v", err)
	}

}
