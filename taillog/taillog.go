package taillog

import (
	"github.com/hpcloud/tail"
	"log-agent/kafka"
)

type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
}

func NewTailTask(path, topic string) *TailTask {
	task := &TailTask{
		path:  path,
		topic: topic,
	}
	task.init()
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
	for {
		select {
		case msg := <-t.instance.Lines:
			kafka.SendToChan(t.topic, msg.Text)
		}
	}
}
