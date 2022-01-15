package common

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbot/lib/config"
	"tgbot/lib/context"
)

type ReplyStruct struct {
	ChatID    int64
	MessageID int
	Msg       string
	Bot       *tgbotapi.BotAPI
}

func Reply(ctx *context.Context, rs *ReplyStruct) {
	rp := tgbotapi.NewMessage(rs.ChatID, rs.Msg)
	rp.ReplyToMessageID = rs.MessageID

	_, err := rs.Bot.Send(rp)
	if err != nil {
		ctx.Logger.Error(fmt.Sprintf("Error sending message: %v", err))
	}
}

func GetSelfAPIEndpoint() string {
	return config.Cfg.ApiServerConfig.ApiEndpoint
}
