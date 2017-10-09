package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	lines       = flag.String("n", "", "display n first lines of a file")
	bytes       = flag.String("c", "", "display n first bytes of a file")
	displayName bool
)

func displayLineResult(names []string, l int) {
	if len(names) == 0 {
		count := 1
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			os.Stdout.Write([]byte(s.Text() + "\n"))
			if count == l {
				return
			}
			count++
		}
		return
	}
	for _, name := range names {
		if displayName {
			io.Copy(os.Stdout, strings.NewReader(fmt.Sprintf("==> %s <==\n", name)))
		}
		file, err := os.Open(name)
		if err != nil {
			io.Copy(os.Stdout, strings.NewReader(err.Error()+"\n"))
		} else {
			scanner := bufio.NewScanner(file)
			var count int
			for scanner.Scan() {
				line := scanner.Text()
				if count < l {
					io.Copy(os.Stdout, strings.NewReader(line+"\n"))
				}
				count++
			}
		}
	}
}

func displayByteResult(names []string, b int) {
	if len(names) == 0 {
		var final string
		for {
			reader := bufio.NewReader(os.Stdin)
			text, err := reader.ReadString('\n')
			if err != nil {
				io.Copy(os.Stdout, strings.NewReader(err.Error()+"\n"))
				break
			}
			if len(final) == 0 {
				final = fmt.Sprintf("%s", text)
			} else {
				final = fmt.Sprintf("%s%s", final, text)
			}

			if len(final) >= b {
				io.Copy(os.Stdout, strings.NewReader(final[:b]))
				break
			}
		}
		return
	}

	for i, name := range names {
		if displayName {
			io.Copy(os.Stdout, strings.NewReader(fmt.Sprintf("==> %s <==\n", name)))
		}
		file, err := os.Open(name)
		if err != nil {
			io.Copy(os.Stdout, strings.NewReader(err.Error()+"\n"))
		} else {
			byteArray := make([]byte, b)
			_, err := file.Read(byteArray)
			if err != nil {
				io.Copy(os.Stdout, strings.NewReader(err.Error()))
			}
			if i == len(names)-1 {
				os.Stdout.Write(byteArray)
			} else {
				os.Stdout.Write([]byte(string(byteArray) + "\n"))
			}

		}
	}
}

func main() {
	flag.Parse()
	var options int

	options = 1
	if *lines != "" && *bytes != "" {
		io.Copy(os.Stdout, strings.NewReader("head: can't combine line and byte counts\n"))
		return
	}
	if *lines != "" {
		lineCount, err := strconv.Atoi(*lines)
		if err != nil || lineCount <= 0 {
			io.Copy(os.Stdout, strings.NewReader(fmt.Sprintf("head: illegal line count -- %s\n", *lines)))
			return
		}
		options += 2
		if len(os.Args[options:]) > 1 {
			displayName = true
		}
		displayLineResult(os.Args[options:], lineCount)
	}
	if *bytes != "" {
		byteCount, err := strconv.Atoi(*bytes)
		if err != nil || byteCount <= 0 {
			io.Copy(os.Stdout, strings.NewReader(fmt.Sprintf("head: illegal byte count -- %s\n", *bytes)))
			return
		}
		options += 2
		if len(os.Args[options:]) > 1 {
			displayName = true
		}
		displayByteResult(os.Args[options:], byteCount)

	}
	if *lines == "" && *bytes == "" {
		if len(os.Args[options:]) > 1 {
			displayName = true
		}
		displayLineResult(os.Args[options:], 10)
	}
}
