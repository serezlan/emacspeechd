package main

import (
	"fmt"
	"testing"
)

type TestCase struct {
	title       string
	input       string
	expectation Command
}

var tests = []TestCase{
	{
		title:       "command without any argument",
		input:       "d\n",
		expectation: Command{Name: "d"},
	},
	{
		title:       "command with one argument",
		input:       "p /home/jameswebb/play.wav",
		expectation: Command{Name: "p", Args: []string{"/home/jameswebb/play.wav"}},
	},
	{
		title:       "command with multiple arguments",
		input:       "tts_set_rate 100 10\n",
		expectation: Command{Name: "tts_set_rate", Args: []string{"100", "10"}},
	},
	{
		title:       "command with arguments within braces",
		input:       "c {C-x C-f }\n",
		expectation: Command{Name: "c", Args: []string{"C-x C-f "}},
	},
}

func TestCommandParse(t *testing.T) {
	for _, test := range tests {
		var actual = Parse(test.input)
		if actual.Name != test.expectation.Name || !equalArgs(actual.Args, test.expectation.Args) {
			t.Errorf("wants %v\ngot %v\non %s", test.expectation, actual, test.title)
		}
	}
}

func equalArgs(a, b []string) bool {
	if len(a) != len(b) {
		fmt.Printf("length is different %d != %d\n", len(a), len(b))
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			fmt.Printf("uchiha diff content '%v' != '%v'\n", a[i], b[i])
			return false
		}
	}
	return true
}
