package main

import (
	"log"
	"sync"
)

type ChangeNotifier struct {
	mux       sync.Mutex
	Observers []chan TodoChangeEvent
}

func NewChangeNotifier() *ChangeNotifier {
	return &ChangeNotifier{
		mux:       sync.Mutex{},
		Observers: []chan TodoChangeEvent{},
	}
}

func (c *ChangeNotifier) CreateObserver() chan TodoChangeEvent {
	ch := make(chan TodoChangeEvent)
	c.mux.Lock()
	defer c.mux.Unlock()
	c.Observers = append(c.Observers, ch)

	log.Printf("observer created: %v", ch)
	return ch
}

func (c *ChangeNotifier) RemoveObserver(ch chan TodoChangeEvent) {
	c.mux.Lock()
	defer c.mux.Unlock()
	target := -1
	for i, e := range c.Observers {
		if e == ch {
			target = i
			break
		}
	}
	if target == -1 {
		log.Printf("Failed to remove chanel : %v", ch)
		return
	}
	c.Observers = append(c.Observers[:target], c.Observers[target+1:]...)
	log.Printf("observer removed: %v left(%d)", ch, len(c.Observers))
}

func (c *ChangeNotifier) Notify(e *TodoChangeEvent) {
	for _, ch := range c.Observers {
		ch <- *e
	}
}
