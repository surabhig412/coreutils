package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() error {
	return nil
}

func TestGetNonPrintingStr(t *testing.T) {
	if GetNonPrintingStr(65) != "A" {
		t.Errorf("Wrong string for byte between 32 and 127")
	}
	if GetNonPrintingStr(127) != "^?" {
		t.Errorf("Wrong string for byte 127")
	}
	if GetNonPrintingStr(130) != "M-^B" {
		t.Errorf("Wrong string for byte between 127 and 160")
	}
	if GetNonPrintingStr(170) != "M-*" {
		t.Errorf("Wrong string for byte between 160 and 255")
	}
	if GetNonPrintingStr(255) != "M-^?" {
		t.Errorf("Wrong string for byte greater than 255")
	}
	if GetNonPrintingStr('\t') != "\t" {
		t.Errorf("Wrong string for byte \t")
	}
	if GetNonPrintingStr(30) != "^^" {
		t.Errorf("Wrong string for byte less than 32")
	}
	*displayTabChar = true
	if GetNonPrintingStr('\t') != "^I" {
		t.Errorf("Wrong string for byte \t and ")
	}
	*displayTabChar = false
}

func TestChoose(t *testing.T) {
	rc, err := Choose("-b")
	if rc != nil || err != nil {
		t.Errorf("Error and input handler is nil for flag argument")
	}
	rc, err = Choose("-")
	if rc != os.Stdin || err != nil {
		t.Errorf("Error is nil and input handler is std input for flag argument")
	}
	rc, err = Choose("Makefile")
	if err != nil {
		t.Errorf("Error is nil and input handler is file for flag argument")
	}
}

func TestCat1(t *testing.T) {
	// when no input handler and flags are sent
	args := []string{"./cat"}
	var rc io.ReadCloser
	cb := &ClosingBuffer{bytes.NewBufferString("abcd\n")}
	rc = cb
	out := mockStdout(args, rc)
	if out != "abcd\n" {
		t.Errorf("Error when no input handler and flags are sent")
	}
}

func TestCat2(t *testing.T) {
	// when no flags are given
	args := []string{"./cat", "testdata/dummy_data"}
	out := mockStdout(args, os.Stdin)
	f, _ := os.OpenFile("testdata/dummy_data", os.O_RDONLY, 0755)
	b := make([]byte, 1024)
	f.Read(b)
	if !equal(out, b) {
		t.Errorf("Error when no flags are sent")
	}
}

func TestCat3(t *testing.T) {
	// error when invalid filename given
	if os.Getenv("EXIT") == "true" {
		args := []string{"./cat", "abcd"}
		Cat(args, os.Stdin)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCat3")
	cmd.Env = append(os.Environ(), "EXIT=true")
	output, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Error(err)
	}

	outBytes, _ := ioutil.ReadAll(output)
	if !equal("open abcd: no such file or directory\n", outBytes) {
		t.Error("error string when invalid filename doesn't match")
	}

	err := cmd.Wait()
	if e, ok := err.(*exec.ExitError); !ok || e.Success() {
		t.Error("error when invalid filename given doesn't exit")
	}
}

func TestCatFlagb(t *testing.T) {
	args := []string{"./cat", "-b", "testdata/dummy_data"}
	*numberNonBlankOutputLines = true
	defer func() {
		*numberNonBlankOutputLines = false
	}()
	out := mockStdout(args, os.Stdin)
	f, _ := os.OpenFile("testdata/dummy_data_flagb", os.O_RDONLY, 0755)
	b := make([]byte, 1024)
	f.Read(b)
	if !equal(out, b) {
		t.Errorf("Error when flag b is sent")
	}
}

func TestCatFlagev(t *testing.T) {
	args := []string{"./cat", "-e", "testdata/dummy_data"}
	*displayDollarAtEnd = true
	defer func() {
		*displayDollarAtEnd = false
	}()
	out := mockStdout(args, os.Stdin)
	f, _ := os.OpenFile("testdata/dummy_data_flage", os.O_RDONLY, 0755)
	b := make([]byte, 1024)
	f.Read(b)
	if !equal(out, b) {
		t.Errorf("Error when flag e is sent")
	}
}

func TestCatFlagn(t *testing.T) {
	args := []string{"./cat", "-n", "testdata/dummy_data"}
	*numberOutputLines = true
	defer func() {
		*numberOutputLines = false
	}()
	out := mockStdout(args, os.Stdin)
	f, _ := os.OpenFile("testdata/dummy_data_flagn", os.O_RDONLY, 0755)
	b := make([]byte, 1024)
	f.Read(b)
	if !equal(out, b) {
		t.Errorf("Error when flag n is sent")
	}
}

func TestCatFlags(t *testing.T) {
	args := []string{"./cat", "-s", "testdata/dummy_data"}
	*singleSpacedOutput = true
	defer func() {
		*singleSpacedOutput = false
	}()
	out := mockStdout(args, os.Stdin)
	f, _ := os.OpenFile("testdata/dummy_data_flags", os.O_RDONLY, 0755)
	b := make([]byte, 1024)
	f.Read(b)
	if !equal(out, b) {
		t.Errorf("Error when flag s is sent")
	}
}

func TestCatMultipleFlags(t *testing.T) {
	args := []string{"./cat", "-b", "-e", "-s", "-u", "-t", "testdata/dummy_data"}
	*numberNonBlankOutputLines = true
	*displayDollarAtEnd = true
	*singleSpacedOutput = true
	*displayTabChar = true
	*disableOutputBuffer = true
	defer func() {
		*numberNonBlankOutputLines = false
		*displayDollarAtEnd = false
		*singleSpacedOutput = false
		*displayTabChar = false
		*disableOutputBuffer = false
	}()
	out := mockStdout(args, os.Stdin)
	f, _ := os.OpenFile("testdata/dummy_data_multiple_flags", os.O_RDONLY, 0755)
	b := make([]byte, 1024)
	f.Read(b)
	if !equal(out, b) {
		t.Errorf("Error when multiple flags are sent")
	}
}

func equal(str string, b []byte) bool {
	for i, value := range b {
		if i < len(str) && value != str[i] {
			fmt.Println("error: ", string(value), " :", i)
			return false
		}
	}
	return true
}

func mockStdout(args []string, reader io.ReadCloser) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	Cat(args, reader)
	output := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		output <- buf.String()
	}()

	w.Close()
	os.Stdout = oldStdout
	return <-output
}
