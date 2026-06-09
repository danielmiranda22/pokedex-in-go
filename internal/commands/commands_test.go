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

func TestShouldCatchPokemon(t *testing.T) {
	tests := []struct {
		name      string
		randomVal int
		expected  bool
	}{
		{"catches when random > 40", 50, true},
		{"escapes when random < 40", 30, false},
		{"escapes when random = 40", 40, false},
		{"catches at boundary", 41, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldCatchPokemon(tt.randomVal)
			if result != tt.expected {
				t.Errorf("Expected %v but got %v", tt.expected, result)
			}
		})
	}
}
