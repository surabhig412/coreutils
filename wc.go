package main

import (
	bytesp "bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

var (
	line       = flag.Bool("l", false, "The number of lines in each input file is written to the standard output.")
	words      = flag.Bool("w", false, "The number of words in each input file is written to the standard output.")
	bytes      = flag.Bool("c", false, "The number of bytes in each input file is written to the standard output.  This will cancel out any prior usage of the -m option.")
	characters = flag.Bool("m", false, "The number of characters in each input file is written to the standard output.  If the current locale does not support multibyte characters, this is equivalent to the -c option.  This will cancel out any prior usage of the -c option.")
)

var wcMap map[int]*fileStat

type fileStat struct {
	Name           string
	LineCount      int64
	WordCount      int64
	ByteCount      int64
	CharacterCount int64
	Error          error
}

func calculateOptions() int {
	totalOptions := 1
	if *line {
		totalOptions++
	}
	if *words {
		totalOptions++
	}
	if *bytes {
		totalOptions++
	}
	if *characters {
		totalOptions++
	}
	if totalOptions == 1 {
		*line, *words, *bytes = true, true, true
	}
	return totalOptions
}
func setReader(options []string) [][]byte {
	var fileBytes [][]byte
	if len(options) == 0 {
		stdinBytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			f := &fileStat{Name: "", Error: errors.New("wc: " + err.Error())}
			wcMap[0] = f
		}
		if wcMap[0] == nil {
			wcMap[0] = &fileStat{}
		}
		fileBytes = append(fileBytes, stdinBytes)
	}
	for i, v := range options {
		fileByte, err := ioutil.ReadFile(v)
		if err != nil {
			f := &fileStat{Name: v, Error: errors.New("wc: " + err.Error())}
			wcMap[i] = f
			return fileBytes
		}
		wcMap[i] = &fileStat{Name: v}
		fileBytes = append(fileBytes, fileByte)
	}
	return fileBytes

}

func countByNewLineDelim(i int, b []byte) {
	var count int
	buff := bytesp.NewBuffer(b)
	for {
		_, err := buff.ReadBytes(byte('\n'))
		switch {
		case err == io.EOF:
			wcMap[i].LineCount = int64(count)
			return

		case err != nil:
			f := &fileStat{Name: "", Error: err}
			wcMap[i] = f
			return
		}
		count++
	}
}
func countLines(wcBytes [][]byte) {
	for i, wcByte := range wcBytes {
		countByNewLineDelim(i, wcByte)
	}

}

func countBytes(wcBytes [][]byte) {
	for i, wcByte := range wcBytes {
		wcMap[i].ByteCount = int64(len(wcByte))
	}
}

func countWords(wcBytes [][]byte) {
	for i, wcByte := range wcBytes {
		wcMap[i].WordCount = int64(len(strings.Fields(string(wcByte))))
	}
}

func countCharacters(wcBytes [][]byte) {
	for i, wcByte := range wcBytes {
		wcMap[i].CharacterCount = int64(utf8.RuneCount(wcByte))
	}
}

func calculateMetrics() (lineSum, wordSum, characterSum, byteSum int64) {
	var keys []int
	for k := range wcMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, keyIndex := range keys {
		if wcMap[keyIndex].Error != nil {
			io.Copy(os.Stdout, strings.NewReader(fmt.Sprintf("%s\n", wcMap[keyIndex].Error.Error())))
		} else {
			var message string
			if *line {

				message = fmt.Sprintf("\t%d", wcMap[keyIndex].LineCount)
				lineSum = lineSum + wcMap[keyIndex].LineCount
			}
			if *words {
				if message != "" {
					message = fmt.Sprintf("\t%s\t%d", message, wcMap[keyIndex].WordCount)
				} else {
					message = fmt.Sprintf("\t%d", wcMap[keyIndex].WordCount)
				}
				wordSum = wordSum + wcMap[keyIndex].WordCount
			}
			if *bytes {
				if message != "" {
					message = fmt.Sprintf("\t%s\t%d", message, wcMap[keyIndex].ByteCount)
				} else {
					message = fmt.Sprintf("\t%d", wcMap[keyIndex].ByteCount)
				}
				byteSum = byteSum + wcMap[keyIndex].ByteCount

			}
			if *characters {
				if message != "" {
					message = fmt.Sprintf("\t%s\t%d", message, wcMap[keyIndex].CharacterCount)
				} else {
					message = fmt.Sprintf("\t%d", wcMap[keyIndex].CharacterCount)
				}
				characterSum = characterSum + wcMap[keyIndex].CharacterCount

			}
			io.Copy(os.Stdout, strings.NewReader(fmt.Sprintf("%s\t%s\n", message, wcMap[keyIndex].Name)))
		}
	}
	return lineSum, wordSum, characterSum, byteSum
}

func printAggregateMetrics(lineSum, wordSum, characterSum, byteSum int64) {
	var totalMessage string
	if *line {
		totalMessage = fmt.Sprintf("\t%d", lineSum)
	}
	if *words {
		if totalMessage != "" {
			totalMessage = fmt.Sprintf("\t%s\t%d", totalMessage, wordSum)
		} else {
			totalMessage = fmt.Sprintf("\t%d", wordSum)
		}
	}
	if *bytes {
		if totalMessage != "" {
			totalMessage = fmt.Sprintf("\t%s\t%d", totalMessage, byteSum)
		} else {
			totalMessage = fmt.Sprintf("\t%d", byteSum)
		}
	}
	if *characters {
		if totalMessage != "" {
			totalMessage = fmt.Sprintf("\t%s\t%d", totalMessage, characterSum)
		} else {
			totalMessage = fmt.Sprintf("\t%d", characterSum)
		}
	}
	io.Copy(os.Stdout, strings.NewReader(fmt.Sprintf("%s\ttotal\n", totalMessage)))

}
func main() {
	flag.Parse()

	wcMap = make(map[int]*fileStat)
	totalOptions := calculateOptions()
	wcBytes := setReader(os.Args[totalOptions:])

	if *line {
		countLines(wcBytes)
	}
	if *words {
		countWords(wcBytes)
	}
	if *bytes {
		countBytes(wcBytes)
	}
	if *characters {
		countCharacters(wcBytes)
	}
	lineSum, wordSum, characterSum, byteSum := calculateMetrics()

	if len(os.Args[totalOptions:]) > 1 {
		printAggregateMetrics(lineSum, wordSum, characterSum, byteSum)
	}
}
