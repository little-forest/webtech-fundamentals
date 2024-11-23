package main

import (
	"errors"
	"log"
	"math/rand"
	"regexp"
	"sync"
	"time"
	_ "time/tzdata"
)

const AccountCleanupInterval = 5 * time.Minute

var (
	ErrUserAlreadyExists   = errors.New("user account already exists")
	ErrInvalidUserIdFormat = errors.New("invalid user id format")
	ErrLoginFailed         = errors.New("login failed")
	ErrAccountExpired      = errors.New("user account expired")
	RegexAccountId         = regexp.MustCompile(`^[A-Za-z0-9_.+@-]{1,32}$`)
)

// ユーザアカウントを管理する構造体。
type UserAccountManager struct {
	userAccounts   map[string]*UserAccount
	location       *time.Location
	lock           sync.Mutex
	cleanerRunning bool
}

// UserAccountManagerを生成する。
func NewUserAccountManager() *UserAccountManager {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}

	m := &UserAccountManager{
		userAccounts: make(map[string]*UserAccount),
		location:     jst,
	}
	return m
}

// ユーザIDの形式を検証する。
func (m *UserAccountManager) ValidateUserId(userId string) bool {
	return RegexAccountId.MatchString(userId)
}

// 新しいユーザアカウントを作成する。
func (m *UserAccountManager) NewUserAccount(userId string, password string) (*UserAccount, error) {
	if !m.ValidateUserId(userId) {
		return nil, ErrInvalidUserIdFormat
	}
	_, exists := m.GetUserAccount(userId)
	if exists {
		return nil, ErrUserAlreadyExists
	}

	expires := time.Now().In(m.location).Add(time.Minute * UserAccountLimitInMinute)
	account := NewUserAccount(userId, password, expires)

	m.userAccounts[userId] = account
	log.Printf("user account created : %s\n", account.Id)
	m.startUserAccountCleaner()
	return account, nil
}

// ユーザアカウントを取得する。
func (m *UserAccountManager) GetUserAccount(userId string) (*UserAccount, bool) {
	a, exists := m.userAccounts[userId]
	return a, exists
}

// ユーザアカウントを削除する。
func (m *UserAccountManager) RemoveUserAccount(userId string) {
	if _, exists := m.userAccounts[userId]; !exists {
		return
	}

	delete(m.userAccounts, userId)
}

// ユーザアカウントを認証する。
func (m *UserAccountManager) Authenticate(userId string, password string) (*UserAccount, error) {
	account, exists := m.GetUserAccount(userId)
	if !exists {
		return nil, ErrLoginFailed
	}
	if !account.Authenticate(password) {
		return nil, ErrLoginFailed
	}
	if account.IsExpired() {
		m.lock.Lock()
		defer m.lock.Unlock()
		m.RemoveUserAccount(userId)
		return nil, ErrAccountExpired
	}
	return account, nil
}

// ランダムなパスワードを生成する。
func MakePassword() string {
	password := make([]byte, PasswordLength)
	for i := 0; i < PasswordLength; i++ {
		password[i] = PasswordChars[rand.Intn(len(PasswordChars))]
	}
	return string(password)
}

// アカウント掃除のバックグラウンドプロセスを開始する。
//
// アカウントが存在しないか、既に実行中の場合は場合は開始しない。
func (m *UserAccountManager) startUserAccountCleaner() {
	if len(m.userAccounts) == 0 || m.cleanerRunning {
		return
	}

	m.cleanerRunning = true
	ticker := time.NewTicker(AccountCleanupInterval)
	go func() {
		log.Printf("user account cleaner started")

	loop:
		for {
			select {
			case <-ticker.C:
				if accountExists := m.CleanExpiredAccounts(); !accountExists {
					break loop
				}
			}
		}
		log.Printf("user account cleaner stopping")
		m.lock.Lock()
		defer m.lock.Unlock()
		m.cleanerRunning = false
		log.Printf("user account cleaner stopped")
	}()
}

// 期限切れのアカウントを削除する。
//
// アカウントが1つも無くなったら false を返す。
func (m *UserAccountManager) CleanExpiredAccounts() bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	log.Printf("cleaning account start. account count : %d", len(m.userAccounts))
	for userId, userAccount := range m.userAccounts {
		if userAccount.IsExpired() {
			log.Printf("remove expired account : %s", userId)
			m.RemoveUserAccount(userId)
		}
	}
	accountNum := len(m.userAccounts)
	log.Printf("cleaning account end. num of counts : %d", accountNum)
	return (accountNum > 0)
}
