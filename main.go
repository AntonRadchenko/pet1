package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Что сделали:
// 1) Установили и запустили PostgreSQL в контейнере
// 2) создали функцию подключения в бд и подключились к бд через GORM
// 3) добавили автомиграцию - GORM создал таблицу в бд по структуре TaskStruct
// 4) переписали CRUD ручки GET и POST (теперь данные берутся из реальной бд, а не из слайса)
// 5) перезаписали также ручки PATCH и DELETE, которые снала проверяют наличие таски, а потом оперируют с ней
// 6) убрали ручную генерацию id и сделали поле ID автоинкрементным (gorm:"primaryKey;autoIncrement"), 
// чтобы база сама присваивала уникальные значения.

// главный объект GORM, через который идут все запросы в бд
var db *gorm.DB

// функция для инициализации подключения и работы с бд
func initDB() {
	// источник данных (инфа о нашей бд)
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable"
	var err error

	// открываем соединение с бд (по нашим данным)
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
// Save (обновить существующую запись в бд)
// Delete (удалить существующую запись из бд)

// структура хранилища тасок
type TaskStruct struct {
	ID   uint `gorm:"primaryKey;autoIncrement"` // autoIncrement говорит GORM, что ID будет генерироваться автоматически
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

	// кладем в структуру нашу новую таску с соответствующим айди
	newTask := TaskStruct{
		Task: requestBody.Task,
	}
	// добавляем новую таску в хранилище
	if err := db.Create(&newTask).Error; err != nil {
		WriteJsonError(w, http.StatusInternalServerError, "Could not add Task")
		return
	}

	log.Printf("[POST] Task %d created successfully", newTask.ID)
	WriteJson(w, http.StatusCreated, newTask)
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

func PatchHandler(w http.ResponseWriter, r *http.Request, id string) {
	var request RequestBody // таска на которую будем менять
	// парсим json
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logError(err, "decode")
		WriteJsonError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	var task TaskStruct
	// проверяем существует ли задача в бд
	if err := db.First(&task, "id = ?", id).Error; err != nil {
		WriteJsonError(w, http.StatusNotFound, "Task not found")
		return
	}

	// обновляем
	task.Task = request.Task
	if err := db.Save(&task).Error; err != nil {
		WriteJsonError(w, http.StatusInternalServerError, "Could not update task")
		return
	}

	log.Printf("[PATCH] Task %d updated successfully!", task.ID)
	WriteJson(w, http.StatusOK, task)

	// если порядок id в JSON-ответе имеет значение, то можно обновить так:
	// db.Order("id asc").Find(&tasks) // чтобы всегда получать задачи в порядке их ID
}

func DeleteHandler(w http.ResponseWriter, r *http.Request, id string) {
	var task TaskStruct
	// проверяем существует ли задача в бд
	if err := db.First(&task, "id = ?", id).Error; err != nil {
		WriteJsonError(w, http.StatusNotFound, "Task not found")
		return
	}

	// удаляем таску
	if err := db.Delete(&task, id).Error; err != nil {
		WriteJsonError(w, http.StatusInternalServerError, "Could not delete task")
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
	initDB()
	http.HandleFunc("/tasks", MainHandler)                    // для списка задач (GET, POST)
	http.HandleFunc("/tasks/", MainHandlerWithID)             // для конкретной задачи (PATCH, DELETE)
	if err := http.ListenAndServe(":9092", nil); err != nil { // слушаем порт 9092
		fmt.Println("Ошибка во время работы HTTP сервера: ", err)
	}
}
