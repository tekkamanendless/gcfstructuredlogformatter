package gcfstructuredlogformatter

import (
	"encoding/json"

	"cloud.google.com/go/logging"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

// Keys defined https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
const (
	// TraceKey is the key for the trace identifier.
	TraceKey = "logging.googleapis.com/trace"
	// SpanKey is the key for the span identifier.
	SpanKey = "logging.googleapis.com/spanId"
	// SeverityKey is the key for the severity.
	SeverityKey = "severity"
	// MessageKey is the key for the message.
	MessageKey = "message"
	// LabelsKey is the key for the labels.
	LabelsKey = "logging.googleapis.com/labels"
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

// New creates a new formatter.
func New() *Formatter {
	f := &Formatter{
		Labels: map[string]string{},
	}
	return f
}

// AddLabel adds a label to the formatter.
func (f *Formatter) AddLabel(key, value string) {
	f.Labels[key] = value
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

	mapEntry := map[string]interface{}{}
	mapEntry[SeverityKey] = severity.String()
	mapEntry[MessageKey] = entry.Message

	if entry.Context != nil {
		// try to get the trace id from the context
		span := trace.SpanFromContext(entry.Context)
		spanContext := span.SpanContext()
		if spanContext.IsValid() {
			mapEntry[TraceKey] = spanContext.TraceID().String()
			mapEntry[SpanKey] = spanContext.SpanID().String()
		}
	}
	if len(f.Labels) > 0 {
		labels := map[string]string{}
		for key, value := range f.Labels {
			labels[key] = value
		}
		mapEntry[LabelsKey] = labels
	}

	for key, value := range entry.Data {
		mapEntry[key] = value
	}
	contents, err := json.Marshal(mapEntry)
	if err != nil {
		return nil, err
	}
	return append(contents, []byte("\n")...), nil
}
