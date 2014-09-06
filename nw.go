// vim: tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab tw=72
package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"os"
	"strconv"
	"strings"
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

func askUser(question string) (requestedNumbers []string, err error) {
	fmt.Println()
	fmt.Print(question)
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
	invert    bool
}

type AskForNumbers struct {
	clip   *bytes.Buffer
	files  []string
	invert bool
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
	for _, v := range ng.numbers {
		n, err := strconv.Atoi(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "nw: %s is not a number\n", v)
			return err
		}
		found := (n == *ng.fileCount)
		if found {
			if found != ng.invert {
				ng.clip.WriteString(file)
				ng.clip.WriteString(" ")
			}
			return nil
		}
	}

	// not found in given numbers, so always a match if inverted
	if ng.invert {
		ng.clip.WriteString(file)
		ng.clip.WriteString(" ")
	}
	return nil
}

func (afn *AskForNumbers) processEnd() error {
	if len(afn.files) == 0 {
		fmt.Println("nw: no file names found, NUMBERWANG!")
		return nil
	}
	question := "to clipboard: "
	if afn.invert {
		question = "NOT " + question
	}
	requestedNumbers, err := askUser(question)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nw: failed to read input: %s\n", err)
		return err
	}

	reqIndexes := make(map[int]struct{})
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
		reqIndexes[i-1] = struct{}{}
	}

	for i, file := range afn.files {
		if _, exists := reqIndexes[i]; exists != afn.invert {
			afn.clip.WriteString(file)
			afn.clip.WriteString(" ")
		}
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
	invert := flag.Bool("i", false, "copy the inverse of the selection to the clipboard")
	flag.Parse()

	extraArgs := flag.Args()

	var processor Processor
	if len(extraArgs) == 0 {
		processor = &AskForNumbers{clip: &clip, invert: *invert}
	} else {
		processor = &NumbersGiven{
			clip:      &clip,
			fileCount: &fileCount,
			numbers:   extraArgs,
			invert:    *invert,
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
		} else if !*short {
			fmt.Print(line)
		}
	}

	err := processor.processEnd()
	if err != nil {
		os.Exit(1)
	}
}
