package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Server-sent Events を表す構造体
type ServerSentEvent struct {
	Id    string `json:"id"`
	Event string `json:"event"`
	Data  any    `json:"data"`
}

// ServerSentEventを生成する関数
func NewServerSentEvent(event string, data any) *ServerSentEvent {
	return &ServerSentEvent{
		Id:    fmt.Sprintf("ttd-%d", time.Now().UnixMicro()),
		Event: event,
		Data:  data,
	}
}

// Server-sent Eventsを送信する関数
func (e ServerSentEvent) Send(w http.ResponseWriter) { // <1>
	data, _ := json.Marshal(e.Data)
	fmt.Fprintf(w, "id: %s\n", e.Id)
	fmt.Fprintf(w, "event: %s\n", e.Event)
	fmt.Fprintf(w, "data: %s\n\n", string(data))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
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
	TodoItem ToDoItem `json:"todoItem"`
}

// TodoChangeEvent を生成する関数
func NewTodoChangeEvent(source string, event string, todoItem ToDoItem) *TodoChangeEvent {
	return &TodoChangeEvent{
		Source:   source,
		Event:    event,
		TodoItem: todoItem,
	}
}

// TodoChangeEvent から ServerSentEvent を生成する関数
func (e TodoChangeEvent) NewServerSentEvent() *ServerSentEvent {
	return NewServerSentEvent(e.Event, e)
}

