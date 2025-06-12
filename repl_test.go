package main

import "testing"

func TestCleanUp(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " shit happens ",
			expected: []string{"shit", "happens"},
		},
		{
			input:    " you will succeed ",
			expected: []string{"you", "will", "succeed"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Lenght of actual are not equals to expected")
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("words are not the same in the slice")
			}
		}
	}
}
