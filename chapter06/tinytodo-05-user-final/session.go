package main

import (
	"net/http"
	"time"
)

const CookieNameSessionId = "sessionId"

// セッション情報を保持する構造体。
type HttpSession struct {
	// セッションID
	SessionId string
	// セッションの有効期限(時刻)
	Expires time.Time
	// Post-Redirect-Getでの遷移先に表示するデータ
	PageData any
	// ユーザアカウント情報への参照
	UserAccount *UserAccount
	// セッションの有効時間
	validityTime time.Duration
	// Cookieにsecure属性を付与するかどうか
	useSecureCookie bool
}

// 新しいセッション情報を生成する。
func NewHttpSession(sessionId string, validityTime time.Duration, useSecureCookie bool) *HttpSession {
	session := &HttpSession{
		SessionId:       sessionId,
		validityTime:    validityTime,
		PageData:        "",
		useSecureCookie: useSecureCookie,
	}
	session.Extend()
	return session
}

// 有効期限を延長する。
func (s *HttpSession) Extend() {
	s.Expires = time.Now().Add(s.validityTime)
}

// ページデータを削除する。
func (s *HttpSession) ClearPageData() {
	s.PageData = ""
}

// セッションIDをCookieへ書き込む。
func (s HttpSession) SetCookie(w http.ResponseWriter) {
	// 同一リクエスト処理内で既にCookieが書き込まれていた場合は削除する
	w.Header().Del("Set-Cookie")
	cookie := &http.Cookie{
		Name:     CookieNameSessionId,
		Value:    s.SessionId,
		Expires:  s.Expires,
		HttpOnly: true,
		Secure:   s.useSecureCookie,
	}
	http.SetCookie(w, cookie)
}
