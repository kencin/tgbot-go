package lib

import (
	"tgbot/lib/config"
	"tgbot/lib/logger"
	"tgbot/lib/request"
)

func init() {
	config.Init()
	logger.Init()
	startDownloadRoutine()
}

func startDownloadRoutine() {
	go request.DownloadRoutine()
}
