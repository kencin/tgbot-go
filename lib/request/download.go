package request

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/panjf2000/ants/v2"
	"os"
	"os/exec"
	"path"
	"sync"
	"sync/atomic"
	"tgbot/bot/common"
	"tgbot/lib/config"
	"tgbot/lib/context"
)

type DownloadTask struct {
	Url      string
	FilePath string
}

type DownloadTgFileLocal struct {
	FileId    string // 文件id，这个每次都不一样
	FileUId   string // 文件唯一id，通过这个避免重复下载
	ChatID    int64  // 当前聊天id
	MessageID int    // 当前消息id
	Bot       *tgbotapi.BotAPI
	Ctx       *context.Context
	Retry     bool // 是否已经重试过了
}

var (
	DownloadTgFileLocals []DownloadTgFileLocal
	lock                 = sync.RWMutex{}
	logger               = context.GlobalContext.Logger
	notifyDownload       = make(chan int)
	runningWorkers       int32
)

func addRunningWorkers() {
	atomic.AddInt32(&runningWorkers, 1)
}

func decRunningWorkers() {
	atomic.AddInt32(&runningWorkers, -1)
}

func AddDownloadTgFileLocalTask(task DownloadTgFileLocal) {
	lock.Lock()
	defer lock.Unlock()
	logger.Info(fmt.Sprintf("添加本地下载任务: %s", task.FileId))

	DownloadTgFileLocals = append(DownloadTgFileLocals, task)

	addRunningWorkers()
	NotifyDownloadRoutine()
}

func GetDownloadTgFileLocalTaskNums() int {
	return int(runningWorkers)
}

// DownloadRoutine 下载线程
func DownloadRoutine() {
	p, _ := ants.NewPool(2)

	for {
		select {
		case <-notifyDownload:
			oneTask := func() *DownloadTgFileLocal {
				lock.Lock()
				defer lock.Unlock()
				if len(DownloadTgFileLocals) == 0 {
					return nil
				}
				task := DownloadTgFileLocals[0]
				DownloadTgFileLocals = DownloadTgFileLocals[1:]
				return &task
			}()

			logger = oneTask.Ctx.Logger

			if oneTask == nil {
				logger.Warn("没有下载任务")
				continue
			}

			reply := func(msg string) string {
				common.Reply(oneTask.Ctx, &common.ReplyStruct{
					ChatID:    oneTask.ChatID,
					MessageID: oneTask.MessageID,
					Msg:       msg,
					Bot:       oneTask.Bot,
				})
				return msg
			}

			err := p.Submit(func() {
				dstPath := path.Join(config.Cfg.Store.FilePath, oneTask.FileUId+".mp4")
				// 先判断文件是否存在
				if _, err := os.Stat(dstPath); err == nil {
					logger.Info(reply(fmt.Sprintf("视频 %s 已存在", dstPath)))
					decRunningWorkers()
					return
				}

				logger.Info(fmt.Sprintf("开始下载, 任务: %s", oneTask.FileId))
				downloadPath, err := getTgFilePath(oneTask.Bot, oneTask.FileId)
				if err != nil {
					logger.Error(reply(fmt.Sprintf("下载视频到本地失败, 任务: %s, 错误: %s", oneTask.FileId, err.Error())))
					decRunningWorkers()
					return
				}
				logger.Info(reply("下载视频到本地成功"))

				go func() {
					cmd := exec.Command("mv", downloadPath, dstPath)
					_, err = cmd.Output()
					if err != nil {
						logger.Error(reply(fmt.Sprintf("上传视频到nas失败, 任务: %s, 错误: %s", oneTask.FileId, err.Error())))
						return
					}
					logger.Info(reply("上传视频到nas成功"))
					decRunningWorkers()
				}()
			})
			if err != nil {
				msg := fmt.Sprintf("提交下载任务失败: %s", err.Error())
				// 将任务重新加入队列
				if oneTask.Retry {
					lock.Lock()
					oneTask.Retry = false
					DownloadTgFileLocals = append(DownloadTgFileLocals, *oneTask)
					lock.Unlock()
					msg += " 任务已尝试重新提交"
				}
				logger.Error(reply(msg))
				continue
			}
		}
	}
}

func getTgFilePath(bot *tgbotapi.BotAPI, fileID string) (string, error) {
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})

	if err != nil {
		return "", err
	}

	return file.FilePath, nil
}

func NotifyDownloadRoutine() {
	notifyDownload <- 1
}
