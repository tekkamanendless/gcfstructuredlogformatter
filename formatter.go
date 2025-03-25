package gcfstructuredlogformatter

import (
	"encoding/json"
	"fmt"

	"cloud.google.com/go/logging"
	"github.com/sirupsen/logrus"
)

// ContextKey is the type for the context key.
// The Go docs recommend not using any built-in type for context keys in order
// to ensure that there are no collisions:
//
//	https://golang.org/pkg/context/#WithValue
type ContextKey string

// ContextKeyTrace defines the key for the trace identifier.
const ContextKeyTrace ContextKey = "trace"

// logrusToGoogleSeverityMap maps a logrus level to a Google severity.
var logrusToGoogleSeverityMap = map[logrus.Level]logging.Severity{
	logrus.PanicLevel: logging.Emergency,
	logrus.FatalLevel: logging.Alert,
	logrus.ErrorLevel: logging.Error,
	logrus.WarnLevel:  logging.Warning,
	logrus.InfoLevel:  logging.Info,
	logrus.DebugLevel: logging.Debug,
	logrus.TraceLevel: logging.Default,
}

// Formatter is the logrus formatter.
type Formatter struct {
	Labels map[string]string // This is an optional map of additional "labels".
}

type logEntry = map[string]any

// Define log entry field keys used for JSON marshaling.
const (
	fieldKeySeverity = "severity"
	fieldKeyTrace    = "logging.googleapis.com/trace"
	fieldKeyLabels   = "labels"
	fieldKeyMessage  = "message"
)

// New creates a new formatter.
func New() *Formatter {
	f := &Formatter{
		Labels: map[string]string{},
	}
	return f
}

// Levels are the available logging levels.
func (f *Formatter) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
}

// Format an entry.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	severity := logging.Default
	if value, okay := logrusToGoogleSeverityMap[entry.Level]; okay {
		severity = value
	}

	newEntry := logEntry{}

	for key, value := range entry.Data {
		newEntry[key] = value
	}

	newEntry[fieldKeySeverity] = severity.String()
	newEntry[fieldKeyMessage] = entry.Message

	if entry.Context != nil {
		if v, okay := entry.Context.Value(ContextKeyTrace).(string); okay {
			newEntry[fieldKeyTrace] = v
		}
	}

	if len(f.Labels) > 0 {
		newEntry[fieldKeyLabels] = f.Labels
	}

	rawJSON, err := json.Marshal(newEntry)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal log entry to JSON: %w", err)
	}

	return append(rawJSON, []byte("\n")...), nil
}
