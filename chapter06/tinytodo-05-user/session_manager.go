package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	ErrSessionExpired   = errors.New("session expired")
	ErrSessionNotFound  = errors.New("session not found")
	ErrInvalidSessionId = errors.New("invalid session id")
)

// セッションを管理する構造体
type HttpSessionManager struct {
	// セッションIDをキーとしてセッション情報を保持するマップ
	sessions map[string]*HttpSession // <1>
}

func NewHttpSessionManager() *HttpSessionManager {
	mgr := &HttpSessionManager{
		sessions: make(map[string]*HttpSession),
	}
	return mgr
}



// セッションを開始してCokkieにセッションIDを書き込む。
func (m *HttpSessionManager) StartSession(w http.ResponseWriter) (*HttpSession, error) {
	// 新しいセッションIDを生成する
	sessionId, err := m.makeSessionId()
	if err != nil {
		return nil, err
	}

	// セッション情報を生成する
	log.Printf("start session : %s", sessionId)
	session := NewHttpSession(sessionId, 10*time.Minute)
	m.sessions[sessionId] = session // <1>
	session.SetCookie(w)            // <2>

	return session, nil
}

// セッションIDを生成する。
func (m *HttpSessionManager) makeSessionId() (string, error) {
	randBytes := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, randBytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randBytes), nil
}


// セッションを削除する
func (m *HttpSessionManager) RevokeSession(w http.ResponseWriter, sessionId string) {
	// セッション情報を削除
	delete(m.sessions, sessionId)
	log.Printf("session revoked : %s", sessionId)

	if w == nil {
		return
	}
	cookie := &http.Cookie{
		Name:    CookieNameSessionId,
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
	}
	http.SetCookie(w, cookie)
}


// セッションが存在するかチェックする。
//
// セッションが存在しなければ、ログイン画面へリダイレクトさせる。
func checkSession(w http.ResponseWriter, r *http.Request) (*HttpSession, error) {
	// CookieのセッションIDに紐付くセッション情報を取得する
	session, err := sessionManager.GetValidSession(r) // <1>
	if err == nil {
		// セッション情報が取得できたら終了
		return session, nil // <2>
	}
	orgErr := err

	// セッションが有効期限切れまたは不正な場合、セッションを作り直す
	log.Printf("session check failed : %s", err.Error())
	session, err = sessionManager.StartSession(w) // <3>
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	// Refererヘッダの有無で他の画面からの遷移かどうかを判定
	// アプリケーショントップのURLに直接アクセスした際は、セッションが存在しないのが
	// 正常であるため、エラーを表示しないための措置
	if r.Referer() != "" {
		page := LoginPageData{}
		page.ErrorMessage = "セッションが不正です。"
		session.PageData = page
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther) // <4>
	return nil, orgErr
}

// Cookieから有効なセッションを取得する。
//
// CookieにセッションIDがなければ ErrSessionNotFound を返す。
// CookieにセッションIDが存在すれば、セッションIDに紐付く HttpSession を返す。
// セッションIDが不正な場合や、セッションの有効期限が切れている場合は、エラーを返す。
func (m *HttpSessionManager) GetValidSession(r *http.Request) (*HttpSession, error) {
	c, err := r.Cookie(CookieNameSessionId) // <5>
	// CookieにセッションIDが存在しない場合
	if err == http.ErrNoCookie {
		return nil, ErrSessionNotFound // <6>
	}
	// CookieにセッションIDが存在する場合
	if err == nil {
		// セッションを取得して返す
		sessionId := c.Value
		session, err := m.getSession(sessionId) // <7>
		return session, err
	}
	return nil, err
}

// セッションIDに紐付くセッション情報を返す。
func (m *HttpSessionManager) getSession(sessionId string) (*HttpSession, error) {
	if session, exists := m.sessions[sessionId]; exists { // <8>
		// セッションの有効期限をチェックする
		if time.Now().After(session.Expires) { // <9>
			// 有効期限が切れていたらセッション情報を削除してエラーを返す
			delete(m.sessions, sessionId)
			return nil, ErrSessionExpired
		}
		return session, nil
	} else {
		return nil, ErrSessionNotFound
	}
}


// セッションが開始されていることを保証する。
//
// セッションが存在しなければ、新しく発行する。
func ensureSession(w http.ResponseWriter, r *http.Request) (*HttpSession, error) {
	session, err := sessionManager.GetValidSession(r)
	if err == nil {
		return session, nil
	}

	// セッションが存在しないか不正な場合は新しく開始する
	log.Printf("session check failed : %s", err.Error())
	session, err = sessionManager.StartSession(w)
	if err != nil {
		writeInternalServerError(w, err)
		return nil, err
	}
	return session, err
}
