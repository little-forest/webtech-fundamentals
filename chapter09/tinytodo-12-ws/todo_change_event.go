package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)


type WebSocketEvent struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

func (e WebSocketEvent) Send(conn *websocket.Conn) error {
	data, _ := json.Marshal(e)
	err := conn.WriteMessage(websocket.TextMessage, data) // <1>
	if err != nil {
		log.Printf("Failed to send : %s", string(data))
		return err
	}
	return nil
}


func NewWebSocketEvent(event string, data any) *WebSocketEvent {
	return &WebSocketEvent{
		Event: event,
		Data:  data,
	}
}

// クライアントへSourceIdを通知するためのイベントデータ
type SourceIdNotification struct {
	Source string `json:"source"`
}

// クライアントへToDoの変化を通知するためのイベントデータ
type TodoChangeEvent struct {
	Source   string   `json:"source"`
	Event    string   `json:"-"`
	ToDoItem ToDoItem `json:"todoItem"`
}

// TodoChangeEvent を生成する関数
func NewTodoChangeEvent(source string, event string, todoItem ToDoItem) *TodoChangeEvent {
	return &TodoChangeEvent{
		Source:   source,
		Event:    event,
		ToDoItem: todoItem,
	}
}

// TodoChangeEvent から WebSocketEvent を生成する関数
func (e TodoChangeEvent) NewWebSocketEvent() *WebSocketEvent {
	return NewWebSocketEvent(e.Event, e)
}
