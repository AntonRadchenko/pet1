package users

import (
	"context"
	"log"
	"strings"

	"github.com/AntonRadchenko/WebPet1/internal/userService"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type UserHandler struct {
	service *userService.UserService
}

func NewUserHandler(s *userService.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) PostUsers(_ context.Context, request PostUsersRequestObject) (PostUsersResponseObject, error) {
	params := userService.CreateUserParams{
		Email: string(request.Body.Email),
		Password: request.Body.Password,
	}

	newUser, err := h.service.CreateUser(params)
	if err != nil {
		return nil, err
	}

	log.Printf("[POST] User %d created successfully", newUser.ID)

	// Конвертируем string в openapi_types.Email для API ответа
	email := openapi_types.Email(newUser.Email)

	// маппим бизнес-модель в апи-модель
	response := PostUsers201JSONResponse{
		Id: &newUser.ID,
		Email: &email,
	}
	return response, nil
}

func (h *UserHandler) GetUsers(_ context.Context, _ GetUsersRequestObject) (GetUsersResponseObject, error) {
	response := make(GetUsers200JSONResponse, 0)

	users, err := h.service.GetUsers()
	if err != nil {
		return nil, err
	}	

	for _, u := range users {
		email := openapi_types.Email(u.Email)

		// маппим бизнес-модель в апи-модель
		response = append(response, User{
			Id: &u.ID,
			Email: &email,
		})
	}
	log.Printf("[GET] Returned %d users", len(users))
	return response, nil
}

func (h *UserHandler) GetUsersIdTasks(ctx context.Context, request GetUsersIdTasksRequestObject) (GetUsersIdTasksResponseObject, error) {
	tasks, err := h.service.GetTasksForUser(request.Id)
	if err != nil {
        if strings.Contains(err.Error(), "user not found") {
            return GetUsersIdTasks404Response{}, nil
        }
        return nil, err
	}

	// маппим бизнес-модель в апи-модель
	response := make(GetUsersIdTasks200JSONResponse, 0, len(tasks))

    for _, t := range tasks {
        response = append(response, Task{
            Id:     &t.ID,
            Task:   &t.Task,
            IsDone: t.IsDone,
            UserId: &t.UserId,
        })
    }

	log.Printf("[GET] Returned %d tasks for user ID %d", len(tasks), request.Id)
	return response, nil
}

func (h *UserHandler) PatchUsersId(_ context.Context, request PatchUsersIdRequestObject) (PatchUsersIdResponseObject, error) {
	params := userService.UpdateUserParams{}

    // Если поля бади не пустые, то кладем эти поля кладем в структурку 
    if request.Body.Email != nil {
        email := string(*request.Body.Email) 
        params.Email = &email
    }
    
    if request.Body.Password != nil {
        params.Password = request.Body.Password 
    }

	updatedUser, err := h.service.UpdateUser(request.Id, params)
	if err != nil {
		return nil, err
	}

	log.Printf("[PATCH] User %d updated successfully", request.Id)

	email := openapi_types.Email(updatedUser.Email)

	// маппим бизнес-модель в апи-модель
	response := PatchUsersId200JSONResponse{
		Id: &updatedUser.ID,
		Email: &email,
	}
	return response, nil
}

func (h *UserHandler) DeleteUsersId(_ context.Context, request DeleteUsersIdRequestObject) (DeleteUsersIdResponseObject, error) {
    urlID := request.Id

    if err := h.service.DeleteUser(urlID); err != nil {
        // Если "user not found" - возвращаем 404
        if strings.Contains(err.Error(), "user not found") {
            return DeleteUsersId404Response{}, nil
        }
        // Другие ошибки - 500
        return nil, err
    }

    log.Printf("[DELETE] User %d deleted successfully", urlID)
    return DeleteUsersId204Response{}, nil
}