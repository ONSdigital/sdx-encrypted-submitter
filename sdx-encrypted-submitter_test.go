package main

import (
	"os/exec"
	"strings"
	"testing"
)

// Note these are integration tests that assume that the sdx-encrypted-submitter has been installed

func TestNotSupplyingAnyArgument(t *testing.T) {

	cmd := exec.Command("sdx-encrypted-submitter", "-e","something.txt", "-s", "somethingelse.txt")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("No error when one was expected")
	}
	var expected = "could not read message body  -  no file name supplied\n"
	var actual = string(output)
	if actual != expected {
		t.Error("expected ", expected, "got ", actual)
	}
}

func TestSupplyingUnknownArgument(t *testing.T) {
	cmd := exec.Command("sdx-encrypted-submitter", "-Y")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("No error when one was expected")
	}
	var expected = "flag provided but not defined"

	var actual = string(output)
	if !strings.Contains(actual, expected) {
		t.Error("'", expected, "' not in the output ")
	}
}

func TestUnableToReadSourceFile(t *testing.T) {
	cmd := exec.Command("sdx-encrypted-submitter", "-f", "AFileThatClearlyDoesNotExist","-e","something.txt", "-s", "somethingelse.txt")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("No error when one was expected")
	}
	var expected = "could not read message body"

	var actual = string(output)
	if !strings.Contains(actual, expected) {
		t.Error("'", expected, "' not in the output ")
	}
}
