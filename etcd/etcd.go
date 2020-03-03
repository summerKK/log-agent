package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var client *clientv3.Client

type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func Init(endpoints []string, timeout time.Duration) {
	var err error
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("etcd connected")
}

func GetLogConf(key string) []*LogEntry {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	response, err := client.Get(ctx, key)
	cancelFunc()
	if err != nil {
		panic(err)
	}
	logEntry := make([]*LogEntry, 0)
	for _, kv := range response.Kvs {
		err := json.Unmarshal(kv.Value, &logEntry)
		if err != nil {
			panic(err)
		}
	}
	return logEntry
}

func WatchLogConf(key string) <-chan []*LogEntry {
	out := make(chan []*LogEntry)
	go func() {
		watch := client.Watch(context.Background(), key)
		for v := range watch {
			for _, event := range v.Events {
				logEntry := make([]*LogEntry, 0)
				if event.Type != clientv3.EventTypeDelete {
					err := json.Unmarshal(event.Kv.Value, &logEntry)
					if err != nil {
						fmt.Printf("json.Unmarsha1 error:%v\n", err)
						continue
					}
					out <- logEntry
				}
			}
		}
		close(out)
	}()

	return out
}

func Put(key string, value string) (*clientv3.PutResponse, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*1)
	response, err := client.Put(ctx, key, value)
	cancelFunc()
	return response, err
}
