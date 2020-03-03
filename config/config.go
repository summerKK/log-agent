package config

import "gopkg.in/ini.v1"

type AppConfig struct {
	Kafka   `ini:"kafka"`
	Taillog `ini:"taillog"`
	Etcd    `ini:"etcd"`
}

type Kafka struct {
	Address []string `ini:"address" delim:","`
	Topic   string   `ini:"topic"`
	MaxSize int      `ini:"max_size"`
}

type Taillog struct {
	Filename string `ini:"filename"`
}

type Etcd struct {
	Endpoints []string `ini:"endpoints" delim:","`
	Timeout   int      `ini:"timeout"`
	Key       string   `ini:"collection_log_key"`
}

var Config *AppConfig

func init() {
	Config = new(AppConfig)
	ini.MapTo(Config, "./config/config.ini")
}
