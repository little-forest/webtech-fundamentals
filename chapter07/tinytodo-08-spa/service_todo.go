package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func handleTodo(w http.ResponseWriter, r *http.Request) {
	// セッション情報を取得
	session, err := checkSession(w, r)
	logRequest(r, session)
	if err != nil {
		return
	}
	// 認証チェック
	if !isAuthenticated(w, r, session) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session.UserAccount.ToDoList)
}


func handleAdd(w http.ResponseWriter, r *http.Request) {
	// セッション情報を取得
	session, err := checkSession(w, r)
	logRequest(r, session)
	if err != nil {
		return
	}
	// 認証チェック
	if !isAuthenticated(w, r, session) {
		return
	}

	// リクエスト内のJSONを解析
	body, _ := io.ReadAll(r.Body)
	item := NewToDoItemFromJson(string(body))

	// TodoListに新しいToDoを登録
	if item.Todo != "" {
		item = session.UserAccount.ToDoList.Append(item.Todo)
	}
	log.Printf("Todo item added. sessionId=%s itemId=%s todo=%s", session.SessionId, item.Id, item.Todo)

	// 登録結果をJSONとしてレスポンスで返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}


func handleEdit(w http.ResponseWriter, r *http.Request) {
	// セッション情報を取得
	session, err := checkSession(w, r)
	logRequest(r, session)
	if err != nil {
		return
	}
	// 認証チェック
	if !isAuthenticated(w, r, session) {
		return
	}

	body, _ := io.ReadAll(r.Body)
	item := NewToDoItemFromJson(string(body))

	_, err = session.UserAccount.ToDoList.Update(item.Id, item.Todo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Todo item updated. sessionId=%s itemId=%s todo=%s", session.SessionId, item.Id, item.Todo)

	w.WriteHeader(http.StatusOK) // <4>
}
