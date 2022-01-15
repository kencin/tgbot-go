package filebot

import (
	"fmt"
	"tgbot/bot/common"
	"tgbot/lib/config"
	"tgbot/lib/context"
	"tgbot/lib/request"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var globalLogger = context.GlobalContext.Logger

func InitFileBot() {
	globalLogger.Debug("Initializing filebot")

	var err error
	handler.bot, err = tgbotapi.NewBotAPIWithAPIEndpoint(config.Cfg.BotToken.FileBotToken, common.GetSelfAPIEndpoint())
	if err != nil {
		globalLogger.Panic(fmt.Sprintf("Initializing filebot err: %v", err))
	}

	globalLogger.Info(fmt.Sprintf("Authorized on account %s", handler.bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := handler.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			ctx := context.NewContext()
			go handler.videoDownloadHandler(ctx, update)
		}
	}
}

var handler = fileBotHandler{}

type fileBotHandler struct {
	bot *tgbotapi.BotAPI
}

func (f *fileBotHandler) videoDownloadHandler(ctx *context.Context, update tgbotapi.Update) {
	ctx.Logger.Info(fmt.Sprintf("[%s] %d", update.Message.From.UserName, update.Message.MessageID))

	if update.Message.Video != nil {
		ctx.Logger.Info(fmt.Sprintf("收到视频下载请求, id: %s", update.Message.Video.FileID))

		// 添加下载任务并通知开始下载
		request.AddDownloadTgFileLocalTask(request.DownloadTgFileLocal{
			FileId:    update.Message.Video.FileID,
			FileUId:   update.Message.Video.FileUniqueID,
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.MessageID,
			Bot:       f.bot,
			Ctx:       ctx,
			Retry:     true,
		})

		common.Reply(ctx, &common.ReplyStruct{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.MessageID,
			Msg:       fmt.Sprintf("添加下载任务成功, 当前剩余下载任务数量: %d", request.GetDownloadTgFileLocalTaskNums()),
			Bot:       f.bot,
		})
	}
}
