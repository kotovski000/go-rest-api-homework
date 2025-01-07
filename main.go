package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

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
func getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasksJSON, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tasksJSON)
}
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Валидация полученных данных
	if newTask.ID == "" || newTask.Description == "" {
		http.Error(w, "ID и Description обязательны для создания задачи", http.StatusBadRequest)
		return
	}

	if _, exists := tasks[newTask.ID]; exists {
		http.Error(w, "Задача с таким ID уже существует", http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func getTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	task, ok := tasks[taskID]
	if !ok {
		http.Error(w, "Task not found", http.StatusBadRequest)
		return
	}
	taskJSON, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(taskJSON)
}

func deleteTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	_, ok := tasks[taskID]
	if !ok {
		http.Error(w, "Task not found", http.StatusBadRequest)
		return
	}

	delete(tasks, taskID)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики

	// Обработчик для получения всех задач
	r.Get("/tasks", getAllTasksHandler)
	// Обработчик для отправки задачи на сервер
	r.Post("/tasks", createTaskHandler)
	// Обработчик для получения задачи по ID
	r.Get("/task/{id}", getTaskByID)
	// Обработчик удаления задачи по ID
	r.Delete("/task/{id}", deleteTaskByIDHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
