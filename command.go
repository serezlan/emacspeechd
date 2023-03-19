package main

import (
	"strings"
)

type Command struct {
	Name string
	Args []string
}

// Parse returns parsed command from emacspeak.
func Parse(s string) Command {
	i := 0

	for {
		if i == len(s)-1 || s[i] == ' ' {
			break
		}
		i++
	}
	trimmedInput := strings.TrimSuffix(s, "\n")

	result := Command{
		Name: trimmedInput[:i],
	}
	if len(s)-1 == i {
		return result
	}
	args := strings.TrimSuffix(s[i+1:], "\n")

	if args[0] == '{' && args[len(args)-1] == '}' {
		result.Args = []string{args[1 : len(args)-1]}
	} else {
		result.Args = strings.Split(args, " ")
	}

	return result
}
