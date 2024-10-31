package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/MirekKrassilnikov/todo_list_api_sql/Database"
	"github.com/MirekKrassilnikov/todo_list_api_sql/config"
	"github.com/MirekKrassilnikov/todo_list_api_sql/server"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	if install {
		database.CreateDatabase()
	} else {
		fmt.Println("Database already exists")
	}
	srv := server.Controller{DB: db}
	// Создаем файловый сервер для директории web
	fs := http.FileServer(http.Dir(config.WebDir))
	// Настраиваем обработчик для всех запросов
	http.Handle("/", fs)
	http.HandleFunc("/api/task", srv.TaskHandler)
	http.HandleFunc("/api/tasks", srv.GetAllTasksHandler)
	http.HandleFunc("/api/nextdate", srv.ApiNextDateHandler)
	http.HandleFunc("/api/task/done", srv.MarkAsDone)
	// Запускаем сервер на указанном порту
	log.Printf("Starting server on :%s\n", config.Port)
	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
