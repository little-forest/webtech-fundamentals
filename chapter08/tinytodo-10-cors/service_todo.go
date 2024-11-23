package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// ToDoリストを返すエンドポイント。
func handleTodos(w http.ResponseWriter, r *http.Request) {
	// セッション情報を取得
	session, err := checkSession(w, r, false)
	logRequest(r, session)
	if err != nil {
		return
	}
	// 認証チェック
	if !isApiAuthenticated(w, r, session) {
		return
	}

	switch r.Method { // <1>
	case http.MethodGet:
		todoId := retrieveToDoId(r) // <2>
		if todoId == "" {           // <3>
			// URLにIDが含まれていない場合はToDoリストを返す
			getToDoList(w, r, session)
		} else { // <4>
			// URLにIDが含まれている場合は単一のToDoを返す
			getToDo(w, r, session)
		}

	case http.MethodPost: // <5>
		// ToDo追加処理
		addToDo(w, r, session)

	case http.MethodPut: // <6>
		// ToDo更新処理
		editToDo(w, r, session)

	case http.MethodDelete: // <7>
		// ToDo削除処理
		deleteToDo(w, r, session)

	default:
		// 上記以外のHTTPメソッドはエラーを返す
		w.WriteHeader(http.StatusMethodNotAllowed) // <8>
		return
	}
}

// リクエストパスからToDo IDを抜き出す。
func retrieveToDoId(r *http.Request) string {
	if !strings.HasPrefix(r.URL.Path, "/todos/") {
		return ""
	}
	return r.URL.Path[len("/todos/"):]
}


// ToDoリストを返す。
func getToDoList(w http.ResponseWriter, r *http.Request, session *HttpSession) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session.UserAccount.ToDoList)
}

// ToDoを返す。
func getToDo(w http.ResponseWriter, r *http.Request, session *HttpSession) {
	todoId := retrieveToDoId(r)
	item, err := session.UserAccount.ToDoList.Get(todoId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// ToDoを追加する。
func addToDo(w http.ResponseWriter, r *http.Request, session *HttpSession) {
	// リクエストボディのJSONからToDo項目を復元する
	body, _ := io.ReadAll(r.Body)
	item := NewToDoItemFromJson(string(body))

	// TodoListに新しいToDoを登録
	if item.Todo != "" {
		item = session.UserAccount.ToDoList.Append(item.Todo)
	}
	log.Printf("Todo item added. sessionId=%s itemId=%s todo=%s", session.SessionId, item.Id, item.Todo)

	// 登録したToDoのURLをLocationヘッダで返す
	w.Header().Set("Location", fmt.Sprintf("/todos/%s", item.Id))
	w.WriteHeader(http.StatusCreated)
}

// ToDoを更新する。
func editToDo(w http.ResponseWriter, r *http.Request, session *HttpSession) {
	// リクエストボディのJSONからToDo項目を復元する
	body, _ := io.ReadAll(r.Body)
	item := NewToDoItemFromJson(string(body))

	// URLからToDo IDを取り出す
	todoId := retrieveToDoId(r)

	// URLのパスで指定されたToDoIdとリクエストボディのToDoIdの一致をチェック
	if todoId != item.Id {
		http.Error(w, "invalid todoId", http.StatusBadRequest)
		return
	}

	// ToDoを更新
	_, err := session.UserAccount.ToDoList.Update(item.Id, item.Todo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Todo item updated. sessionId=%s itemId=%s todo=%s", session.SessionId, item.Id, item.Todo)

	w.WriteHeader(http.StatusOK)
}

// ToDoを更新する。
func deleteToDo(w http.ResponseWriter, r *http.Request, session *HttpSession) {
	// URLからToDo IDを取り出す
	todoId := retrieveToDoId(r)

	// ToDoを削除
	err := session.UserAccount.ToDoList.Delete(todoId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.WriteHeader(http.StatusNoContent)
}
