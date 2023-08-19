package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"slices"
	"strings"
)

func Exec() {
	ParseParameters()
	baseBranchName := params.BaseBranchName

	err := exec.Command("git", "-v").Run()
	if err != nil {
		log.Fatal("Command not found: git")
	}

	output, err := execCommand("git", "for-each-ref", "refs/heads/", "--format=%(refname:short)")
	if err != nil {
		log.Fatal(err)
	}

	branchNames := strings.Split(output, "\n")
	fmt.Println(branchNames)

	if !slices.Contains(branchNames, baseBranchName) {
		log.Fatalf("Base branch not found: %s", baseBranchName)
	}

	for _, branchName := range branchNames {
		if branchName == baseBranchName {
			continue
		}

		fmt.Println(branchName)
	}
}

func execCommand(name string, arg ...string) (string, error) {
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
