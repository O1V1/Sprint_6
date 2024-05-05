package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// Обработчик 1 getTasks должен вернуть все задачи, которые хранятся в мапе.
func getTasks(w http.ResponseWriter, r *http.Request) {
	//сериаллизируем  данные из tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//записываем тип контента в заголовок
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//записываем в тело ответа сериализованные в формат json данные
	w.Write(resp)
}

// Обработчик 2 должен принимать задачу в теле запроса и сохранять ее в мапе
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// Обработчик 3 getTask должен вернуть задачу с указанным в запросе пути ID, если такая есть в мапе
func getTask(w http.ResponseWriter, r *http.Request) {
	//возвращает значение параметра изURL
	id := chi.URLParam(r, "id")
	//проверка,есть ли этот id
	task, ok := tasks[id]

	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Обработчик 4 должен удалить задачу из мапы по её ID
func delTask(w http.ResponseWriter, r *http.Request) {
	//возвращает значение параметра изURL
	id := chi.URLParam(r, "id")

	//проверка,есть ли этот id
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//w.Write(resp)
}

func main() {
	// создание нового роутера
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// регистрируем endpoint `/tasks` с методом GET, для которого создан getTasks
	r.Get("/tasks", getTasks)
	//регистрируем endpoint `/tasks` с методом POST, для которого создан postTask
	r.Post("/tasks", postTask)
	//регистрируем endpoint `/tasks/{id}` с методом GET, для которого создан getTask
	r.Get("/tasks/{id}", getTask)
	//регистрируем endpoint `/tasks/{id}` с методом DELETE, для которого создан
	r.Delete("/tasks/{id}", delTask)

	//запускается сервер
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
