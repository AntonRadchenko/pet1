package main

import (
	"fmt"
	"net/http"
	"log"

	"github.com/AntonRadchenko/WebPet1/internal/db"
	"github.com/AntonRadchenko/WebPet1/internal/handlers"
	"github.com/AntonRadchenko/WebPet1/internal/service"
)

// 5. верхний слой (все связывается вместе)

func main() {
	db.InitDB()
	// автомиграция (в бд автоматически создастся модель (таблица) на основе структуры TaskStruct)
	if err := db.DB.AutoMigrate(&service.TaskStruct{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	repo := &service.TaskRepo{}
	svc := service.NewService(repo)
	h := handlers.NewTaskHandler(svc)

	http.HandleFunc("/tasks", h.MainHandler)
	http.HandleFunc("/tasks/", h.MainHandlerWithID)
	if err := http.ListenAndServe(":9092", nil); err != nil { // слушаем порт 9092
		fmt.Println("Ошибка во время работы HTTP сервера: ", err)
	}
}