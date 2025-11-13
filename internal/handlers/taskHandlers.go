package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/AntonRadchenko/WebPet1/internal/service"
)

// 4. taskHandlers-слой (рот)

// TaskHandler — слой HTTP (есть ссылка на TaskService) --
// -- то есть принимает HTTP запросы и вызывает нужные методы сервиса (+ отправляет ответ клиенту)
// Этот слой ничего не знает о БД — он работает только через TaskService.

type TaskHandler struct {
	service *service.TaskService
}

// конструктор NewTaskHandler - связывает HTTP и сервис
func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

// универсальная отправка JSON ответа клиенту
func WriteJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		fmt.Println("err: ", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}

// универсальная отправка ошибки в JSON клиенту
func WriteJsonError(w http.ResponseWriter, status int, msg string) {
	WriteJson(w, status, service.ErrorStruct{Error: msg})
}

// логирование ошибок (в терминал)
func logError(err error, context string) {
	if err != nil {
		log.Printf("[%s] error: %v", context, err)
	}
}

// ----- CRUD ------

func (h *TaskHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	var req service.RequestBody // таска которую будем записывать 

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logError(err, "decode")
		WriteJsonError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	
	newTask, err := h.service.CreateTask(req.Task, req.IsDone)
	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, err.Error())
		return 
	}
	log.Printf("[POST] Task %d created successfully", newTask.ID)
	WriteJson(w, http.StatusCreated, newTask)
}

func (h *TaskHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetTasks()
	if err != nil {
		WriteJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("[GET] Tasks printed")
	WriteJson(w, http.StatusOK, tasks)
}

func (h *TaskHandler) PatchHandler(w http.ResponseWriter, r *http.Request, id string) {
	var req service.RequestBody // таска на которую будем менять

	uid, err := strconv.ParseUint(id, 10, 64) // превращаем строковый id в uint64 (а его потом превратим в uint, как в бд)
	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, "invalid JSON")
		return 
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {	
		logError(err, "decode")
		WriteJsonError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	task, err := h.service.UpdateTask(uint(uid), req.Task, req.IsDone)
	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, err.Error())
		return 
	}

	log.Printf("[PATCH] Task %d updated successfully!", task.ID)
	WriteJson(w, http.StatusOK, task)
}

func (h *TaskHandler) DeleteHandler(w http.ResponseWriter, r *http.Request, id string) {
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		WriteJsonError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if err := h.service.DeleteTask(uint(uid)); err != nil {
		WriteJsonError(w, http.StatusNotFound, err.Error())
		return
	}
	log.Printf("[DELETE] Task %s deleted successfully", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) MainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetHandler(w, r)
	case http.MethodPost:
		h.PostHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed) // ошибка метода
	}
}

func (h *TaskHandler) MainHandlerWithID(w http.ResponseWriter, r *http.Request) {
	// оставляем в пути только ID (чтобы к нему обратиться)
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 || parts[0] != "tasks" || parts[1] == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id := parts[1]

	switch r.Method {
	case http.MethodPatch:
		h.PatchHandler(w, r, id)
	case http.MethodDelete:
		h.DeleteHandler(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed) // ошибка метода
	}
}