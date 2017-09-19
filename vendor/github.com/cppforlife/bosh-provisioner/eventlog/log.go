package eventlog

import (
	"time"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

const logLogTag = "Log"

type Log struct {
	device Device
	logger boshlog.Logger
}

type LogEntry struct {
	Time int64 `json:"time"`

	Stage string   `json:"stage"`
	Task  string   `json:"task"`
	Tags  []string `json:"tags"`

	Total int `json:"total"`
	Index int `json:"index"`

	State    string `json:"state"`
	Progress int    `json:"progress"`

	// Might contain error key
	Data map[string]interface{} `json:"data,omitempty"`
}

type ErrorEntry struct {
	Time int64 `json:"time"`

	Body ErrorEntryBody `json:"error"`
}

type ErrorEntryBody struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func NewLog(device Device, logger boshlog.Logger) Log {
	return Log{device: device, logger: logger}
}

func (l Log) BeginStage(name string, total int) *Stage {
	return &Stage{
		log:   l,
		name:  name,
		total: total,
	}
}

func (l Log) WriteErr(err error) {
	entry := ErrorEntry{
		Time: time.Now().Unix(),

		Body: ErrorEntryBody{
			Code:    0,
			Message: err.Error(),
		},
	}

	l.logger.Error(logLogTag, "Error occurred: %s", err)

	writeErr := l.device.WriteErrorEntry(entry)
	if writeErr != nil {
		l.logger.Error(logLogTag, "Failed writing error entry %s", writeErr)
	}
}

func (l Log) WriteLogEntryNoErr(entry LogEntry) {
	err := l.device.WriteLogEntry(entry)
	if err != nil {
		l.logger.Error(logLogTag, "Failed writing log entry %s", err)
	}
}
