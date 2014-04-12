// vim: tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab tw=72
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"os"
	"strconv"
	"strings"
	"errors"
)

type existsFunc func(string) bool

func osStatExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

var ignoreList = [...]string{"/", ".", "./", "..", "../"}

//var rootListing, _ = ioutil.ReadDir("/")
//var pwdListing, _ = ioutil.ReadDir(".")

func ignored(file string) bool {
	for _, val := range ignoreList {
		if file == val {
			return true
		}
	}
	return false
}

func longestFileEndIndex(line []rune, exists existsFunc) int {
	// Possible optimisations:
	// 1. this should start at the end - the longest substring first
	// 2. it could be a good strategy to list files and try to
	//    find a common prefix - if not just stop right there and then
	//    from / if it starts with / and if not from `pwd`
	//    do file listing from / and `pwd` only once
	//    need to consider relative dirs though which is annoying

	maxIndex := 0
	for i, _ := range line {
		slice := line[0 : i+1]
		file := string(slice)
		if !ignored(file) {
			if exists(file) {
				// TODO if this is not a dir, stop here
				maxIndex = i
			}
		}
	}
	return maxIndex
}

func longestFileInLine(line string, exists existsFunc) (firstCharIndex int, lastCharIndex int) {
	for searchStartIndex, _ := range line {
		searchSpace := []rune(line[searchStartIndex:len(line)])
		lastCharIndexInSlice := longestFileEndIndex(searchSpace, exists)
		lastCharIndexInLine := lastCharIndexInSlice + searchStartIndex
		if lastCharIndexInSlice > 0 && lastCharIndexInLine > lastCharIndex {
			lastCharIndex = lastCharIndexInLine
			firstCharIndex = searchStartIndex
		}
	}

	return firstCharIndex, lastCharIndex
}

func askUser() (requestedNumbers []string, err error) {
	fmt.Println()
	fmt.Print("to clipboard: ")
	ttyFile, err := os.Open("/dev/tty")
	if err != nil {
		return nil, err
	}
	defer ttyFile.Close()
	ttyReader := bufio.NewReader(ttyFile)
	s, err := ttyReader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return strings.Fields(s), nil
}

type PrintFunction func(string, string, int, int)

type Processor interface {
	processFile(file string) error
	processEnd() error
}

type NumbersGiven struct {
	clip      *bytes.Buffer
	fileCount *int
	numbers   []string
}

type AskForNumbers struct {
	clip  *bytes.Buffer
	files []string
}

func printShortFormat(fileCount *int) PrintFunction {
	return func(file string, line string, firstCharIndex int, lastCharIndex int) {
		fmt.Println(strconv.Itoa(*fileCount), file)
	}
}

func printLongFormat(fileCount *int) PrintFunction {
	return func(file string, line string, firstCharIndex int, lastCharIndex int) {
		fmt.Print("{")
		fmt.Print(strconv.Itoa(*fileCount))
		fmt.Print("} ")
		fmt.Print(line[:firstCharIndex])
		fmt.Print("{")
		fmt.Print(file)
		fmt.Print("}")
		fmt.Print(line[lastCharIndex+1:])
	}
}

func (ng *NumbersGiven) processEnd() error {
	writeToClipboard(ng.clip)
	return nil
}

func (ng *NumbersGiven) processFile(file string) error {
	// collect any file position arguments to copy to the
	// clipboard later
	for _, v := range ng.numbers {
		n, err := strconv.Atoi(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "nw: %s is not a number\n", v)
			return err
		}
		if n == *ng.fileCount {
			ng.clip.WriteString(file)
			ng.clip.WriteString(" ")
		}
	}
	return nil
}

func (afn *AskForNumbers) processEnd() error {
	if len(afn.files) == 0 {
		fmt.Println("nw: no files names found, NUMBERWANG!")
		return nil
	}
	requestedNumbers, err := askUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "nw: failed to read input: %s\n", err)
		return err
	}

	for _, n := range requestedNumbers {
		i, err := strconv.Atoi(n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "nw: %s is not a number\n", n)
			return err
		}
		if i <= 0 || i > len(afn.files) {
			fmt.Fprintf(os.Stderr, "nw: %s is not a valid choice\n", n)
			return errors.New("invalid choice")
		}
		afn.clip.WriteString(afn.files[i-1])
		afn.clip.WriteString(" ")
	}

	writeToClipboard(afn.clip)
	return nil
}

func (afn *AskForNumbers) processFile(file string) error {
	afn.files = append(afn.files, file)
	return nil
}

func writeToClipboard(buffer *bytes.Buffer) {
	clipboardOutput := buffer.String()
	if clipboardOutput != "" {
		clipboard.WriteAll(clipboardOutput)
		fmt.Printf("nw: wrote \"%s\" to clipboard\n", clipboardOutput)
	}
}

func main() {
	var fileCount int
	var clip bytes.Buffer

	short := flag.Bool("s", false, "short format, only display file names")
	flag.Parse()

	extraArgs := flag.Args()

	var processor Processor
	if len(extraArgs) == 0 {
		processor = &AskForNumbers{clip: &clip}
	} else {
		processor = &NumbersGiven{
			clip:      &clip,
			fileCount: &fileCount,
			numbers:   extraArgs,
		}
	}
	var printer PrintFunction
	if *short {
		printer = printShortFormat(&fileCount)
	} else {
		printer = printLongFormat(&fileCount)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		firstCharIndex, lastCharIndex := longestFileInLine(line, osStatExists)

		if lastCharIndex > 0 {
			fileCount++
			file := line[firstCharIndex : lastCharIndex+1]

			printer(file, line, firstCharIndex, lastCharIndex)
			err := processor.processFile(file)
			if err != nil {
				os.Exit(1)
			}
		} else if (!*short) {
			fmt.Print(line)
		}
	}

	err := processor.processEnd()
	if err != nil {
		os.Exit(1)
	}
}
