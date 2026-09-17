package main

import (
	"log"
	"net/http"
	"os"

	"go_final_project/internal/api"
	"go_final_project/internal/database"
)

func main() {
	port := "7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	dbFile := "./scheduler.db"
	if envDbFile := os.Getenv("TODO_DBFILE"); envDbFile != "" {
		dbFile = envDbFile
	}

	err := database.Init(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("База данных инициализирована")

	http.HandleFunc("/api/nextdate", api.NextDateHandler)
	http.HandleFunc("/api/task/done", api.DoneTaskHandler)
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.GetTaskHandler(w, r)
		case http.MethodPost:
			api.AddTaskHandler(w, r)
		case http.MethodPut:
			api.UpdateTaskHandler(w, r)
		case http.MethodDelete:
			api.DeleteTaskHandler(w, r)
		default:
			http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/tasks", api.TasksHandler)

	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Printf("Сервер запущен на порту %s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
