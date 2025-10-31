package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Что сделали для задания 2:
// 1) Раделили структуры тела запроса и хранилища тасок (в теле просто таска, а в хранилище с айдишкой);
// 2) Изменили глобальную переменную со слойса строк на слайс структур
// 3) Изменили Post метод: таски теперь добавляются в новую глобальную переменную, и вместе с тем генерируются айдишки
// 4) Реализовали метод Patch - сравнивает айди из url с айди из хранилища с тасками и обновляет таску по этому айди.
// 5) В Get методе реализовали метод через структуру хранилища, вмето структуры тела запроса.
// 6) Также изменили Get так, чтобы он выводил вместо последней таски хранилища, таску по айди (по query параметру)
// 7) Реализовали REST-логику с помощью объединения хендлеров в основной. Теперь один путь, разные методы.

// Что осталось сделать для 2-го задания:
// 1) Сделать Delete метод

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

// хранилище тасок (новая глобальная переменная)
var idTasks = []TaskStruct{}

func PostHandler(w http.ResponseWriter, r *http.Request) {
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

	// добавляем в хранилище тасок новую таску
	idTasks = append(idTasks, newTask)
	msg := "Task created successfully!"
	fmt.Println(msg)
	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("fail to write HTTP response: ", err)
		return 
	}
}

func PatchHandler(w http.ResponseWriter, r *http.Request) {
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
	if len(idTasks) == 0 {
		msg := "No tasks to update yet!"
		fmt.Println(msg)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg + "\n"))
		return
	}
	// обновляем задачу (сопоставляем id из URL и таску из тела запроса)
	updated := false
	for i := range idTasks {
		// если id в URL совпадает с id таски в хранилище
		if idTasks[i].ID == id {
			idTasks[i].Task = request.Task // обновляем
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

	msg := "Task updated successfully"
	fmt.Println(msg)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("fail to write HTTP response: ", err)
		return 
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	// получаем id из query параметра в URL
	id := r.URL.Query().Get("id")

	// проверяем был ли вообще передан id в URL
	if id == "" {
		msg := "missing id!"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg + "\n"))
		return
	}

	// если в хранилище пока не задач, то выводить пока нечего
	if len(idTasks) == 0 {
		msg := "No tasks to print yet!"
		fmt.Println(msg)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg + "\n"))
		return
	}

	// выводим задачу по id
	printed := false
	msg := ""
	for i := range idTasks {
		if idTasks[i].ID == id {
			msg += "Hello, " + idTasks[i].Task
			printed = true
			break
		}
	}

	// если не найден такой id
	if !printed {
		msg := "Task not found!"
		fmt.Println(msg)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg + "\n"))
		return
	}

	fmt.Println(msg)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("fail to write HTTP response: ", err)
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
	if len(idTasks) == 0 {
		msg := "No tasks to print yet!"
		fmt.Println(msg)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg + "\n"))
		return 
	}

	// удаляем задачу по айди
	deleted := false
	for i := range idTasks {
		if idTasks[i].ID == id {
			idTasks = append(idTasks[:i], idTasks[i+1:]...) 
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

	msg := "Task deleted successfully!"
	fmt.Println(msg)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("fail to write HTTP response: ", err)
		return
	}
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