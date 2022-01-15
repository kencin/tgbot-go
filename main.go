package main

import (
	"tgbot/bot"
	_ "tgbot/lib"
	"tgbot/lib/context"
)

func main() {
	context.GlobalContext.Logger.Info("Starting bot")

	bot.NewBot()

	select {}
}
