// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// Logger provides structured logging
type Logger struct {
	level  string
	format string
}

// New creates a new logger
func New(level, format string) *Logger {
	return &Logger{
		level:  level,
		format: format,
	}
}

// LogEntry represents a log entry
type LogEntry struct {
	Time    string                 `json:"time"`
	Level   string                 `json:"level"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	l.log("info", message, fields...)
}

// Error logs an error message
func (l *Logger) Error(message string, err error, fields ...map[string]interface{}) {
	if len(fields) == 0 {
		fields = []map[string]interface{}{{"error": err.Error()}}
	} else {
		fields[0]["error"] = err.Error()
	}
	l.log("error", message, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	l.log("warn", message, fields...)
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	if l.level == "debug" {
		l.log("debug", message, fields...)
	}
}

func (l *Logger) log(level, message string, fields ...map[string]interface{}) {
	entry := LogEntry{
		Time:    time.Now().Format(time.RFC3339),
		Level:   level,
		Message: message,
	}

	if len(fields) > 0 {
		entry.Fields = fields[0]
	}

	if l.format == "json" {
		data, err := json.Marshal(entry)
		if err != nil {
			log.Printf("Failed to marshal log entry: %v", err)
			return
		}
		fmt.Fprintf(os.Stderr, "%s\n", data)
	} else {
		if entry.Fields != nil {
			log.Printf("[%s] %s %v", level, message, entry.Fields)
		} else {
			log.Printf("[%s] %s", level, message)
		}
	}
}

