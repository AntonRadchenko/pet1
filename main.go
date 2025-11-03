package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// главный объект GORM, через который идут все запросы в бд
var db *gorm.DB

// функция для инициализации подключения и работы с бд
func initDB() {
	// источник данных
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable"
	var err error

	// открываем соединение с бд
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to Database: %v", err)
	}

	// автомиграция (в бд автоматически создастся модель (таблица) на основе структуры TaskStruct)
	if err := db.AutoMigrate(&TaskStruct{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
}

// Основные методы ORM, c которыми будем работать:
// Find (найти записи в бд и заполнить переданный срез)
// Create (записать новый объект в бд)
// Update (обновить существующую запись в бд)
// Delete (удалить существующую запись из бд)

// структура хранилища тасок
type TaskStruct struct {
	ID   string `gorm:"primaryKey" json:"id"`
	Task string `json:"task"`
}

// генератор айдишек для тасок
var idCounter int

// структура тела запроса
type RequestBody struct {
	Task string `json:"task"`
}

// единый формат ошибок
type ErrorStruct struct {
	Error string `json:"error"`
}

// появилась бд => слайс больше не нужен
// var tasks = []TaskStruct{}

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
		WriteJsonError(w, http.StatusBadRequest, "task is empty")
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
	if err := db.Create(&newTask).Error; err != nil {
		WriteJsonError(w, http.StatusInternalServerError, "Could not add Task")
		return 
	}

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
	var tasks []TaskStruct

	// ищем все записи в таблице в бд и заполняем их в tasks
	if err := db.Find(&tasks).Error; err != nil { 
		WriteJsonError(w, http.StatusInternalServerError, "Could not get Tasks")
		return 
	}
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
	if len(parts) != 2 || parts[0] != "tasks" || parts[1] == "" {
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
	http.HandleFunc("/tasks", MainHandler)                    // для списка задач (GET, POST)
	http.HandleFunc("/tasks/", MainHandlerWithID)             // для конкретной задачи (PATCH, DELETE)
	if err := http.ListenAndServe(":9092", nil); err != nil { // слушаем порт 9092
		fmt.Println("Ошибка во время работы HTTP сервера: ", err)
	}
}
