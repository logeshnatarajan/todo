package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Todo struct {
	WORKS string `json:"work"`
	ID    string `json:"id"`
}

var Todos []Todo

func main() {
	Todos = append(Todos, Todo{WORKS: "brushing", ID: "1"}, Todo{WORKS: "bathing", ID: "2"})

	r := mux.NewRouter()
	r.HandleFunc("/get", gettodo).Methods("GET")
	r.HandleFunc("/add", addtodo).Methods("POST")
	r.HandleFunc("/delete/{id}", deletetodo).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))

}
func gettodo(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Todos)

}
func addtodo(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var todo Todo
	_ = json.NewDecoder(req.Body).Decode(&todo)
	todo.ID = strconv.Itoa((rand.Intn(100)))

	Todos = append(Todos, todo)
	json.NewEncoder(w).Encode(Todos)

}

func deletetodo(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parames := mux.Vars(req)
	for index, value := range Todos {

		if value.ID == parames["id"] {

			Todos = append(Todos[:index], Todos[index+1:]...)
			break

		}

	}
	json.NewEncoder(w).Encode(Todos)

}
