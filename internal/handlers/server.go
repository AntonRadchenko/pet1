package handlers

import (
	"github.com/AntonRadchenko/WebPet1/internal/taskService"
	"github.com/AntonRadchenko/WebPet1/internal/userService"
)

// Server - структура, реализующая интерфейс StrictServerInterface.
// Связывает HTTP-обработчики (OpenAPI слой) с бизнес-логикой (сервисы).

// p.s: Сделали так, потому что сгенерированный StrictServerInterface 
// требует все 8 методов в одной структуре (4 tasks + 4 users)
type Server struct {
	taskService *taskService.TaskService // сервис для работы с тасками
	userService *userService.UserService // сервис для работы с юзерами
}

func NewServer(ts *taskService.TaskService, us *userService.UserService) *Server {
	return &Server{	
		taskService: ts,
		userService: us,
	}
}

