package gcfstructuredlogformatter

import (
	"encoding/json"

	"cloud.google.com/go/logging"
	"github.com/sirupsen/logrus"
)

// ContextKey is the type for the context key.
// The Go docs recommend not using any built-in type for context keys in order
// to ensure that there are no collisions:
//    https://golang.org/pkg/context/#WithValue
type ContextKey string

// ContextKey constants.
const (
	ContextKeyTrace ContextKey = "trace" // This is the key for the trace identifier.
)

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

// logEntry is an abbreviated version of the Google "structured logging" data structure.
type logEntry struct {
	Message  string            `json:"message"`
	Severity string            `json:"severity,omitempty"`
	Trace    string            `json:"logging.googleapis.com/trace,omitempty"`
	Labels   map[string]string `json:"labels,omitempty"`
}

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

	newEntry := logEntry{
		Message:  entry.Message,
		Severity: severity.String(),
		Labels:   map[string]string{},
	}
	if entry.Context != nil {
		if v, okay := entry.Context.Value(ContextKeyTrace).(string); okay {
			newEntry.Trace = v
		}
	}
	for key, value := range f.Labels {
		newEntry.Labels[key] = value
	}

	contents, err := json.Marshal(newEntry)
	if err != nil {
		return nil, err
	}
	return append(contents, []byte("\n")...), nil
}
