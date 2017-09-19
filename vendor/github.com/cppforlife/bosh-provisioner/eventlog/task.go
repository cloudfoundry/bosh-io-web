package eventlog

import (
	"time"
)

type Task struct {
	log Log

	stageName string
	name      string

	total int
	index int
}

func (t Task) Start() {
	entry := LogEntry{
		Time: time.Now().Unix(),

		Stage: t.stageName,
		Task:  t.name,

		Total: t.total,
		Index: t.index,

		State:    "started",
		Progress: 0,
	}

	t.log.WriteLogEntryNoErr(entry)
}

func (t Task) End(err error) error {
	entry := LogEntry{
		Time: time.Now().Unix(),

		Stage: t.stageName,
		Task:  t.name,

		Total: t.total,
		Index: t.index,

		State:    "finished",
		Progress: 100,
	}

	if err != nil {
		entry.State = "failed"
		entry.Data = map[string]interface{}{
			"error": err.Error(),
		}
	}

	t.log.WriteLogEntryNoErr(entry)

	return err
}
