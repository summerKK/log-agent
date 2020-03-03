package taillog

import (
	"fmt"
	"log-agent/etcd"
)

var Manager *manager

type manager struct {
	logEntry []*etcd.LogEntry
	tailTask map[string]*TailTask
	confChan chan *etcd.LogEntry
}

func Init(logEntry []*etcd.LogEntry, newConfChan <-chan []*etcd.LogEntry) {

	Manager = &manager{
		logEntry: logEntry,
		tailTask: make(map[string]*TailTask, 16),
		confChan: make(chan *etcd.LogEntry),
	}

	// 监听到来的配置
	go Manager.watchConf()

	for _, conf := range logEntry {
		Manager.confChan <- conf
	}

	// 热更新配置
	go Manager.run(newConfChan)

}

func (m *manager) watchConf() {
	for v := range m.confChan {
		tailTask := NewTailTask(v.Path, v.Topic)
		mk := mk(v.Path, v.Topic)
		Manager.tailTask[mk] = tailTask
	}
}

func (m *manager) run(ch <-chan []*etcd.LogEntry) {
	for newConfs := range ch {
		newConfsMap := make(map[string]struct{})
		for _, newConf := range newConfs {
			mk := mk(newConf.Path, newConf.Topic)
			newConfsMap[mk] = struct{}{}
			// 配置已经存在
			if _, ok := m.tailTask[mk]; ok {
				continue
			} else {
				// 新配置
				m.confChan <- newConf
				fmt.Printf("got new config,%+v\n", newConf)
			}
		}

		for k, v := range m.tailTask {
			// 删除配置
			if _, ok := newConfsMap[k]; !ok {
				v.Close()
			}
		}
	}
}

func mk(path, topic string) string {
	return fmt.Sprintf("%s_%s", path, topic)
}
