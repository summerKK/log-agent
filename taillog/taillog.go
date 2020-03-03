package taillog

import (
	"fmt"
	"github.com/hpcloud/tail"
	"log-agent/kafka"
)

type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
	close    chan struct{}
}

func NewTailTask(path, topic string) *TailTask {
	task := &TailTask{
		path:  path,
		topic: topic,
		close: make(chan struct{}),
	}
	task.init()
	fmt.Printf("tail topic:%v,path:%v start\n", task.topic, task.path)
	return task
}

func (t *TailTask) init() error {
	tailFile, err := tail.TailFile(t.path, tail.Config{
		Follow: true,
	})
	if err != nil {
		return err
	}
	t.instance = tailFile
	go t.run()

	return nil
}

func (t *TailTask) run() {
Loop:
	for {
		select {
		case msg := <-t.instance.Lines:
			kafka.SendToChan(t.topic, msg.Text)
		case <-t.close:
			break Loop
		}
	}

	fmt.Printf("tail topic:%v,path:%v stop\n", t.topic, t.path)
}

func (t *TailTask) Close() {
	close(t.close)
}
