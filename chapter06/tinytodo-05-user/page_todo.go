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
	// セッション情報が存在するかチェック
	session, err := checkSession(w, r) // <1>
	if err != nil {
		return
	}
	// 認証済みかどうかをチェック
	if !isAuthenticated(w, r, session) { // <2>
		return
	}

	// ToDoリストの表示
	pageData := TodoPageData{ // <3>
		UserId:   session.UserAccount.Id,
		Expires:  session.UserAccount.ExpiresText(),
		ToDoList: session.UserAccount.ToDoList,
	}

	templates.ExecuteTemplate(w, "todo.html", pageData)
}


func handleAdd(w http.ResponseWriter, r *http.Request) {
	session, err := checkSession(w, r)
	if err != nil {
		return
	}
	if !isAuthenticated(w, r, session) {
		return
	}

	r.ParseForm()
	todo := strings.TrimSpace(html.EscapeString(r.Form.Get("todo")))
	if todo != "" {
		session.UserAccount.ToDoList = append(session.UserAccount.ToDoList, todo)
	}
	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := checkSession(w, r)
	if err != nil {
		return
	}

	sessionManager.RevokeSession(w, session.SessionId)
	sessionManager.StartSession(w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
