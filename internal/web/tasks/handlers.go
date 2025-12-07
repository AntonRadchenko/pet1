package tasks

import (
	"context"
	"log"

	"github.com/AntonRadchenko/WebPet1/internal/taskService"
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

type TaskHandler struct {
	service *taskService.TaskService
}

func NewTaskHandler(s *taskService.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

func (h *TaskHandler) PostTasks(_ context.Context, req PostTasksRequestObject) (PostTasksResponseObject, error) {
	params := taskService.CreateTaskParams{
		Task: req.Body.Task,
		IsDone: req.Body.IsDone,
		UserId: req.Body.UserId,
	}

	// передаем данные с тела запроса в сервис (который уже передаст их в репозиторий)
	newTask, err := h.service.CreateTask(params) // передаю таску и флаг из тела запроса
	if err != nil {
		return nil, err
	}

	log.Printf("[POST] Task %d created successfully", newTask.ID)

	// маппим бизнес-модель в апи-модель
	response := PostTasks201JSONResponse{
		Id:     &newTask.ID,
		Task:   &newTask.Task,
		IsDone: newTask.IsDone,
		UserId: &newTask.UserId,
	}
	return response, nil // отправляем клиенту ответ
}

func (h *TaskHandler) GetTasks(_ context.Context, _ GetTasksRequestObject) (GetTasksResponseObject, error) {
	// инициализируем слайс данным способом, чтобы при ошибке вернулся пустой массив, вместо null
	response := make(GetTasks200JSONResponse, 0)

	// получаем модель бд
	tasks, err := h.service.GetTasks()
	if err != nil {
		return nil, err
	}

	for _, t := range tasks {
		// маппинг в API-модель
		response = append(response, Task{
			Id:     &t.ID,
			Task:   &t.Task,
			IsDone: t.IsDone,
			UserId: &t.UserId,
		})

	}
	log.Printf("[GET] Returned %d tasks", len(tasks))
	return response, nil
}

func (h *TaskHandler) PatchTasksId(_ context.Context, req PatchTasksIdRequestObject) (PatchTasksIdResponseObject, error) {
	params := taskService.UpdateTaskParams{}

	if req.Body.Task != nil {
		task := *req.Body.Task
		params.Task = &task
	}

	if req.Body.IsDone != nil {
		isDone := *req.Body.IsDone
		params.IsDone = &isDone
	}

	if req.Body.UserId != nil {
		userId := *req.Body.UserId
		params.UserId = &userId
	}

	updatedTask, err := h.service.UpdateTask(req.Id, params)
	if err != nil {
		return nil, err
	}

	log.Printf("[PATCH] Task %d updated successfully", req.Id)

	// маппим бизнес-модель в апи-модель
	response := PatchTasksId200JSONResponse{
		Id:     &updatedTask.ID,
		Task:   &updatedTask.Task,
		IsDone: updatedTask.IsDone,
		UserId: &updatedTask.UserId,
	}
	return response, nil
}

func (h *TaskHandler) DeleteTasksId(_ context.Context, req DeleteTasksIdRequestObject) (DeleteTasksIdResponseObject, error) {
	urlID := req.Id

	if err := h.service.DeleteTask(urlID); err != nil {
		return nil, err
	}

	log.Printf("[DELETE] Task %d deleted successfully", urlID)

	return DeleteTasksId204Response{}, nil
}
