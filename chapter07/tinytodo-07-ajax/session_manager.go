package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const SessionCleanupInterval = 5 * time.Minute
const SessionValidityTime = 30 * time.Minute

var (
	ErrSessionExpired    = errors.New("session expired")
	ErrSessionNotFound   = errors.New("session not found")
	ErrSessionIdNotGiven = errors.New("session id not given in cookie")
	ErrInvalidSessionId  = errors.New("invalid session id")
)

// セッションを管理する構造体
type HttpSessionManager struct {
	lock           sync.Mutex
	cleanerRunning bool
	// セッションIDをキーとしてセッション情報を保持するマップ
	sessions map[string]*HttpSession
	// セッションの有効時間
	validityTime time.Duration
	// セッションID検証用のシークレット
	sessionIdSecret uint64
	// Cookieにsecure属性を付与するかどうか
	useSecureCookie bool
}

func NewHttpSessionManager(sessionIdSecret uint64, useSecureCookie bool) *HttpSessionManager {
	mgr := &HttpSessionManager{
		sessions:        make(map[string]*HttpSession),
		validityTime:    SessionValidityTime,
		sessionIdSecret: sessionIdSecret,
		useSecureCookie: useSecureCookie,
	}
	return mgr
}

// セッションIDに紐付くセッション情報を返す。
func (m *HttpSessionManager) getSession(sessionId string) (*HttpSession, error) {
	if session, exists := m.sessions[sessionId]; exists {
		// セッションの有効期限をチェックする
		if time.Now().After(session.Expires) {
			// 有効期限が切れていたらセッション情報を削除してエラーを返す
			delete(m.sessions, sessionId)
			return nil, ErrSessionExpired
		}
		return session, nil
	} else {
		return nil, ErrSessionNotFound
	}
}

// 新しいセッション情報を生成して登録する
func (m *HttpSessionManager) newSession(sessionId string) *HttpSession {
	log.Printf("start session : %s", sessionId)
	session := NewHttpSession(sessionId, m.validityTime, m.useSecureCookie)
	m.sessions[sessionId] = session
	m.startSessionCleaner()
	return session
}

// Cookieから有効なセッションを取得する。
//
// CookieにセッションIDがなければ ErrSessionNotFound を返す。
// CookieにセッションIDが存在すれば、セッションIDに紐付く HttpSession を返す。
// セッションIDが不正な場合や、セッションの有効期限が切れている場合は、エラーを返す。
func (m *HttpSessionManager) GetValidSession(r *http.Request) (*HttpSession, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	c, err := r.Cookie(CookieNameSessionId)
	// CookieにセッションIDが存在しない場合
	if err == http.ErrNoCookie {
		return nil, ErrSessionNotFound
	}
	// CookieにセッションIDが存在する場合
	if err == nil {
		sessionId := c.Value
		ok, err := m.verifySessionId(sessionId)
		if err != nil || !ok {
			return nil, ErrInvalidSessionId
		}

		// セッションを取得して返す
		session, err := m.getSession(sessionId)
		return session, err
	}
	return nil, err
}

// セッションを開始してCokkieにセッションIDを書き込む。
func (m *HttpSessionManager) StartSession(w http.ResponseWriter) (*HttpSession, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	// 新しいセッションIDを生成する
	var sessionId string
	for {
		sid, err := m.makeSessionId()
		if err != nil {
			return nil, err
		}
		// 既存のセッションと重複していないことを確認する
		if _, err := m.getSession(sid); err == ErrSessionNotFound {
			sessionId = sid
			break
		}
	}

	// セッション情報を生成する
	session := m.newSession(sessionId)
	session.SetCookie(w)

	log.Printf("start session : %s", sessionId)

	return session, nil
}

// セッションIDを生成する。
func (m *HttpSessionManager) makeSessionId() (string, error) {
	keyBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(keyBytes, m.sessionIdSecret)

	randBytes := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, randBytes); err != nil {
		return "", err
	}

	hashBytes := md5.Sum(append(keyBytes, randBytes...))

	sessionId := append(randBytes, hashBytes[:]...)
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(sessionId), nil
}

// セッションIDが正当なものかどうか検証する。
func (m *HttpSessionManager) verifySessionId(sessionId string) (bool, error) {
	b, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(sessionId)
	if err != nil {
		return false, err
	}
	if len(b) < 17 {
		return false, nil
	}

	randBytes := b[:16]
	mac := b[16:]

	keyBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(keyBytes, m.sessionIdSecret)

	hashBytes := md5.Sum(append(keyBytes, randBytes...))

	return bytes.Equal(mac, hashBytes[:]), nil
}

// セッションを削除する
func (m *HttpSessionManager) RevokeSession(w http.ResponseWriter, sessionId string) {
	m.lock.Lock()
	defer m.lock.Unlock()

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
func checkSessionForPage(w http.ResponseWriter, r *http.Request) (*HttpSession, error) {
	session, err := sessionManager.GetValidSession(r)
	if err == nil {
		// セッションが存在すれば、延長する
		session.Extend()
		session.SetCookie(w)
		return session, nil
	}
	orgErr := err

	// セッションが有効期限切れまたは不正な場合、セッションを作り直す
	log.Printf("session check failed : %s", err.Error())
	session, err = sessionManager.StartSession(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	// Refererヘッダの有無で他の画面からの遷移かどうかを判定
	// アプリケーショントップのURLに直接アクセスした際は、セッションが存在しないのが
	// 正常であるため、エラーを表示しないための措置
	if r.Referer() != "" {
		page := LoginPageData{}
		page.ErrorMessage = fmt.Sprintf("セッションが不正です。(%s)", orgErr.Error())
		session.PageData = page
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil, orgErr
}

// セッションが存在するかチェックする。
//
// セッションが存在しなければ、エラーを返す。
func checkSession(w http.ResponseWriter, r *http.Request) (*HttpSession, error) {
	session, err := sessionManager.GetValidSession(r)
	if err == nil {
		// セッションが存在すれば、延長する
		session.Extend()
		session.SetCookie(w)
		return session, nil
	}

	// セッションが有効期限切れまたは不正な場合、エラーを返す
	log.Printf("session check failed : %s", err.Error())
	http.Error(w, err.Error(), http.StatusBadRequest)

	return nil, err
}

// セッションが開始されていることを保証する。
//
// セッションが存在しなければ、新しく発行する。
func ensureSession(w http.ResponseWriter, r *http.Request) (*HttpSession, error) {
	session, err := sessionManager.GetValidSession(r)
	if err == nil {
		// セッションが存在すれば、延長する
		session.Extend()
		session.SetCookie(w)
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

// セッション掃除のバックグラウンドプロセスを開始する。
//
// セッションが存在しないか、既に実行中の場合は場合は開始しない。
func (m *HttpSessionManager) startSessionCleaner() {
	if len(m.sessions) == 0 || m.cleanerRunning {
		return
	}

	m.cleanerRunning = true
	ticker := time.NewTicker(SessionCleanupInterval)
	go func() {
		log.Printf("session cleaner started")

	loop:
		for {
			select {
			case <-ticker.C:
				if sessionExists := m.CleanSessions(); !sessionExists {
					break loop
				}
			}
		}
		log.Printf("session cleaner stopping")
		m.lock.Lock()
		defer m.lock.Unlock()
		m.cleanerRunning = false
		log.Printf("session cleaner stopped")
	}()
}

// 期限切れのセッションを削除する。
//
// セッションが1つも無くなったら false を返す。
func (m *HttpSessionManager) CleanSessions() bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	log.Printf("cleaning session start. session count : %d", len(m.sessions))
	for id := range m.sessions {
		_, err := m.getSession(id)
		if err != nil {
			log.Printf("cleaning session : %s %v", id, err)
		}
	}
	sessionNum := len(m.sessions)
	log.Printf("cleaning session end. session count : %d", sessionNum)
	return (sessionNum > 0)
}
