package gcfstructuredlogformatter

import (
	"encoding/json"

	"cloud.google.com/go/logging"
	"github.com/sirupsen/logrus"
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

// GoogleCloudFunctionFormatter is the logrus formatter.
type GoogleCloudFunctionFormatter struct {
	Labels map[string]string // This is an optional map of additional "labels".
}

// logEntry is an abbreviated version of the Google "structured logging" data structure.
type logEntry struct {
	Message  string            `json:"message"`
	Severity string            `json:"severity"`
	Labels   map[string]string `json:"labels"`
}

// New creates a new formatter.
func New() *GoogleCloudFunctionFormatter {
	f := &GoogleCloudFunctionFormatter{
		Labels: map[string]string{},
	}
	return f
}

// Levels are the available logging levels.
func (f *GoogleCloudFunctionFormatter) Levels() []logrus.Level {
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

// Fire sends an entry.
func (f *GoogleCloudFunctionFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	severity := logging.Default
	if value, okay := logrusToGoogleSeverityMap[entry.Level]; okay {
		severity = value
	}

	newEntry := logEntry{
		Message:  entry.Message,
		Severity: severity.String(),
		Labels:   map[string]string{},
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
