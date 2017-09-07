package main

import (
	"os"
	"testing"
)

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

func TestCat(t *testing.T) {
	args := []string{"./cat"}
	Cat(args)
}
