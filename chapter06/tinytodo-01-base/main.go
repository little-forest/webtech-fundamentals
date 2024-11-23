package main

import (
	"html/template"
	"log"
	"net/http"
)

var todoList []string // <1>

func handleTodo(w http.ResponseWriter, r *http.Request) { // <5>
	t, _ := template.ParseFiles("templates/todo.html") // <6>
	t.Execute(w, todoList)                             // <7>
}

func main() {
	todoList = append(todoList, "顔を洗う", "朝食を食べる", "歯を磨く") // <2>

	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // <3>

	http.HandleFunc("/todo", handleTodo) // <4>

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("failed to start : ", err)
	}
}
