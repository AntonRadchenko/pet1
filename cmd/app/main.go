package main

import (
	"log"
	"net/http"

	"github.com/AntonRadchenko/WebPet1/internal/db"
	"github.com/AntonRadchenko/WebPet1/internal/taskService"
	"github.com/AntonRadchenko/WebPet1/internal/userService"
    "github.com/AntonRadchenko/WebPet1/internal/web/tasks"
    "github.com/AntonRadchenko/WebPet1/internal/web/users" // users пакет // users API
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

	// создаём handlers (TaskHandler и UserHandler)
	taskHandler := tasks.NewTaskHandler(tasksService)
	userHandler := users.NewUserHandler(usersSevice)

	// оборачиваем API-хендлеры в strict-server 
    strictTaskHandler := tasks.NewStrictHandler(taskHandler, nil)
    strictUserHandler := users.NewStrictHandler(userHandler, nil)

	// создаём наш router
	mux := http.NewServeMux()

	// регистрируем OpenAPI маршруты в mux
	tasks.HandlerFromMux(strictTaskHandler, mux)
	users.HandlerFromMux(strictUserHandler, mux)

	// запускаем сервер
	log.Println("Server is running on :9092")
	if err := http.ListenAndServe(":9092", mux); err != nil { // слушаем порт 9092
		log.Fatal(err)
	}
}
