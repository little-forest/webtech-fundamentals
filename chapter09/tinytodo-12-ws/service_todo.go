package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
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

	switch r.Method {
	case http.MethodGet:
		todoId := retrieveToDoId(r)
		if todoId == "" {
			// URLにIDが含まれていない場合はToDoリストを返す
			getToDoList(w, r, session)
		} else {
			// URLにIDが含まれている場合は単一のToDoを返す
			getToDo(w, r, session)
		}

	case http.MethodPost:
		// ToDo追加処理
		addToDo(w, r, session)

	case http.MethodPut:
		// ToDo更新処理
		editToDo(w, r, session)

	case http.MethodDelete:
		// ToDo削除処理
		deleteToDo(w, r, session)

	default:
		// 上記以外のHTTPメソッドはエラーを返す
		w.WriteHeader(http.StatusMethodNotAllowed)
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

	todoList := session.UserAccount.ToDoList

	// TodoListに新しいToDoを登録
	if item.Todo != "" {
		item = todoList.Append(item.Todo)
	}
	log.Printf("Todo item added. sessionId=%s itemId=%s todo=%s", session.SessionId, item.Id, item.Todo)

	// ToDo追加をチャネルを通じてobserveToDoへ通知する
	todoList.ChangeNotifier.Notify(NewTodoChangeEvent(getTinyTodoSourceId(r), "add", *item)) // <1>

	// 登録したToDoのURLをLocationヘッダで返す
	w.Header().Set("Location", fmt.Sprintf("/todos/%s", item.Id))
	w.WriteHeader(http.StatusCreated)
}

func getTinyTodoSourceId(r *http.Request) string { // <2>
	return r.Header.Get("X-Tinytodo-Sourceid")
}


// ToDoを更新する。
func editToDo(w http.ResponseWriter, r *http.Request, session *HttpSession) {
	// リクエストボディのJSONからToDo項目を復元する
	body, _ := io.ReadAll(r.Body)
	item := NewToDoItemFromJson(string(body))

	todoList := session.UserAccount.ToDoList

	// URLからToDo IDを取り出す
	todoId := retrieveToDoId(r)

	// URLのパスで指定されたToDoIdとリクエストボディのToDoIdの一致をチェック
	if todoId != item.Id {
		http.Error(w, "invalid todoId", http.StatusBadRequest)
		return
	}

	// ToDoを更新
	_, err := todoList.Update(item.Id, item.Todo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Todo item updated. sessionId=%s itemId=%s todo=%s", session.SessionId, item.Id, item.Todo)

	// ToDo追加をチャネルを通じてobserveToDoへ通知する
	todoList.ChangeNotifier.Notify(NewTodoChangeEvent(getTinyTodoSourceId(r), "update", *item))

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

var upgrader = websocket.Upgrader{} // <1>

func handleObserve(w http.ResponseWriter, r *http.Request) {
	// セッション情報を取得
	session, err := checkSession(w, r, false) // <3>
	logRequest(r, session)
	if err != nil {
		return
	}
	// 認証チェック
	if !isApiAuthenticated(w, r, session) { // <4>
		return
	}

	// WebSocketへのプロトコルアップグレード
	conn, err := upgrader.Upgrade(w, r, nil) // <2>
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	observeTodo(conn, session) // <5>
}


func observeTodo(conn *websocket.Conn, session *HttpSession) {
	todoList := session.UserAccount.ToDoList

	// チャネルを生成する
	eventReceiver := todoList.ChangeNotifier.CreateObserver()

	// SourceIdの生成
	sourceId := fmt.Sprintf("%d", time.Now().UnixMicro())
	log.Printf("WebSocket start (session:%s sourceId:%s)", session.SessionId, sourceId)

	// initialイベントを送信し、sourceIdを通知する
	ev := NewWebSocketEvent("initial", SourceIdNotification{Source: sourceId}) // <1>
	ev.Send(conn)                                                              // <3>

	// チャネルからの通知待ちループ
LOOP:
	for {
		select {
		// チャネルから通知を受けたら、イベントを通知する
		case ev := <-eventReceiver:
			log.Printf("Reveived Todo change. (session:%s sourceId:%s %v", session.SessionId, sourceId, ev)
			wse := ev.NewWebSocketEvent() // <2>
			err := wse.Send(conn)         // <4>
			if err != nil {
				break LOOP
			}
		}
	}

	// チャネルを削除する
	todoList.ChangeNotifier.RemoveObserver(eventReceiver)
	log.Printf("Connection closed. (session:%s sourceId:%s)", session.SessionId, sourceId)
}

