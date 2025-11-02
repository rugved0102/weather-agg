package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type LogEntry struct {
	Timestamp string      `json:"timestamp"`
	Level     string      `json:"level"`
	Event     string      `json:"event"`
	Module    string      `json:"module"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}

var logger = log.New(os.Stdout, "", 0)

func logEvent(level, event, module, message string, data interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Event:     event,
		Module:    module,
		Message:   message,
		Data:      data,
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}

	logger.Println(string(jsonBytes))
}

func Info(event, module, message string, data interface{}) {
	logEvent("INFO", event, module, message, data)
}

func Error(event, module, message string, data interface{}) {
	logEvent("ERROR", event, module, message, data)
}

func Debug(event, module, message string, data interface{}) {
	logEvent("DEBUG", event, module, message, data)
}
