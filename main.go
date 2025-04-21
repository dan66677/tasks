package main

import (
	"log"
	"net/http"
	"time"

	"tz/internal/api"
	"tz/internal/repository"
	"tz/internal/service"
	"tz/internal/worker"
)

func main() {

	taskRepo := repository.NewTaskRepository()
	taskWorker := worker.NewTaskWorker(taskRepo)
	taskService := service.NewTaskService(taskRepo, taskWorker) 

	go taskWorker.Start()

	router := api.NewRouter(taskService)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server on :8080")
	log.Fatal(server.ListenAndServe())
}
