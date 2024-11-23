package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var (
	sessionManager *HttpSessionManager
	accountManager *UserAccountManager
	templates      *template.Template

	ErrMethodNotAllowed = errors.New("method not allowed")
)

func main() {
	sessionManager = NewHttpSessionManager()

	accountManager = NewUserAccountManager()

	// テンプレートを読み込む
	templates = template.Must(template.ParseGlob("templates/*.html"))

	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/create-user-account", handleCreateUserAccount)

	http.HandleFunc("/new-user-account", handleNewUserAccount)

	http.HandleFunc("/login", handleLogin)

	http.HandleFunc("/logout", handleLogout)

	http.HandleFunc("/todo", handleTodo)

	http.HandleFunc("/add", handleAdd)

	http.HandleFunc("/favicon.ico", handleNotFound)

	http.HandleFunc("/", handleRoot)

	port := getPortNumber()
	fmt.Printf("listening port : %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// ログイン画面へリダイレクトする。
func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
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

// ログイン済みかどうかを調べる。
//
// ログイン済みでない場合はログイン画面へ遷移する。
func isAuthenticated(w http.ResponseWriter, r *http.Request, session *HttpSession) bool {
	if session.UserAccount != nil {
		return true
	}

	page := LoginPageData{}
	page.ErrorMessage = "未ログインです。"
	session.PageData = page

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return false
}
