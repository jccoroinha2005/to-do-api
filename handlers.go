package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		rows, err := db.Query(r.Context(), "SELECT id, title, done FROM tasks ORDER BY id")
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		list := []Task{}
		for rows.Next() {
			var t Task
			if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
				http.Error(w, "database error", http.StatusInternalServerError)
				return
			}
			list = append(list, t)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)

	case http.MethodPost:
		var t titleReceiver
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, "invalid request body: expected JSON with a title", http.StatusBadRequest)
			return
		}
		if t.Title == "" {
			http.Error(w, "title cannot be empty", http.StatusBadRequest)
			return
		}

		var newTask Task
		err = db.QueryRow(r.Context(),
			"INSERT INTO tasks (title) VALUES ($1) RETURNING id, title, done",
			t.Title).Scan(&newTask.ID, &newTask.Title, &newTask.Done)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed) //imprima method not allowed e torne o valor do erro = 405
	}
}

func idTaskHandler(w http.ResponseWriter, r *http.Request) {
	idValue := r.PathValue("id")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		http.Error(w, "invalid task ID: it must be a number", http.StatusBadRequest)
		return
	}

	switch r.Method {

	case http.MethodGet:
		var t Task
		err = db.QueryRow(r.Context(),
			"SELECT id, title, done FROM tasks WHERE id = $1", id).Scan(&t.ID, &t.Title, &t.Done)
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)

	case http.MethodPut:
		var u updatedTask
		err = json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, "invalid request body: expected JSON with a title and done", http.StatusBadRequest)
			return
		}
		if u.Title == "" {
			http.Error(w, "title cannot be empty", http.StatusBadRequest)
			return
		}

		var t Task
		err = db.QueryRow(r.Context(),
			"UPDATE tasks SET title = $1, done = $2 WHERE id = $3 RETURNING id, title, done",
			u.Title, u.Done, id).Scan(&t.ID, &t.Title, &t.Done)
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)

	case http.MethodDelete:
		tag, err := db.Exec(r.Context(), "DELETE FROM tasks WHERE id = $1", id)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		if tag.RowsAffected() == 0 {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
