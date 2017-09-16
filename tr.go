package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	deleteChar          = flag.Bool("d", false, "Delete characters in string1 from the input.")
	oneOccurrenceChar   = flag.Bool("s", false, "Squeeze multiple occurrences of the characters listed in the last operand (either string1 or string2) in the input into a single instance of the character.  This occurs after all deletion and translation is completed.")
	disableOutputBuffer = flag.Bool("u", false, "Guarantee that any output is unbuffered.")
	str1                string
	str2                string
	noFlag              bool
)

func main() {
	flag.Parse()
	noFlag = true
	checkUsage(os.Args)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var result, output string
		if *deleteChar || noFlag {
			for _, value := range scanner.Bytes() {
				index := strings.Index(str1, string(value))
				if *deleteChar {
					if index == -1 {
						result += string(value)
					}
				} else if noFlag {
					if index > -1 {
						if len(str2) <= index {
							index = len(str2) - 1
						}
						result += string(str2[index])
					} else {
						result += string(value)
					}
				}
			}
		} else {
			result = scanner.Text()
		}

		if *oneOccurrenceChar {
			str := str1
			if *deleteChar {
				str = str2
			}
			for i := 0; i < len(result); i++ {
				index := strings.Index(str, string(result[i]))
				if len(output) == 0 {
					output += string(result[i])
				} else if index != -1 && output[len(output)-1] == result[i] {
					continue
				} else {
					output += string(result[i])
				}
			}
		} else {
			output = result
		}
		io.Copy(os.Stdout, strings.NewReader(output+"\n"))
	}
}

func checkUsage(args []string) {
	if len(args) < 2 {
		usageError()
	}
	stringArgs := 2
	stringArgsCopy := 99
	squeeze := false
	deleteChar := false
	for _, arg := range args[1:] {
		if flagArg(arg) {
			noFlag = false
			if arg == "-d" {
				deleteChar = true
				stringArgs = 1
			}
			if arg == "-s" {
				squeeze = true
				stringArgs = 1
			}
		} else {
			if squeeze && deleteChar && stringArgs == 1 {
				stringArgs = 2
			}
			if stringArgsCopy > stringArgs {
				stringArgsCopy = stringArgs
			}
			if stringArgsCopy > 0 {
				stringArgsCopy--
			} else {
				usageError()
			}
		}
	}
	if stringArgsCopy > 0 {
		usageError()
	}
	if stringArgs == 1 {
		str1 = args[len(args)-1]
	} else {
		str1 = args[len(args)-2]
		str2 = args[len(args)-1]
	}
}

func flagArg(arg string) bool {
	splitArg := strings.Split(arg, "-")
	if len(splitArg) == 2 && flag.Lookup(splitArg[1]) != nil {
		return true
	}
	return false
}

func usageError() {
	usageStr := "\nusage: tr string1 string2 \ntr -d string1 \ntr -s string1 \ntr -d -s string1 string2"
	fmt.Println("Please check usage. Use -h flag for help.", usageStr)
	os.Exit(1)
}
