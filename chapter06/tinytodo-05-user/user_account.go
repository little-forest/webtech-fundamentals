package main

import (
	"time"
	_ "time/tzdata"

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
	ToDoList []string
}

// ユーザアカウント情報を生成する。
func NewUserAccount(userId string, plainPassword string, expires time.Time) *UserAccount {
	// bcryptアルゴリズムでパスワードをハッシュ化する
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	account := &UserAccount{
		Id:             userId,
		HashedPassword: string(hashedPassword),
		Expires:        expires,
		ToDoList:       make([]string, 0, 10),
	}
	return account
}


func (u UserAccount) ExpiresText() string {
	return u.Expires.Format("2006/01/02 15:04:05")
}
