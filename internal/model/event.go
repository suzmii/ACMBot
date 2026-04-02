package model

type Event struct {
	Commands []string
	Handler  Handler
}

type Handler func(*Context) error
