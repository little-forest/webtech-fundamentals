package main

import "sync"

// ToDo Listを保持する構造体。
type ToDoList struct {
	lock  sync.Mutex
	Items []string
}

// 新しいToDoListを生成する。
func NewToDoList() *ToDoList {
	list := &ToDoList{
		Items: make([]string, 0, 10),
	}
	return list
}

// ToDoを追加する。
func (t *ToDoList) Append(todo string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.Items = append(t.Items, todo)
}
