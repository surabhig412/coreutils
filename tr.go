package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) <= 2 {
		usageError()
		return
	}
	str1 := os.Args[1]
	str2 := os.Args[2]
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var result string
		for _, value := range scanner.Bytes() {
			index := strings.Index(str1, string(value))
			if index > -1 {
				if len(str2) <= index {
					index = len(str2) - 1
				}
				result += string(str2[index])
			} else {
				result += string(value)
			}
		}
		io.Copy(os.Stdout, strings.NewReader(result))
	}
}

func usageError() {
	fmt.Println("Please check usage. Use -h flag for help.")
}
