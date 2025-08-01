package model

type Adapter interface {
	Bind(events []Event)
	Start()
	Info
}

type Info interface {
	Platform() string
}
