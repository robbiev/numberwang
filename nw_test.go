// vim: tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab tw=72
package main

import (
	"testing"
)

func TestLastCharacterGetsConsidered(t *testing.T) {
	line := "blah/.gitkeep"
	start, end := longestFileInLine(line, func(file string) bool {
		return file == ".git" || file == line
	})
	result := line[start : end+1]

	if result != line {
		t.Errorf("Expected: '%s', found: '%s'", line, result)
	}
}
