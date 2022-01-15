package config

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type cfg struct {
	AppName         string          `yaml:"name"`
	BotToken        botToken        `yaml:"bottoken"`
	LogConfig       logConfig       `yaml:"log"`
	Store           store           `yaml:"store"`
	ApiServerConfig apiServerConfig `yaml:"selfapiserver"`
}

type botToken struct {
	FileBotToken string `yaml:"fileBotToken"`
}

type logConfig struct {
	FileLogPath string `yaml:"fileLogPath"`
	Level       string `yaml:"level"`
	Out         string `yaml:"out"`
	ShowCaller  bool   `yaml:"showCaller"`
}

type store struct {
	FilePath string `yaml:"filePath"`
}

type apiServerConfig struct {
	ApiEndpoint string `yaml:"apiEndpoint"`
}

var Cfg = &cfg{}

func Init() {
	config.WithOptions(config.ParseEnv)

	config.WithOptions(func(opt *config.Options) {
		opt.DecoderConfig.TagName = "yaml"
	})

	// add driver for support yaml content
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("config/config.yml")
	if err != nil {
		panic(err)
	}

	err = config.BindStruct("", &Cfg)
	if err != nil {
		panic(err)
	}
}
