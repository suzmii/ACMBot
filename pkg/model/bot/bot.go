package bot

import (
	"context"

	"github.com/suzmii/ACMBot/internal/bot/message"
)

type Handler func(context.Context) error

type ApiCaller interface {
	Send(message message.Message)
	GetCallerInfo() CallerInfo
}

type CallerInfo struct {
	NickName string
	ID       int64
	Group    GroupInfo
}

type GroupInfo struct {
	ID          int64
	Name        string
	MemberCount int64
}

type Platform int

const (
	PlatformQQ Platform = iota
)
