package main

import (
	"html"
	"html/template"
	"net/http"
	"strings"
)

// セッションIDをキーとしてToDoリストを保持するマップ
var todoLists = make(map[string][]string) // <1>

// セッションIDに紐付くToDoリストを取得する。
func getTodoList(sessionId string) []string { // <2>
	todos, ok := todoLists[sessionId]
	if !ok {
		todos = []string{} // <3>
		todoLists[sessionId] = todos
	}
	return todos
}

// ToDoリストを返す。
func handleTodo(w http.ResponseWriter, r *http.Request) {
	sessionId, err := ensureSession(w, r) // <4>
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	todos := getTodoList(sessionId) // <5>

	t, _ := template.ParseFiles("templates/todo.html")
	t.Execute(w, todos)
}

// セッション上のToDoリストにToDoを追加する。
func handleAdd(w http.ResponseWriter, r *http.Request) {
	sessionId, err := ensureSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	todos := getTodoList(sessionId)

	r.ParseForm()
	todo := strings.TrimSpace(html.EscapeString(r.Form.Get("todo")))
	if todo != "" {
		todoLists[sessionId] = append(todos, todo)
	}
	http.Redirect(w, r, "/todo", 303)
}
