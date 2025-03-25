package gcfstructuredlogformatter

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogEntry(t *testing.T) {
	rows := []struct {
		description string
		input       logEntry
	}{
		{
			description: "Empty",
			input:       logEntry{},
		},
		{
			description: "All fields",
			input: logEntry{
				"message":                      "my-message",
				"severity":                     "my-severity",
				"logging.googleapis.com/trace": "my-trace",
				"labels": map[string]any{
					"label-1": "value-1",
					"label-2": "value-2",
				},
			},
		},
	}
	for rowIndex, row := range rows {
		t.Run(fmt.Sprintf("%d/%s", rowIndex, row.description), func(t *testing.T) {
			contents, err := json.Marshal(row.input)
			require.Nil(t, err)

			var testEntry logEntry
			err = json.Unmarshal(contents, &testEntry)
			require.Nil(t, err)

			assert.Equal(t, row.input, testEntry)
		})
	}
}

func TestFormat(t *testing.T) {
	logger := logrus.New()

	rows := []struct {
		description string
		input       *logrus.Entry
		output      logEntry
	}{
		{
			description: "Empty",
			input:       logrus.NewEntry(logger),
			output: logEntry{
				"message":  "",
				"severity": "Emergency", // logrus's 0th level is PanicLevel.
			},
		},
		{
			description: "Info",
			input: func() *logrus.Entry {
				e := logrus.NewEntry(logger)
				e.Message = "test"
				e.Level = logrus.InfoLevel

				return e
			}(),
			output: logEntry{
				"message":  "test",
				"severity": "Info",
			},
		},
		{
			description: "Warning",
			input: func() *logrus.Entry {
				e := logrus.NewEntry(logger)
				e.Message = "test"
				e.Level = logrus.WarnLevel

				return e
			}(),
			output: logEntry{
				"message":  "test",
				"severity": "Warning",
			},
		},
		{
			description: "Info with trace",
			input: func() *logrus.Entry {
				ctx := context.WithValue(context.Background(), ContextKeyTrace, "trace-1")
				e := logrus.NewEntry(logger).WithContext(ctx)
				e.Message = "test"
				e.Level = logrus.InfoLevel

				return e
			}(),
			output: logEntry{
				"message":                      "test",
				"severity":                     "Info",
				"logging.googleapis.com/trace": "trace-1",
			},
		},
		{
			description: "Info with bogus trace",
			input: func() *logrus.Entry {
				ctx := context.WithValue(context.Background(), ContextKeyTrace, 123456) // Not string.
				e := logrus.NewEntry(logger).WithContext(ctx)
				e.Message = "test"
				e.Level = logrus.InfoLevel

				return e
			}(),
			output: logEntry{
				"message":  "test",
				"severity": "Info",
			},
		},
		{
			description: "With logrus fields",
			input: func() *logrus.Entry {
				e := logrus.NewEntry(logger).
					WithFields(logrus.Fields{
						"key-1": "value-1",
						"nested": logrus.Fields{
							"key-2": "value-2",
						},
					})

				e.Message = "test"
				e.Level = logrus.InfoLevel

				return e
			}(),
			output: logEntry{
				"message": "test",
				"key-1":   "value-1",
				"nested": map[string]any{
					"key-2": "value-2",
				},
				"severity": "Info",
			},
		},
	}
	for rowIndex, row := range rows {
		t.Run(fmt.Sprintf("%d/%s", rowIndex, row.description), func(t *testing.T) {
			formatter := New()
			result, err := formatter.Format(row.input)
			require.Nil(t, err)
			if assert.NotNil(t, result) {
				var output logEntry
				err = json.Unmarshal(result, &output)
				require.Nil(t, err)
				assert.Equal(t, row.output, output)
			}
		})
	}
}
