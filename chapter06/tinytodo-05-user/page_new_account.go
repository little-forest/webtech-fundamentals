package main

import (
	"net/http"
)

type NewUserAccountPageData struct {
	UserId       string
	Password     string
	Expires      string
	ErrorMessage string
}

// 生成したアカウント情報の表示画面のリクエスト処理
func handleNewUserAccount(w http.ResponseWriter, r *http.Request) {
	if err := checkMethod(w, r, http.MethodGet); err != nil {
		return
	}
	session, err := checkSession(w, r)
	if err != nil {
		return
	}
	if pageData, ok := session.PageData.(NewUserAccountPageData); ok {
		err = templates.ExecuteTemplate(w, "new-user-account.html", pageData)
		if err != nil {
			writeInternalServerError(w, err)
		}
		session.ClearPageData()
	} else {
		// 不正な画面遷移
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
