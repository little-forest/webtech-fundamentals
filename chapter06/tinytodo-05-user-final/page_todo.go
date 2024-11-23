package main

import (
	"html"
	"net/http"
	"strings"
)

type TodoPageData struct {
	UserId   string
	Expires  string
	ToDoList []string
}

func handleTodo(w http.ResponseWriter, r *http.Request) {
	session, err := checkSession(w, r)
	logRequest(r, session)
	if err != nil {
		return
	}
	if !isAuthenticated(w, r, session) {
		return
	}

	pageData := TodoPageData{
		UserId:   session.UserAccount.Id,
		Expires:  session.UserAccount.ExpiresText(),
		ToDoList: session.UserAccount.ToDoList.Items,
	}

	templates.ExecuteTemplate(w, "todo.html", pageData)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	session, err := checkSession(w, r)
	logRequest(r, session)
	if err != nil {
		return
	}
	if !isAuthenticated(w, r, session) {
		return
	}

	r.ParseForm()
	todo := strings.TrimSpace(html.EscapeString(r.Form.Get("todo")))
	if todo != "" {
		session.UserAccount.ToDoList.Append(todo)
	}
	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := checkSession(w, r)
	logRequest(r, session)
	if err != nil {
		return
	}

	sessionManager.RevokeSession(w, session.SessionId)
	sessionManager.StartSession(w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
