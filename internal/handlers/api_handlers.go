package handlers

import (
	"context"
	"log"

	"github.com/AntonRadchenko/WebPet1/openapi"
	"github.com/AntonRadchenko/WebPet1/internal/service"
)

// ApiHandler — слой HTTP согласно OpenAPI.
// Он реализует интерфейс StrictServerInterface, который сгенерирован oapi-codegen.

// ApiHandler:
//   • Принимает уже РАСПАРСЕННЫЕ данные из HTTP (готовые структуры из api.gen.go)
//   • Вызывает соответствующие методы TaskService (бизнес-логика)
//   • Маппит TaskStruct (из БД) → Task (из OpenAPI)
//   • Возвращает сгенерированные типы ответов (например PostTasks201JSONResponse)
//   • НЕ работает с БД напрямую (только через сервис)
//   • НЕ парсит JSON
//   • НЕ пишет JSON
//   • НЕ устанавливает HTTP-коды
//   • НЕ управляет роутингом

type ApiHandler struct {
	service *service.TaskService
}

// конструктор NewApiHandler - связывает OpenApi-слой и сервис
func NewApiHandler(s *service.TaskService) *ApiHandler {
	return &ApiHandler{service: s}
}

func (h *ApiHandler) PostTasks(ctx context.Context, req api.PostTasksRequestObject) (api.PostTasksResponseObject, error) {
	body := req.Body

	// получаем модель бд
	newTask, err := h.service.CreateTask(body.Task, body.IsDone)
	if err != nil {
		return nil, err
	}

	log.Printf("[POST] Task %d created successfully", newTask.ID)

	id := newTask.ID
	task := newTask.Task
	is_done := newTask.IsDone

	// маппинг в API-модель
	// api.Task - модель API, а TaskStruct - это модель бд.        <-- для себя, чтобы не путаться
	response := api.PostTasks201JSONResponse{
		Id: &id, 
		Task: &task, 
		IsDone: &is_done,
	}
	return response, nil

}

func (h *ApiHandler) GetTasks(ctx context.Context, req api.GetTasksRequestObject) (api.GetTasksResponseObject, error) {
	// инициализируем слайс данным способом, чтобы при ошибке вернулся пустой массив, вместо null
	response := make(api.GetTasks200JSONResponse, 0) 
	
	// получаем модель бд
	tasks, err := h.service.GetTasks()
	if err != nil {
		return nil, err
	}

	for _, t := range tasks {
		id := t.ID
		task := t.Task
		is_done := t.IsDone

		// маппинг в API-модель
		response = append(response, api.Task{
			Id: &id,
			Task: &task,
			IsDone: &is_done,
		})

	}
	log.Printf("[GET] Returned %d tasks", len(tasks))
	return response, nil
}
