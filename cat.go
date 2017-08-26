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
	numberNonBlankOutputLines = flag.Bool("b", false, "Number the non-blank output lines, starting at 1")
	displayDollarAtEnd        = flag.Bool("e", false, "Display non-printing characters (see the -v option), and display a dollar sign (`$') at the end of each line.")
	displayNonPrintingChars   = flag.Bool("v", false, "Display non-printing characters so they are visible.")
	numberOutputLines         = flag.Bool("n", false, "Number the output lines, starting at 1.")
	singleSpacedOutput        = flag.Bool("s", false, "Squeeze multiple adjacent empty lines, causing the output to be single spaced.")
	displayTabChar            = flag.Bool("t", false, "Display non-printing characters (see the -v option), and display tab characters as `^I'.")
	disableOutputBuffer       = flag.Bool("u", false, "Disable output buffering.")
)

func main() {
	flag.Parse()
	if len(os.Args) == 1 {
		io.Copy(os.Stdout, os.Stdin)
		return
	}

	for _, fn := range os.Args[1:] {
		fh, err := choose(fn)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		if fh == nil {
			continue
		}

		defer fh.Close()
		if *numberNonBlankOutputLines {
			printNumberedNonBlankLines(fh)
		} else if *displayDollarAtEnd {
			printDollarLines(fh)
		} else if *displayNonPrintingChars {
			printNonPrintingChars(fh)
		} else if *numberOutputLines {
			printNumberedLines(fh)
		} else if *singleSpacedOutput {
			printSingleSpacedLines(fh)
		} else if *displayTabChar {
			printNonPrintingChars(fh)
		} else {
			io.Copy(os.Stdout, fh)
		}
	}

}

func choose(name string) (io.ReadCloser, error) {
	splitArg := strings.Split(name, "-")
	if len(splitArg) == 2 && flag.Lookup(splitArg[1]) != nil {
		return nil, nil
	}
	if name == "-" {
		return os.Stdin, nil
	} else {
		return os.Open(name)
	}
}

func printNumberedNonBlankLines(rc io.ReadCloser) {
	scanner := bufio.NewScanner(rc)
	counter := 1
	for scanner.Scan() {
		var output string
		if scanner.Text() != "" {
			output = fmt.Sprintf("\t%d %s\n", counter, scanner.Text())
			counter++
		} else {
			output = "\n"
		}
		io.Copy(os.Stdout, strings.NewReader(output))
	}
}

func printNumberedLines(rc io.ReadCloser) {
	scanner := bufio.NewScanner(rc)
	counter := 1
	for scanner.Scan() {
		io.Copy(os.Stdout, strings.NewReader(fmt.Sprintf("\t%d %s\n", counter, scanner.Text())))
		counter++
	}
}

func printDollarLines(rc io.ReadCloser) {
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		var output string
		for _, value := range scanner.Bytes() {
			output += getNonPrintingChar(value)
		}
		io.Copy(os.Stdout, strings.NewReader(fmt.Sprintf("%s$\n", output)))
	}
}

func printNonPrintingChars(rc io.ReadCloser) {
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		var output string
		for _, value := range scanner.Bytes() {
			output += getNonPrintingChar(value)
		}
		io.Copy(os.Stdout, strings.NewReader(fmt.Sprintf("%s\n", output)))
	}
}

func printSingleSpacedLines(rc io.ReadCloser) {
	scanner := bufio.NewScanner(rc)
	counter := 0
	flag := false
	for scanner.Scan() {
		if scanner.Text() != "" {
			if flag && counter >= 1 {
				flag = false
				counter = 0
			}
			io.Copy(os.Stdout, strings.NewReader(fmt.Sprintf("%s\n", scanner.Text())))
		} else {
			if !flag {
				io.Copy(os.Stdout, strings.NewReader("\n"))
			}
			counter++
			flag = true
		}
	}
}

/*
	Non-printing character:
		Control characters print as `^X' for control-X;
		the delete character (octal 0177) prints as `^?'.
		Non-ASCII characters (with the high bit set) are printed as `M-' (for meta) followed by the character for the low 7 bits.
*/
func getNonPrintingChar(b byte) string {
	if b >= 32 {
		if b < 127 {
			return string(b)
		} else if b == 127 {
			return "^?"
		} else {
			str := "M-"
			if b >= 128+32 {
				if b < 128+127 {
					str += string(b - 128)
				} else {
					str += "^?"
				}
			} else {
				str += "^" + string(b-128+64)
			}
			return str
		}
	} else if b == '\t' && *displayTabChar {
		return "^I"
	} else if b == '\t' {
		return "\t"
	} else {
		return "^" + string(b+64)
	}
	return ""
}
