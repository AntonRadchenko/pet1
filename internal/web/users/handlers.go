package users

import (
	"context"
	"log"

	"github.com/AntonRadchenko/WebPet1/internal/userService"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type UserHandler struct {
	service *userService.UserService
}

func NewUserHandler(s *userService.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) PostUsers(_ context.Context, req PostUsersRequestObject) (PostUsersResponseObject, error) {
	params := userService.CreateUserParams{
		Email: string(req.Body.Email),
		Password: req.Body.Password,
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

func (h *UserHandler) PatchUsersId(_ context.Context, req PatchUsersIdRequestObject) (PatchUsersIdResponseObject, error) {
	params := userService.UpdateUserParams{}

    // Если поля бади не пустые, то кладем эти поля кладем в структурку 
    if req.Body.Email != nil {
        email := string(*req.Body.Email) 
        params.Email = &email
    }
    
    if req.Body.Password != nil {
        params.Password = req.Body.Password 
    }

	updatedUser, err := h.service.UpdateUser(req.Id, params)
	if err != nil {
		return nil, err
	}

	log.Printf("[PATCH] User %d updated successfully", req.Id)

	email := openapi_types.Email(updatedUser.Email)

	// маппим бизнес-модель в апи-модель
	response := PatchUsersId200JSONResponse{
		Id: &updatedUser.ID,
		Email: &email,
	}
	return response, nil
}

func (h *UserHandler) DeleteUsersId(_ context.Context, req DeleteUsersIdRequestObject) (DeleteUsersIdResponseObject, error) {
	urlID := req.Id

	if err := h.service.DeleteUser(urlID); err != nil {
		return nil, err
	}

	log.Printf("[DELETE] User %d deleted successfully", urlID)

	return DeleteUsersId204Response{}, nil
}