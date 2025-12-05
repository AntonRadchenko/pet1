package handlers

import (
	"context"
	"log"

	"github.com/AntonRadchenko/WebPet1/openapi"
)

// слой handlers:
//   • Принимает уже РАСПАРСЕННЫЕ данные из HTTP (готовые структуры из api.gen.go)
//   • Вызывает соответствующие методы TaskService (бизнес-логика)
//   • Маппит TaskStruct (из БД) → Task (из OpenAPI)
//   • Возвращает сгенерированные типы ответов (например PostTasks201JSONResponse)
//   • НЕ работает с БД напрямую (только через сервис)
//   • НЕ парсит JSON
//   • НЕ пишет JSON
//   • НЕ устанавливает HTTP-коды
//   • НЕ управляет роутингом

func (s *Server) PostTasks(_ context.Context, req openapi.PostTasksRequestObject) (openapi.PostTasksResponseObject, error) {
	body := req.Body

	// передаем данные с тела запроса в сервис (который уже передаст их в репозиторий)
	newTask, err := s.taskService.CreateTask(*body) // передаю таску и флаг из тела запроса
	if err != nil {
		return nil, err
	}

	log.Printf("[POST] Task %d created successfully", newTask.ID)

	// маппинг в API-модель
	// openapi.Task - модель API, а service.TaskStruct - это модель бд.        <-- для себя, чтобы не путаться
	response := openapi.PostTasks201JSONResponse{
		Id:     &newTask.ID,
		Task:   &newTask.Task,
		IsDone: &newTask.IsDone,
	}
	return response, nil // отправляем клиенту ответ
}

func (s *Server) GetTasks(_ context.Context, _ openapi.GetTasksRequestObject) (openapi.GetTasksResponseObject, error) {
	// инициализируем слайс данным способом, чтобы при ошибке вернулся пустой массив, вместо null
	response := make(openapi.GetTasks200JSONResponse, 0)

	// получаем модель бд
	tasks, err := s.taskService.GetTasks()
	if err != nil {
		return nil, err
	}

	for _, t := range tasks {
		// маппинг в API-модель
		response = append(response, openapi.Task{
			Id:     &t.ID,
			Task:   &t.Task,
			IsDone: &t.IsDone,
		})

	}
	log.Printf("[GET] Returned %d tasks", len(tasks))
	return response, nil
}

func (s *Server) PatchTasksId(_ context.Context, req openapi.PatchTasksIdRequestObject) (openapi.PatchTasksIdResponseObject, error) {
	urlID := req.Id
	body := req.Body

	updatedTask, err := s.taskService.UpdateTask(urlID, *body)
	if err != nil {
		return nil, err
	}

	log.Printf("[PATCH] Task %d updated successfully", urlID)

	// маппинг в API-модель
	response := openapi.PatchTasksId200JSONResponse{
		Id:     &updatedTask.ID,
		Task:   &updatedTask.Task,
		IsDone: &updatedTask.IsDone,
	}
	return response, nil
}

func (s *Server) DeleteTasksId(_ context.Context, req openapi.DeleteTasksIdRequestObject) (openapi.DeleteTasksIdResponseObject, error) {
	urlID := req.Id

	if err := s.taskService.DeleteTask(urlID); err != nil {
		return nil, err
	}

	log.Printf("[DELETE] Task %d deleted successfully", urlID)

	return openapi.DeleteTasksId204Response{}, nil
}
