package commands

import (
	"testing"
)

func TestSortCommandsByOrder(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]Command
		expected []string // Just check the order of names
	}{
		{
			name: "sorts commands by order field",
			input: map[string]Command{
				"exit": {Name: "exit", Order: 3},
				"help": {Name: "help", Order: 1},
				"map":  {Name: "map", Order: 2},
			},
			expected: []string{"help", "map", "exit"},
		},
		{
			name: "handles duplicate names",
			input: map[string]Command{
				"help": {Name: "help", Order: 1},
				"h":    {Name: "help", Order: 1}, // duplicate
			},
			expected: []string{"help"}, // only one should appear
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortCommandByOrder(tt.input)
			// verify the order
			for i, cmd := range result {
				if cmd.Name != tt.expected[i] {
					t.Errorf("Position %d: expected %s, got %s", i, tt.expected[i], cmd.Name)
				}
			}
		})
	}
}
