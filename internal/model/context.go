package model

import (
	"context"
)

type Context struct {
	Ctx      context.Context
	Adapter  Info
	Messager Messager
	Params   []string
	UserID   int64
	Username string
}
