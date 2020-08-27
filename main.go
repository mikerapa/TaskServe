package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Task struct {
	Title string `json:"title"`
	ID    int64  `json:"id"`
}

var taskList []Task

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/tasks", getTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	r.HandleFunc("/tasks", createTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	taskList = append(taskList, Task{Title: "Do your stuff", ID: 1})

	log.Fatal(srv.ListenAndServe())

}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(taskList)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	soughtID, _ := strconv.ParseInt(params["id"], 10, 32)

	for _, currentTask := range taskList {
		if currentTask.ID == soughtID {
			json.NewEncoder(w).Encode(currentTask)

		}
	}
	// none were found, send empty task
	json.NewEncoder(w).Encode(&Task{})
}

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newTask Task
	_ = json.NewDecoder(r.Body).Decode(&newTask)

	newTask.ID = rand.Int63n(999999)
	taskList = append(taskList, newTask)
	json.NewEncoder(w).Encode(newTask)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	soughtID, _ := strconv.ParseInt(params["id"], 10, 32)

	for i, currentTask := range taskList {
		if currentTask.ID == soughtID {
			// remove this task from the list
			taskList = append(taskList[:i], taskList[i+1])
			break
		}
	}

	var newTask Task
	_ = json.NewDecoder(r.Body).Decode(&newTask)
	newTask.ID = soughtID
	taskList = append(taskList, newTask)
	json.NewEncoder(w).Encode(newTask)

}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	soughtID, _ := strconv.ParseInt(params["id"], 10, 32)

	for i, currentTask := range taskList {
		if currentTask.ID == soughtID {
			// remove this task from the list
			taskList = append(taskList[:i], taskList[i+1:]...)
			break
		}
	}
	// return all
	json.NewEncoder(w).Encode(taskList)
}
