package main

import (
	"fmt"
	"log-agent/config"
	"log-agent/etcd"
	"log-agent/kafka"
	"log-agent/taillog"
	"log-agent/utils"
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
	wg.Wait()
}
