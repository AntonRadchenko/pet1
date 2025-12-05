package main

import (
	"log"
	"net/http"

	"github.com/AntonRadchenko/WebPet1/internal/db"
	"github.com/AntonRadchenko/WebPet1/internal/handlers"
	"github.com/AntonRadchenko/WebPet1/internal/taskService"
	"github.com/AntonRadchenko/WebPet1/internal/userService"
	"github.com/AntonRadchenko/WebPet1/openapi"
)

// 5. верхний слой (все связывается вместе)

func main() {
	// инициализируем бд
	db.InitDB()

	// собираем слои
	// tasks-слои (repo -> service)
	tasksRepo := &taskService.TaskRepo{}
	tasksService := taskService.NewTaskService(tasksRepo)

	// users-слои (repo -> service)
	usersRepo := &userService.UserRepo{}
	usersSevice := userService.NewUserService(usersRepo)

	// создаем общий сервер для обоих сервисов (tasks и users)
	server := handlers.NewServer(tasksService, usersSevice)

	// оборачиваем API-хендлер в strict-server (server реализует StrictServerInterface)
	strictHandler := openapi.NewStrictHandler(server, nil)	

	// создаём наш router
	mux := http.NewServeMux()

	// регистрируем OpenAPI маршруты в mux
	openapi.HandlerFromMux(strictHandler, mux)

	// запускаем сервер
	log.Println("Server is running on :9092")
	if err := http.ListenAndServe(":9092", mux); err != nil { // слушаем порт 9092
		log.Fatal(err)
	}
}
