package main

import (
	"log-agent/config"
	"log-agent/etcd"
	"log-agent/kafka"
	"log-agent/taillog"
	"sync"
	"time"
)

func main() {
	// 初始化etcd
	etcd.Init(config.Config.Endpoints, time.Duration(config.Config.Timeout)*time.Second)
	logConf := etcd.GetLogConf(config.Config.Key)
	// 初始化kafka
	kafka.Init(config.Config.Address, config.Config.MaxSize)
	watchLogConfChan := etcd.WatchLogConf(config.Config.Key)
	var wg sync.WaitGroup
	wg.Add(1)
	taillog.Init(logConf, watchLogConfChan)
	wg.Wait()
}
