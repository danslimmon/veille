package main

// InputMode is the interface that input modes must implement.
type InputMode interface {
	Start() error
	LineChan() chan string
}
