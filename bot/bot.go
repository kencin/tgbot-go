package bot

import (
	"tgbot/bot/filebot"
	"tgbot/lib/context"
)

var globalLogger = context.GlobalContext.Logger

type botFunc func()

var botFuncs = []botFunc{
	filebot.InitFileBot,
}

func NewBot() {
	globalLogger.Debug("Start init bot")

	for _, f := range botFuncs {
		go f()
	}

	globalLogger.Debug("Finish init bot")
}
