package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"time"
)

const cookieSessionId = "sessionId"

const sessionValidityTime = 600 * time.Second

const sessionIdSecret = 123456789

func ensureSession(w http.ResponseWriter, r *http.Request) (sessionId string, err error) { // <1>
	c, err := r.Cookie(cookieSessionId)
	if err == http.ErrNoCookie { // <2>
		sessionId, err = startSession(w)
		return
	}
	if err == nil { // <3>
		sessionId = c.Value
		if ok, _ := verifySessionId(sessionIdSecret, sessionId); !ok {
			return "", fmt.Errorf("invalid session id")
		}
		return
	}
	return
}

func startSession(w http.ResponseWriter) (string, error) { // <4>
	sessionId, err := makeSessionId(sessionIdSecret)
	if err != nil {
		return "", err
	}

	cookie := &http.Cookie{ // <5>
		Name:     cookieSessionId,
		Value:    sessionId,
		Expires:  time.Now().Add(sessionValidityTime),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	return sessionId, nil
}

func makeSessionId(secret uint64) (string, error) {
	keyBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(keyBytes, secret)

	randBytes := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, randBytes); err != nil {
		return "", err
	}

	hashBytes := md5.Sum(append(keyBytes, randBytes...))

	sessionId := append(randBytes, hashBytes[:]...)
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(sessionId), nil
}

func verifySessionId(secret uint64, sessionId string) (bool, error) {
	b, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(sessionId)
	if err != nil {
		return false, err
	}

	randBytes := b[:16]
	mac := b[16:]

	keyBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(keyBytes, secret)

	hashBytes := md5.Sum(append(keyBytes, randBytes...))

	return bytes.Equal(mac, hashBytes[:]), nil
}
