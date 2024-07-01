package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/gorilla/mux"
)

type ToDo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

var todos []ToDo
var idCounter int32

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo ToDo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo.ID = int(atomic.AddInt32(&idCounter, 1))
	todos = append(todos, todo)
	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	for index, item := range todos {
		if item.ID == id {
			var updatedTodo ToDo
			if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
				http.Error(w, "Error decoding request body", http.StatusBadRequest)
				return
			}

			// Validate the updatedTodo, e.g., check if the name field is not empty
			if updatedTodo.Name == "" {
				http.Error(w, "The name field cannot be empty", http.StatusBadRequest)
				return
			}

			// Preserve the ID and any other fields that should not change
			updatedTodo.ID = item.ID

			// Update the todo item in the slice
			todos[index] = updatedTodo

			// Respond with the updated todo item
			json.NewEncoder(w).Encode(updatedTodo)
			return
		}
	}

	// If we didn't find the todo item, return a 404 Not Found
	http.NotFound(w, r)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for index, item := range todos {
		if item.ID == id {
			todos = append(todos[:index], todos[index+1:]...)
			break
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/todos", getTodos).Methods("GET")
	r.HandleFunc("/todos", createTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
