package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func Exec() {
	ParseParameters()

	err := exec.Command("git", "-v").Run()
	if err != nil {
		log.Fatal("command not found: git")
	}

	output, err := exec.Command("git", "for-each-ref", "refs/heads/", "--format=%(refname:short)").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	outputString := strings.TrimSpace(string(output))
	branchNames := strings.Split(outputString, "\n")

	for _, branchName := range branchNames {
		fmt.Println(branchName)
	}
}
