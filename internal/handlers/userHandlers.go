package handlers

import (
	"context"
	"log"

	"github.com/AntonRadchenko/WebPet1/openapi"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *Server) PostUsers(_ context.Context, req openapi.PostUsersRequestObject) (openapi.PostUsersResponseObject, error) {
	body := req.Body

	newUser, err := s.userService.CreateUser(*body)
	if err != nil {
		return nil, err
	}

	log.Printf("[POST] User %d created successfully", newUser.ID)

	// преобразуем тип Email к нужному типу (как в api - модели)
	email := openapi_types.Email(newUser.Email)

	// маппинг
	response := openapi.PostUsers201JSONResponse{
		Id: &newUser.ID,
		Email: &email,
	}
	return response, nil
}

func (s *Server) GetUsers(_ context.Context, _ openapi.GetUsersRequestObject) (openapi.GetUsersResponseObject, error) {
	response := make(openapi.GetUsers200JSONResponse, 0)

	users, err := s.userService.GetUsers()
	if err != nil {
		return nil, err
	}	

	for _, u := range users {
		email := openapi_types.Email(u.Email)
		// маппинг
		response = append(response, openapi.User{
			Id: &u.ID,
			Email: &email,
		})
	}
	log.Printf("[GET] Returned %d users", len(users))
	return response, nil
}

func (s *Server) PatchUsersId(_ context.Context, req openapi.PatchUsersIdRequestObject) (openapi.PatchUsersIdResponseObject, error) {
	urlID := req.Id
	body := req.Body

	updatedUser, err := s.userService.UpdateUser(urlID, *body)
	if err != nil {
		return nil, err
	}

	log.Printf("[PATCH] User %d updated successfully", urlID)

	email := openapi_types.Email(updatedUser.Email)
	// маппинг
	response := openapi.PatchUsersId200JSONResponse{
		Id: &updatedUser.ID,
		Email: &email,
	}
	return response, nil
}

func (s *Server) DeleteUsersId(_ context.Context, req openapi.DeleteUsersIdRequestObject) (openapi.DeleteUsersIdResponseObject, error) {
	urlID := req.Id

	if err := s.userService.DeleteUser(urlID); err != nil {
		return nil, err
	}

	log.Printf("[DELETE] User %d deleted successfully", urlID)

	return openapi.DeleteUsersId204Response{}, nil
}