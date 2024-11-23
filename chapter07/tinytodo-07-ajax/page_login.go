package main

import (
	"log"
	"net/http"
	"time"
)

const CookieNameUserId = "tinyToDoUserId"

type LoginPageData struct {
	UserId       string
	ErrorMessage string
}

// ログインに関するリクエスト処理
func handleLogin(w http.ResponseWriter, r *http.Request) {
	session, err := ensureSession(w, r)
	logRequest(r, session)
	if err != nil {
		return
	}

	switch r.Method {
	// GETリクエスト:ログイン画面の表示
	case http.MethodGet:
		showLogin(w, r, session)
		return

	// POSTリクエスト:ログイン処理
	case http.MethodPost:
		login(w, r, session)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// ログイン画面を表示する。
func showLogin(w http.ResponseWriter, r *http.Request, session *HttpSession) {
	// すでにログイン中の場合はtodo画面へ遷移する
	if session.UserAccount != nil {
		log.Printf("already login : %s\n", session.UserAccount.Id)
		http.Redirect(w, r, "/todo", http.StatusSeeOther)
		return
	}

	var pageData LoginPageData
	if p, ok := session.PageData.(LoginPageData); ok {
		pageData = p
	} else {
		pageData = LoginPageData{}
	}
	pageData.UserId = getUserIdFromCookie(r)

	templates.ExecuteTemplate(w, "login.html", pageData)
	session.ClearPageData()
}

// ログイン処理を行う。
func login(w http.ResponseWriter, r *http.Request, session *HttpSession) {
	// POSTパラメータを取得
	r.ParseForm()
	userId := r.Form.Get("userId")
	password := r.Form.Get("password")

	// セッションを新しく作る
	sessionManager.RevokeSession(nil, session.SessionId)
	session, err := sessionManager.StartSession(w)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	// 認証処理
	log.Printf("login attempt : %s\n", userId)
	account, err := accountManager.Authenticate(userId, password)
	if err != nil {
		if err == ErrAccountExpired {
			log.Printf("account expired : %s\n", userId)
			session.PageData = LoginPageData{
				ErrorMessage: "アカウントの有効期限が切れています",
			}
		} else {
			log.Printf("login failed : %s\n", userId)
			session.PageData = LoginPageData{
				ErrorMessage: "ユーザIDまたはパスワードが違います",
			}
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// ログイン成功
	session.UserAccount = account
	saveUserIdToCookie(w, account)

	log.Printf("login success : %s\n", account.Id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

// アカウントIDをCookieへ保存する。
func saveUserIdToCookie(w http.ResponseWriter, account *UserAccount) {
	cookie := &http.Cookie{
		Name:     CookieNameUserId,
		Value:    account.Id,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

// Cookieに保存されたアカウントIDを取得する。
func getUserIdFromCookie(r *http.Request) string {
	c, err := r.Cookie(CookieNameUserId)
	if err != nil {
		return ""
	}
	return c.Value
}
