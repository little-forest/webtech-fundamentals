package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

// ToDo項目を表す構造体。
type ToDoItem struct {
	Id   string `json:"id"` // <1>
	Todo string `json:"todo"`
}

// 新しいToDoItemを生成する。
func NewToDoItem(todo string) *ToDoItem {
	id := MakeToDoId(todo)
	return &ToDoItem{
		Id:   id,
		Todo: todo,
	}
}

// ToDo項目のIDを生成する。
func MakeToDoId(todo string) string { // <3>
	timeBytes := []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
	hasher := md5.New()
	hasher.Write(timeBytes)
	hasher.Write([]byte(todo))
	return hex.EncodeToString(hasher.Sum(nil))
}
