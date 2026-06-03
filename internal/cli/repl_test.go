package cli

import "testing"

func TestCleanInput(t *testing.T) {
	type testCase struct {
		input    string
		expected []string
	}

	cases := []testCase{
		{input: "  hello world  ", expected: []string{"hello", "world"}},
		{input: "  this is beautiful  ", expected: []string{"this", "is", "beautiful"}},
		{input: "", expected: []string{}},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("expected %v, got %v", c.expected, actual)
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("expected %v, got %v", expectedWord, word)
				return
			}
		}
	}
}
