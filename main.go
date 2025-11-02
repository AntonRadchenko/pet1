package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Что сделали:
// 1) универсальная отправка JSON ответа клиенту, вместо текста (в том числе ошибки)
// 2) логирование ошибок (в терминал рзработчика)

// структура хранилища тасок
type TaskStruct struct {
	ID   string `json:"id"`
	Task string `json:"task"`
}

// структура тела запроса
type RequestBody struct {
	Task string `json:"task"`
}

// единый формат ошибок
type ErrorStruct struct {
	Error string `json:"error"`
}

// генератор айдишек для тасок
var idCounter int

// хранилище тасок (глобальная переменная)
var tasks = []TaskStruct{}

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
	WriteJson(w, status, ErrorStruct{Error: msg})
}

// логирование ошибок (в терминал)
func logError(err error, context string) {
	if err != nil {
		log.Printf("[%s] error: %v", context, err)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody
	// парсим json
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		logError(err, "decode")
		WriteJsonError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// проверяем что переданная таска не пустая
	if strings.TrimSpace(requestBody.Task) == "" {
		WriteJsonError(w, http.StatusInternalServerError, "task is empty")
		return
	}

	// увеличиваем номер айдишки
	idCounter++
	// кладем в структуру нашу новую таску с соответствующим айди
	newTask := TaskStruct{
		ID:   strconv.Itoa(idCounter),
		Task: requestBody.Task,
	}
	// добавляем новую таску в хранилище
	tasks = append(tasks, newTask)

	log.Printf("[POST] Task %s created successfully", newTask.ID)
	WriteJson(w, http.StatusCreated, newTask)
}

func PatchHandler(w http.ResponseWriter, r *http.Request, id string) {
	var request RequestBody // таска на которую будем менять
	// парсим json
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logError(err, "decode")
		WriteJsonError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	// если задач в хранилище пока нет - то обновлять пока нечего
	if len(tasks) == 0 {
		WriteJsonError(w, http.StatusNotFound, "No tasks to update")
		return
	}
	// обновляем задачу
	updated := false           // флаг обновленной задачи
	var updatedTask TaskStruct // переменная для хранения обновленной задачи
	for i := range tasks {
		// если id в URL совпадает с id таски в хранилище
		if tasks[i].ID == id {
			tasks[i].Task = request.Task // обновляем
			updatedTask = tasks[i]       // сохраняем копию обновленной таски для ответа клиенту
			updated = true
			break
		}
	}
	// если не найден такой id
	if !updated {
		WriteJsonError(w, http.StatusNotFound, "Task not found")
		return
	}

	log.Printf("[PATCH] Task %s updated successfully!", updatedTask.ID)
	WriteJson(w, http.StatusOK, updatedTask)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GET] Tasks printed")
	WriteJson(w, http.StatusOK, tasks)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request, id string) {
	// если в хранилище пока нет задач, то удалять пока нечего
	if len(tasks) == 0 {
		WriteJsonError(w, http.StatusNotFound, "No tasks to delete")
		return
	}
	// удаляем задачу по айди
	deleted := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			deleted = true
			break
		}
	}

	// если не найден такой id
	if !deleted {
		WriteJsonError(w, http.StatusNotFound, "Task not found")
		return
	}

	log.Printf("[DELETE] Task %s deleted successfully", id)
	w.WriteHeader(http.StatusNoContent)
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetHandler(w, r)
	case http.MethodPost:
		PostHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed) // ошибка метода
	}
}

func MainHandlerWithID(w http.ResponseWriter, r *http.Request) {
	// оставляем в пути только ID (чтобы к нему обратиться)
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 || parts[0] != "task" || parts[1] == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id := parts[1]
	// проверяем, был ли вообще передан id в URL
	if id == "" {
		msg := "missing id!"
		fmt.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg + "\n"))
		return
	}

	switch r.Method {
	case http.MethodPatch:
		PatchHandler(w, r, id)
	case http.MethodDelete:
		DeleteHandler(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed) // ошибка метода
	}
}

func main() {
	http.HandleFunc("/task", MainHandler)                     // для списка задач (GET, POST)
	http.HandleFunc("/task/", MainHandlerWithID)              // для конкретной задачи (PATCH, DELETE)
	if err := http.ListenAndServe(":9092", nil); err != nil { // слушаем порт 9092
		fmt.Println("Ошибка во время работы HTTP сервера: ", err)
	}
}
