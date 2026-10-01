package main

import (
	"encoding/json" //transformar dados de Golang para JSON
	"log"           //mensagens no terminal sobre o servidor
	"net/http"      //criar servidor e lidar com requisiçoes
	"strconv"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks = []Task{
	{ID: 1, Title: "Update project documentation", Done: false},
	{ID: 2, Title: "Deploy staging environment", Done: true},
	{ID: 3, Title: "Review pull request", Done: true},
	{ID: 4, Title: "Schedule dentist appointment", Done: false},
}

type titleReceiver struct {
	Title string `json:"title"`
}

type updatedTask struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json") //se o metodo for igual a GET, set um header
		json.NewEncoder(w).Encode(tasks)                   // crie um encoder e transforme a variavel tasks em JSON
	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json")
		var t titleReceiver
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, "invalid request body!", http.StatusBadRequest)
			return
		}
		newID := len(tasks) + 1
		newTask := Task{ID: newID, Title: t.Title, Done: false}
		tasks = append(tasks, newTask)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)
	default:
		http.Error(w, "method not allowed!", http.StatusMethodNotAllowed) //imprima method not allowed e torne o valor do erro = 405
		return
	}
}

func findTaskIndex(id int) int {
	for i := range tasks {
		if tasks[i].ID == id {
			return i
		}
	}
	return -1
}
func idTaskHandler(w http.ResponseWriter, r *http.Request) {
	idValue := r.PathValue("id")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		http.Error(w, "write a valid ID please!", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		i := findTaskIndex(id)
		if i == -1 {
			http.Error(w, "sorry, there's no task with this ID!", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks[i])
	case http.MethodPut:
		i := findTaskIndex(id)
		if i == -1 {
			http.Error(w, "sorry, there's no task with this ID!", http.StatusNotFound)
			return
		}
		var u updatedTask
		err = json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, "invalid request body!", http.StatusBadRequest)
			return
		}
		if u.Title == "" {
			http.Error(w, "write a valid Title please!", http.StatusBadRequest)
			return
		}
		tasks[i].Title = u.Title
		tasks[i].Done = u.Done
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks[i])
	/*case http.MethodDelete:
	i := findTaskIndex(id)
	if i == -1 {
		http.Error(w, "sorry, there's no task with this ID!", http.StatusNotFound)
		return
	}

	*/
	default:
		http.Error(w, "method not allowed!", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/tasks/{id}", idTaskHandler)
	log.Println("server running on http://localhost:8080/tasks")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
