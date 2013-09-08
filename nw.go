//http://stackoverflow.com/questions/8757389/reading-file-line-by-line-in-go
package main

import (
    "fmt"
    "os"
//    "io"
    "strconv"
//    "log"
    "bufio"
    "github.com/atotto/clipboard"
    "bytes"
)

func longestFile(line []rune) int {
    max := 0
    for i, _ := range line {
        slice := line[0:i]
        file := string(slice)
        if _, err := os.Stat(file); err == nil {
//          log.Println("file exists")
//          log.Println(file)
          max = i
        }
    }
    return max
}

func main() {
    var clip bytes.Buffer
    argsWithoutProg := os.Args[1:]
    //fmt.Println(argsWithoutProg)
    reader := bufio.NewReader(os.Stdin)

    i := 0

    for {
        line, err := reader.ReadString('\n')

        if err != nil {
            // You may check here if err == io.EOF
            //if err == io.EOF {
            //}
            break
        }

        //file := make([]rune, len(line))
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
            //fmt.Println(start)
            //fmt.Println(longest)
            fmt.Println(strconv.Itoa(i), line[start:longest - 1])
            for _, v := range argsWithoutProg {
              n, _ := strconv.Atoi(v)
              if n == i  {
                clip.WriteString(line[start:longest-1] + " ")
                //clipboard.WriteAll(line[start:longest-1])
              }
            }
            //err:=os.Setenv("NW"+ strconv.Itoa(i), line[start:longest -1])
            if err != nil {
              fmt.Println("ERROR")
            }
        } else {
            i = i - 1
            fmt.Print(line)
        }
    }
    clipboard.WriteAll(clip.String())
}
