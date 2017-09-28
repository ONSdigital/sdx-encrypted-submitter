package main

import (
	"testing"
)

var errorMessage string

func TestNotSupplyingSourceFileName(t *testing.T) {
	testArgs = []string{" "}
	onExit = onExitHandler

	defer func() {
		r := recover()
		if r != nil {
			var expected = "could not read message body: no file name supplied"
			if errorMessage != expected {
				t.Errorf("Expected = %s : actual = %s", expected, errorMessage)
			}
		}
	}()

	main()

}

// Panics but does not exit , so bypasses all the code in main after panic and is caught in deferred panic handler
func onExitHandler(msg string) {
	errorMessage = msg
	panic(msg)
}
