package main

import (
	"log"
	"net/http"
)

type LoginPageData struct {
	UserId       string
	ErrorMessage string
}

// ログインに関するリクエスト処理
func handleLogin(w http.ResponseWriter, r *http.Request) {
	session, err := ensureSession(w, r) // <1>
	if err != nil {
		return
	}

	switch r.Method {
	// GETリクエスト:ログイン画面の表示
	case http.MethodGet:
		showLogin(w, r, session) // <2>
		return

	// POSTリクエスト:ログイン処理 // <3>
	case http.MethodPost:
		login(w, r, session)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed) // <4>
		return
	}
}


// ログイン画面を表示する。
func showLogin(w http.ResponseWriter, r *http.Request, session *HttpSession) {
	var pageData LoginPageData
	if p, ok := session.PageData.(LoginPageData); ok {
		pageData = p
	} else {
		pageData = LoginPageData{}
	}

	templates.ExecuteTemplate(w, "login.html", pageData)
	session.ClearPageData()
}


// ログイン処理を行う。
func login(w http.ResponseWriter, r *http.Request, session *HttpSession) {
	// POSTパラメータを取得 <1>
	r.ParseForm()
	userId := r.Form.Get("userId")
	password := r.Form.Get("password")

	// 認証処理
	log.Printf("login attempt : %s\n", userId)
	account, err := accountManager.Authenticate(userId, password) // <2>
	if err != nil {                                               // <3>
		log.Printf("login failed : %s\n", userId)
		session.PageData = LoginPageData{ // <6>
			ErrorMessage: "ユーザIDまたはパスワードが違います",
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther) // <7>
		return
	}

	// ログイン成功
	session.UserAccount = account // <4>

	log.Printf("login success : %s\n", account.Id)
	http.Redirect(w, r, "/todo", http.StatusSeeOther) // <5>
	return
}

