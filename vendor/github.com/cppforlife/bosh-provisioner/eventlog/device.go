package eventlog

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Device interface {
	WriteLogEntry(LogEntry) error
	WriteErrorEntry(ErrorEntry) error
}

// JSONDevice writes events as JSON log entries
type JSONDevice struct {
	writer io.Writer
}

// TextDevice writes events in user friendly format
type TextDevice struct {
	writer io.Writer
}

func NewJSONDevice(writer io.Writer) JSONDevice {
	return JSONDevice{writer: writer}
}

func NewTextDevice(writer io.Writer) TextDevice {
	return TextDevice{writer: writer}
}

func (d JSONDevice) WriteLogEntry(entry LogEntry) error {
	return d.writeEntry(entry)
}

func (d JSONDevice) WriteErrorEntry(entry ErrorEntry) error {
	return d.writeEntry(entry)
}

func (d JSONDevice) writeEntry(entry interface{}) error {
	bytes, err := json.Marshal(entry)
	if err != nil {
		return bosherr.WrapError(err, "Marshalling entry")
	}

	bytes = append(bytes, []byte("\n")...)

	_, err = d.writer.Write(bytes)
	if err != nil {
		return bosherr.WrapError(err, "Writing entry")
	}

	return nil
}

func (d TextDevice) WriteLogEntry(entry LogEntry) error {
	_, err := fmt.Fprintf(
		d.writer,
		"%s %s > %s\n",
		strings.Title(entry.State),
		strings.ToLower(entry.Stage),
		entry.Task,
	)
	if err != nil {
		return bosherr.WrapError(err, "Writing log entry")
	}

	return nil
}

func (d TextDevice) WriteErrorEntry(entry ErrorEntry) error {
	_, err := fmt.Fprintf(d.writer, "\n-----\nError: %s\n-----\n", entry.Body.Message)
	if err != nil {
		return bosherr.WrapError(err, "Writing error entry")
	}

	return nil
}
