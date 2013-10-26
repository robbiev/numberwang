// vim: tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab tw=72
//http://stackoverflow.com/questions/8757389/reading-file-line-by-line-in-go
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/atotto/clipboard"
	"os"
	"strconv"
)

func longestFileEndIndex(line []rune) int {
	maxIndex := 0
	for i, _ := range line {
		slice := line[0:i]
		file := string(slice)
		if file != "/" && file != "." && file != "./" && file != ".." && file != "../" {
			if _, err := os.Stat(file); err == nil {
				maxIndex = i - 1
			}
		}
	}
	return maxIndex
}

func longestFileInLine(line string) (firstCharIndex int, lastCharIndex int) {
	for searchStartIndex, _ := range line {
		searchSpace := []rune(line[searchStartIndex:len(line)])
		lastCharIndexInSlice := longestFileEndIndex(searchSpace)
		lastCharIndexInLine := lastCharIndexInSlice + searchStartIndex
		if lastCharIndexInSlice > 0 && lastCharIndexInLine > lastCharIndex {
			lastCharIndex = lastCharIndexInLine
			firstCharIndex = searchStartIndex
		}
	}

	return firstCharIndex, lastCharIndex
}

func main() {
	var clip bytes.Buffer

	argsWithoutProg := os.Args[1:]

	reader := bufio.NewReader(os.Stdin)

	fileCount := 0

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			// check here if err == io.EOF
			break
		}

		firstCharIndex, lastCharIndex := longestFileInLine(line)

		if lastCharIndex > 0 {
			fileCount++
			file := line[firstCharIndex : lastCharIndex+1]

			fmt.Println(strconv.Itoa(fileCount), file)

			// collect any file position arguments to copy to the
			// clipboard later
			for _, v := range argsWithoutProg {
				n, _ := strconv.Atoi(v)
				if n == fileCount {
					clip.WriteString(file)
					clip.WriteString(" ")
				}
			}
		} else {
			fmt.Print(line)
		}
	}

	clipboardOutput := clip.String()
	if clipboardOutput != "" {
		clipboard.WriteAll(clipboardOutput)
	}
}
