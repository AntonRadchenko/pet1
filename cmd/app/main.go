package main

import (
	"log"
	"net/http"

	"github.com/AntonRadchenko/WebPet1/openapi"
	"github.com/AntonRadchenko/WebPet1/internal/db"
	"github.com/AntonRadchenko/WebPet1/internal/handlers"
	"github.com/AntonRadchenko/WebPet1/internal/service"
)

// 5. верхний слой (все связывается вместе)

func main() {
	// инициализируем бд
	db.InitDB()

	// собираем слои: repo -> service -> api handler
	repo := &service.TaskRepo{}
	svc := service.NewService(repo)
	apiHandler := handlers.NewApiHandler(svc)

	// оборачиваем API-хендлер в strict-server
	strictHandler := openapi.NewStrictHandler(apiHandler, nil)

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
