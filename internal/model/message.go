package model

type Message interface {
	isMessage()
}

type Messager interface {
	Send(Message)
}

type TextMessage struct {
	Text string
}

func (TextMessage) isMessage() {}

type ImageMessage struct {
	Image []byte
}

func (ImageMessage) isMessage() {}
