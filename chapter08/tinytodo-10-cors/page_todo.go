package main

import (
	"net/http"
)

type TodoPageData struct {
	UserId  string
	Expires string
}

// ToDo画面を表示する。
func handleRoot(w http.ResponseWriter, r *http.Request) {
	session, err := checkSessionForPage(w, r)
	logRequest(r, session)
	if err != nil {
		return
	}
	if !isAuthenticated(w, r, session) {
		return
	}
	pageData := TodoPageData{
		UserId:  session.UserAccount.Id,
		Expires: session.UserAccount.ExpiresText(),
	}

	templates.ExecuteTemplate(w, "todo.html", pageData)
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
