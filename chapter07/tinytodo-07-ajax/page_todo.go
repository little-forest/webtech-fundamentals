package main

import (
	"html"
	"log"
	"net/http"
	"strings"
)

type TodoPageData struct {
	UserId   string
	Expires  string
	ToDoList []*ToDoItem
}

// ToDo画面を表示する。
func handleTodo(w http.ResponseWriter, r *http.Request) {
	session, err := checkSessionForPage(w, r)
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
	// セッション情報を取得
	session, err := checkSession(w, r)
	logRequest(r, session)
	if err != nil {
		return
	}
	// 認証チェック
	if !isAuthenticated(w, r, session) {
		return
	}

	// TodoListに新しいToDoを登録
	r.ParseForm()
	todo := strings.TrimSpace(html.EscapeString(r.Form.Get("todo")))
	if todo != "" {
		item := session.UserAccount.ToDoList.Append(todo)
		log.Printf("Todo item added. sessionId=%s itemId=%s todo=%s", session.SessionId, item.Id, item.Todo)
	}
	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}


func handleEdit(w http.ResponseWriter, r *http.Request) {
	// POSTメソッドによるリクエストであることの確認
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// セッション情報を取得
	session, err := checkSession(w, r)
	logRequest(r, session)
	if err != nil {
		return
	}
	// 認証チェック
	if !isAuthenticated(w, r, session) {
		return
	}

	// POSTパラメータを解析
	r.ParseForm() // <1>
	todoId := r.Form.Get("id")
	todo := r.Form.Get("todo")

	// ToDo項目を更新
	_, err = session.UserAccount.ToDoList.Update(todoId, todo) // <2>
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Todo item updated. sessionId=%s itemId=%s todo=%s", session.SessionId, todoId, todo)

	// レスポンスの返却
	w.WriteHeader(http.StatusOK) // <3>
}


func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := checkSessionForPage(w, r)
	logRequest(r, session)
	if err != nil {
		return
	}

	sessionManager.RevokeSession(w, session.SessionId)
	sessionManager.StartSession(w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
