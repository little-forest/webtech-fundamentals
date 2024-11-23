package main

import (
	"html/template"
	"log"
	"net/http"
)

var todoList []string

func handleTodo(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/todo.html")
	t.Execute(w, todoList)
}

func handleAdd(w http.ResponseWriter, r *http.Request) { // <1>
	r.ParseForm()
	todo := r.Form.Get("todo")
	todoList = append(todoList, todo) // <3>
	handleTodo(w, r)                  // <4>
}

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/todo", handleTodo)

	http.HandleFunc("/add", handleAdd) // <2>

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("failed to start : ", err)
	}
}
