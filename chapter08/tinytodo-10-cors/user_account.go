package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Account expire limit in minutes
const UserAccountLimitInMinute = 60

const PasswordLength = 10

const PasswordChars = "23456789abcdefghijkmnpqrstuvwxyz"

// ユーザアカウント情報を保持する構造体。
type UserAccount struct {
	// ユーザID
	Id string
	// ハッシュ化されたパスワード
	HashedPassword string
	// アカウントの有効期限
	Expires time.Time
	// ToDoリスト
	ToDoList *ToDoList
}

// ユーザアカウント情報を生成する。
func NewUserAccount(userId string, plainPassword string, expires time.Time) *UserAccount {
	// bcryptアルゴリズムでパスワードをハッシュ化する
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	account := &UserAccount{
		Id:             userId,
		HashedPassword: string(hashedPassword),
		Expires:        expires,
		ToDoList:       NewToDoList(),
	}
	return account
}

func (u UserAccount) ExpiresText() string {
	return u.Expires.Format("2006/01/02 15:04:05")
}

// 認証する。
func (u UserAccount) Authenticate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return (err == nil)
}

// 有効期限が切れているかどうかを確認する。
func (u UserAccount) IsExpired() bool {
	return time.Now().After(u.Expires)
}
