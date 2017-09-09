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
	Cat(os.Args, os.Stdin)
}

func Cat(args []string, r io.Reader) {
	if len(args) == 1 {
		io.Copy(os.Stdout, r)
		return
	}

	for _, fn := range args[1:] {
		fh, err := Choose(fn)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		if fh == nil {
			continue
		}

		defer fh.Close()
		counter := 1
		lineCounter := 0
		flag := false
		scanner := bufio.NewScanner(fh)
		for scanner.Scan() {
			output := scanner.Text()
			if *numberOutputLines && !*numberNonBlankOutputLines {
				output = fmt.Sprintf("    %d %s", counter, output)
				counter++
			}
			if *numberNonBlankOutputLines {
				if output != "" {
					output = fmt.Sprintf("    %d %s", counter, output)
					counter++
				}
			}
			if *displayDollarAtEnd || *displayNonPrintingChars || *displayTabChar {
				var nonPrintingStr string
				for _, value := range []byte(output) {
					nonPrintingStr += GetNonPrintingStr(value)
				}
				output = nonPrintingStr
			}
			if *singleSpacedOutput {
				if output != "" {
					if flag && lineCounter >= 1 {
						flag = false
						lineCounter = 0
					}
				} else {
					var str string
					if *displayDollarAtEnd {
						str = fmt.Sprintf("%s$", str)
					}
					if !flag {
						io.Copy(os.Stdout, strings.NewReader(str+"\n"))
					}
					lineCounter++
					flag = true
					continue
				}
			}
			if *displayDollarAtEnd {
				output = fmt.Sprintf("%s$", output)
			}
			output = fmt.Sprintf("%s\n", output)
			io.Copy(os.Stdout, strings.NewReader(output))
		}
	}
}

func Choose(name string) (io.ReadCloser, error) {
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

/*
	Non-printing character:
		Control characters print as `^X' for control-X;
		the delete character (octal 0177) prints as `^?'.
		Non-ASCII characters (with the high bit set) are printed as `M-' (for meta) followed by the character for the low 7 bits.

		160 = 128+32
		255 = 128+127
*/
func GetNonPrintingStr(b byte) string {
	switch {
	case b >= 32 && b < 127:
		return string(b)
	case b == 127:
		return "^?"
	case b > 127 && b < 160:
		return "M-^" + string(b-128+64)
	case b >= 160 && b < 255:
		return "M-" + string(b-128)
	case b >= 255:
		return "M-^?"
	case b == '\t' && *displayTabChar:
		return "^I"
	case b == '\t':
		return "\t"
	default:
		return "^" + string(b+64)
	}
}
