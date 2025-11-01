package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Что сделали:
// 1) Раделили структуры тела запроса и хранилища тасок (в теле просто таска, а в хранилище с айдишкой);
// 2) Изменили глобальную переменную со слойса строк на слайс структур
// 3) Изменили Post метод: таски теперь добавляются в новую глобальную переменную, и вместе с тем генерируются айдишки
// 4) В Get методе реализовали метод через структуру хранилища, вмето структуры тела запроса.
// 5) Реализовали метод Patch - сравнивает айди из url с айди из хранилища с тасками и обновляет таску по этому айди.
// 6) Реализовали REST-логику с помощью объединения хендлеров в основной. Теперь один путь, разные методы.
// 7) Реализовали Delete метод
// 8) отправка сущности клиенту

// структура хранилища тасок
type TaskStruct struct {
	ID   string `json:"id"`
	Task string `json:"task"`
}

// структура тела запроса
type RequestBody struct {
	Task string `json:"task"`
}

// генератор айдишек для тасок
var idCounter int

// хранилище тасок (глобальная переменная)
var tasks = []TaskStruct{}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody RequestBody
	// парсим json
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Println("err: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// увеличиваем номер айдишки
	idCounter++
	// кладем в структуру нашу новую распаршенную таску с соответствующим айди
	newTask := TaskStruct{
		ID:   strconv.Itoa(idCounter),
		Task: requestBody.Task,
	}
	// добавляем новую таску в хранилище
	tasks = append(tasks, newTask)

	fmt.Println("Task created!")
	w.WriteHeader(http.StatusCreated)

	// отправляем ответ клиенту
	if err := json.NewEncoder(w).Encode(newTask); err != nil {
		fmt.Println("err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func PatchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// получаем id из query параметра в URL
	id := r.URL.Query().Get("id")

	// проверяем, а был ли вообще передан id в URL
	if id == "" {
		msg := "missing id!"
		fmt.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg + "\n"))
		return
	}

	// получаем саму таску из тела запроса (на которую будем менять старую)
	var request RequestBody
	// парсим json
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println("err: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// если задач в хранилище пока нет - то обновлять пока нечего
	if len(tasks) == 0 {
		msg := "No tasks to update yet!"
		fmt.Println(msg)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg + "\n"))
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
		msg := "Task not found!"
		fmt.Println(msg)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg + "\n"))
		return
	}

	fmt.Println("Task updated successfully!")
	w.WriteHeader(http.StatusOK)

	// отправляем ответ клиенту
	if err := json.NewEncoder(w).Encode(updatedTask); err != nil {
		fmt.Println("err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// если в хранилище пока не задач, то выводить пока нечего
	if len(tasks) == 0 {
		msg := "No tasks to print yet!"
		fmt.Println(msg)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg + "\n"))
		return
	}
	fmt.Println("Task printed")
	w.WriteHeader(http.StatusOK)
	// отправляем ответ клиенту
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		fmt.Println("err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// получаем id из query параметра в URL
	id := r.URL.Query().Get("id")

	// проверяем, был ли вообще передан id в URL
	if id == "" {
		msg := "missing id!"
		fmt.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg + "\n"))
		return
	}

	// если в хранилище пока нет задач, то удалять пока нечего
	if len(tasks) == 0 {
		msg := "No tasks to print yet!"
		fmt.Println(msg)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg + "\n"))
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
		msg := "Task not found!"
		fmt.Println(msg)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg + "\n"))
		return
	}

	msg := "Task deleted!"
	fmt.Println(msg)
	w.WriteHeader(http.StatusNoContent)
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetHandler(w, r)
	case http.MethodPost:
		PostHandler(w, r)
	case http.MethodPatch:
		PatchHandler(w, r)
	case http.MethodDelete:
		DeleteHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed) // ошибка метода
	}
}

func main() {
	http.HandleFunc("/task", MainHandler)
	if err := http.ListenAndServe(":9092", nil); err != nil { // слушаем порт 9092
		fmt.Println("Ошибка во время работы HTTP сервера: ", err)
	}
}
