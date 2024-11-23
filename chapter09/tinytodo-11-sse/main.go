package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

var (
	sessionManager *HttpSessionManager
	accountManager *UserAccountManager
	templates      *template.Template

	ErrMethodNotAllowed = errors.New("method not allowed")
)

func main() {
	sessionManager = NewHttpSessionManager(getSessionSecret(), getSecureCookie())

	accountManager = NewUserAccountManager()

	// テンプレートを読み込む
	templates = template.Must(template.ParseGlob("templates/*.html"))

	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/create-user-account", handleCreateUserAccount)

	http.HandleFunc("/new-user-account", handleNewUserAccount)

	http.HandleFunc("/login", handleLogin)

	http.HandleFunc("/logout", handleLogout)

	http.HandleFunc("/favicon.ico", handleNotFound)

	// `/todos`で始まるURLをすべて処理する
	http.HandleFunc("/todos/", handleTodos)

	http.HandleFunc("/observe", handleObserve)

	http.HandleFunc("/", handleRoot)

	port := getPortNumber()
	fmt.Printf("listening port : %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

// HTTPリクエストが指定したメソッドかどうかをチェックする。
//
// 想定したメソッドでなければ、Method not allowedを返す。
func checkMethod(w http.ResponseWriter, r *http.Request, method string) error {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return ErrMethodNotAllowed
	}
	return nil
}

// エラーを出力する。
func writeInternalServerError(w http.ResponseWriter, err error) {
	msg := fmt.Sprintf("500 Internal Server Error\n\n%s", err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(msg))
}

// SECURE_COOKIE環境変数が設定されているかチェックする。
//
// yesに設定されている場合、発行するCookieにsecure属性を付与する。
func getSecureCookie() bool {
	return os.Getenv("SECURE_COOKIE") == "yes"
}

// ログイン済みかどうかを調べる。
//
// ログイン済みでない場合はログイン画面へ遷移する。
func isAuthenticated(w http.ResponseWriter, r *http.Request, session *HttpSession) bool {
	if session.UserAccount != nil {
		return true
	}

	log.Printf("not authenticated %s", session.SessionId)

	page := LoginPageData{}
	page.ErrorMessage = "未ログインです。"
	session.PageData = page

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return false
}

// APIコール時にログイン済みかどうかを調べる。
//
// ログイン済みでない場合 401 Unauthorized を返す。
func isApiAuthenticated(w http.ResponseWriter, r *http.Request, session *HttpSession) bool {
	if session.UserAccount != nil {
		return true
	}

	log.Printf("not authenticated %s", session.SessionId)

	http.Error(w, "API Requires authenticated session cookie", http.StatusUnauthorized)
	return false
}

// セッションID検証用のシークレットを返す。
//
// 環境変数 SESSION_SECRET で指定されていれば、そちらを優先する。
func getSessionSecret() uint64 {
	secretStr := os.Getenv("SESSION_SECRET")
	if secretStr != "" {
		if secretInt, err := strconv.Atoi(secretStr); err == nil {
			return uint64(secretInt)
		}
	}
	return rand.Uint64()
}

// リクエストをログに記録する。
func logRequest(r *http.Request, session *HttpSession) {
	var logMsg string
	if session != nil {
		var userId string
		if session.UserAccount != nil {
			userId = session.UserAccount.Id
		}
		logMsg = fmt.Sprintf("%s %s %s %s", r.Method, r.RequestURI, session.SessionId, userId)
	} else {
		logMsg = fmt.Sprintf("%s %s", r.Method, r.RequestURI)
	}
	log.Println(logMsg)
}
