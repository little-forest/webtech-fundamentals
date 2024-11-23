package main

import (
	"fmt"
	"sync"
)

// ToDo Listを保持する構造体。
type ToDoList struct {
	lock  sync.Mutex  `json:"-"`
	Items []*ToDoItem `json:"items"`
}

// 新しいToDoListを生成する。
func NewToDoList() *ToDoList {
	list := &ToDoList{
		Items: make([]*ToDoItem, 0, 10),
	}
	return list
}

// ToDoを追加する。
func (t *ToDoList) Append(todo string) *ToDoItem {
	t.lock.Lock()
	defer t.lock.Unlock()

	todoItem := NewToDoItem(todo)
	t.Items = append(t.Items, todoItem)
	return todoItem
}

// ToDo項目を取得する。
func (t *ToDoList) Get(id string) (*ToDoItem, error) {
	for _, todo := range t.Items {
		if todo.Id == id {
			return todo, nil
		}
	}
	return nil, fmt.Errorf("todo not found. itemId=%s", id)
}

// ToDo項目を更新する。
func (t *ToDoList) Update(id string, newTodo string) (*ToDoItem, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	todoItem, err := t.Get(id)
	if err != nil {
		return nil, err
	}

	todoItem.Todo = newTodo
	return todoItem, nil
}

// ToDo項目を削除する。
func (t *ToDoList) Delete(id string) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	deleteIndex := -1
	for i, todo := range t.Items {
		if todo.Id == id {
			deleteIndex = i
			break
		}
	}

	if deleteIndex == -1 {
		return fmt.Errorf("todo not found. itemId=%s", id)
	}

	t.Items = append(t.Items[:deleteIndex], t.Items[deleteIndex+1:]...)
	return nil
}
