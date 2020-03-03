package taillog

import (
	"fmt"
	"log-agent/etcd"
)

var Manager *manager

type manager struct {
	logEntry    []*etcd.LogEntry
	tailTask    map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

func Init(logEntry []*etcd.LogEntry, ch <-chan []*etcd.LogEntry) {

	Manager = &manager{
		logEntry:    logEntry,
		tailTask:    make(map[string]*TailTask, 16),
		newConfChan: make(chan []*etcd.LogEntry),
	}

	for _, conf := range logEntry {
		tailTask := NewTailTask(conf.Path, conf.Topic)
		Manager.tailTask[conf.Topic] = tailTask
	}

	go Manager.run(ch)

}

func (m *manager) run(ch <-chan []*etcd.LogEntry) {
	for v := range ch {
		fmt.Printf("got new config,%+v\n", v)
		m.newConfChan <- v
	}
}
