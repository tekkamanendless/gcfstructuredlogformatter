package gcfstructuredlogformatter

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	logger := logrus.New()
	rows := []struct {
		description string
		input       *logrus.Entry
		output      []byte
	}{
		{
			description: "Empty Entry",
			input: func() *logrus.Entry {
				e := logrus.NewEntry(logger)
				e.Level = logrus.InfoLevel
				return e
			}(),
			output: []byte(`{"message":"","severity":"Info"}` + "\n"),
		},
		{
			description: "Info Entry",
			input: func() *logrus.Entry {
				e := logger.WithFields(logrus.Fields{"prop": "value"})
				e.Message = "test"
				e.Level = logrus.InfoLevel
				return e
			}(),
			output: []byte(`{"message":"test","prop":"value","severity":"Info"}` + "\n"),
		},
	}

	for _, row := range rows {
		t.Run(row.description, func(t *testing.T) {
			formatter := New()
			result, err := formatter.Format(row.input)
			require.Nil(t, err)
			assert.Equal(t, row.output, result)
		})
	}
}

func TestFormatWithLabels(t *testing.T) {
	logger := logrus.New()
	rows := []struct {
		description string
		input       *logrus.Entry
		output      []byte
	}{

		{
			description: "Entry with Labels",
			input: func() *logrus.Entry {
				e := logger.WithFields(logrus.Fields{"prop": "value"})
				e.Message = "test"
				e.Level = logrus.InfoLevel
				return e
			}(),
			output: []byte(`{"logging.googleapis.com/labels":{"key":"value"},"message":"test","prop":"value","severity":"Info"}` + "\n"),
		},
	}

	for _, row := range rows {
		t.Run(row.description, func(t *testing.T) {
			formatter := New()
			formatter.AddLabel("key", "value")
			result, err := formatter.Format(row.input)
			require.Nil(t, err)
			assert.Equal(t, row.output, result)
		})
	}
}
