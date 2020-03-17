package main

import (
	"fmt"
	"github.com/summerKK/go-code-snippet-library/log-agent/config"
	"github.com/summerKK/go-code-snippet-library/log-agent/etcd"
	"github.com/summerKK/go-code-snippet-library/log-agent/kafka"
	"github.com/summerKK/go-code-snippet-library/log-agent/taillog"
	"github.com/summerKK/go-code-snippet-library/log-agent/utils"
	"sync"
	"time"
)

func main() {
	// 初始化etcd
	etcd.Init(config.Config.Endpoints, time.Duration(config.Config.Timeout)*time.Second)
	logConfKey := config.Config.Key
	// 获取每个服务器对应的配置文件
	logConfKey = fmt.Sprintf(logConfKey, utils.GetOutboundIP())
	logConf := etcd.GetLogConf(logConfKey)
	// 初始化kafka
	kafka.Init(config.Config.Address, config.Config.MaxSize)
	watchLogConfChan := etcd.WatchLogConf(logConfKey)
	var wg sync.WaitGroup
	wg.Add(1)
	taillog.Init(logConf, watchLogConfChan)

	// kafka消费者
	for _, entry := range logConf {
		kafka.InitConsumer(config.Config.Address, entry.Topic)
	}

	wg.Wait()
}
