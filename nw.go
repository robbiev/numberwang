//http://stackoverflow.com/questions/8757389/reading-file-line-by-line-in-go
package main

import (
	"fmt"
	"os"
	"strconv"
	"bufio"
	"bytes"
	"github.com/atotto/clipboard"
)

func longestFile(line []rune) int {
	max := 0
	for i, _ := range line {
		slice := line[0:i]
		file := string(slice)
		if _, err := os.Stat(file); err == nil {
			max = i
		}
	}
	return max
}

func main() {
	var clip bytes.Buffer
	argsWithoutProg := os.Args[1:]

	reader := bufio.NewReader(os.Stdin)

	i := 0

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			// check here if err == io.EOF
			break
		}

		longest := 0
		start := 0
		for i, _ := range line {
			search := []rune(line[i:len(line)])
			found := longestFile(search)
			pos := found + i + 1
			if found > 0 && pos > longest {
				longest = pos
				start = i
			}
		}

		i = i + 1

		if longest > 0 {
			fmt.Println(strconv.Itoa(i), line[start:longest-1])
			for _, v := range argsWithoutProg {
				n, _ := strconv.Atoi(v)
				if n == i {
					clip.WriteString(line[start:longest-1] + " ")
				}
			}
		} else {
			i = i - 1
			fmt.Print(line)
		}
	}
	clipboard.WriteAll(clip.String())
}
