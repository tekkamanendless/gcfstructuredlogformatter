package gcfstructuredlogformatter

import (
	"encoding/json"
	"fmt"
	"testing"

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
				Message:  "my-message",
				Severity: "my-severity",
				Labels: map[string]string{
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
