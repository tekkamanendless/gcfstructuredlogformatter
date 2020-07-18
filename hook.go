package gcfstructuredlogformatter

import (
	"encoding/json"
	"net/http"
	"os"

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
	ExecutionID  string            // This is the execution ID, as found in the HTTP header `Function-Execution-Id`.
	FunctionName string            // This is the function name, as found in the `FUNCTION_NAME` environment variable.
	Labels       map[string]string // This is an optional map of additional "labels".
}

// logEntry is an abbreviated version of the Google "structured logging" data structure.
type logEntry struct {
	message  string            `json:"message"`
	severity string            `json:"severity"`
	labels   map[string]string `json:"labels"`
}

// New creates a new formatter.
func New() *GoogleCloudFunctionFormatter {
	f := &GoogleCloudFunctionFormatter{
		FunctionName: os.Getenv("FUNCTION_NAME"),
		Labels:       map[string]string{},
	}
	return f
}

// NewForRequest creates a new formatter that will include the "execution ID" of the request
// as a label with each log message.
func NewForRequest(r *http.Request) *GoogleCloudFunctionFormatter {
	f := New()

	// If we can get the execution ID, then use it when we fire the log messages.
	f.ExecutionID = r.Header.Get("Function-Execution-Id")

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
		message:  entry.Message,
		severity: severity.String(),
		labels:   map[string]string{},
	}
	for key, value := range f.Labels {
		newEntry.labels[key] = value
	}
	if f.ExecutionID != "" {
		newEntry.labels["execution_id"] = f.ExecutionID
	}
	if f.FunctionName != "" {
		newEntry.labels["function_name"] = f.FunctionName
	}

	return json.Marshal(newEntry)
}

