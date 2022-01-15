package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
	"tgbot/lib/config"

	nested "github.com/antonfisher/nested-logrus-formatter"
)

func Init() {
	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"traceId"},
	})

	// set log output type
	switch config.Cfg.LogConfig.Out {
	case "file":
		f, err := os.OpenFile(config.Cfg.LogConfig.FileLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		log.SetOutput(f)
	case "stdout":
		log.SetOutput(os.Stdout)
	default:
		panic("Unknown log output")
	}

	// set log level
	level, err := log.ParseLevel(config.Cfg.LogConfig.Level)
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)

	log.SetReportCaller(config.Cfg.LogConfig.ShowCaller)
}

func Logger(traceId string) *log.Entry {
	return log.WithFields(log.Fields{
		//"app":     config.Cfg.AppName,
		"traceId": traceId,
	})
}
