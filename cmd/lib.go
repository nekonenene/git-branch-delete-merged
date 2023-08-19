package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Wrapper method of exec.Command
// error contains stderr message
func ExecCommand(name string, arg ...string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("%v\n\n%v", err, stderr.String())
		return stdout.String(), fmt.Errorf("%s command failed:\n\n%s", name, stderr.String())
	}

	trimmedString := strings.TrimRight(stdout.String(), "\n") // Remove newlines from end of string

	return trimmedString, nil
}
