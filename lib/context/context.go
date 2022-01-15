package context

import (
	"context"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"tgbot/lib/logger"
)

var GlobalContext = NewContext("global")

type Context struct {
	context.Context
	Logger *log.Entry
}

func NewContext(traceId ...string) *Context {
	tid := uuid.NewString()
	if len(traceId) > 0 {
		tid = traceId[0]
	}

	return &Context{
		Context: context.Background(),
		Logger:  logger.Logger(tid),
	}
}
