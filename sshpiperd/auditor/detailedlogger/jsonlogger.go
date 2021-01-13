package detailedlogger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type jsonLogger struct {
	Filename string
	LogLevel int
	Params   map[string]interface{}
}

func (l *jsonLogger) log(level int, message string, args ...interface{}) error {
	if level > l.LogLevel {
		return nil
	}
	message = fmt.Sprintf(message, args...)
	l.Params["timestamp"] = time.Now().UTC().Format("2006-01-02 15:04:05.000")
	l.Params["message"] = message
	data, err := json.Marshal(l.Params)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(l.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(append(data, '\n')); err != nil {
		return err
	}
	return nil
}

func (l *jsonLogger) clone() *jsonLogger {
	newLogger := &jsonLogger{
		Filename: l.Filename,
		LogLevel: l.LogLevel,
		Params:   map[string]interface{}{},
	}
	for k, v := range l.Params {
		newLogger.Params[k] = v
	}
	return newLogger
}

func newJsonLogger(feature string) (*jsonLogger, error) {
	l := &jsonLogger{
		Filename: "/dev/stdout",
		LogLevel: 1,
		Params: map[string]interface{}{
			"feature": feature,
		},
	}
	return l, nil
}
