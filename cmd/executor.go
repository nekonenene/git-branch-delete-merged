package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"slices"
	"strings"
)

func Exec() {
	ParseParameters()

	baseBranchName := params.BaseBranchName
	var targetBranchNames []string

	err := exec.Command("git", "-v").Run()
	if err != nil {
		log.Fatal("Command not found: git")
	}

	localBranchNamesWithNewLine, err := ExecCommand("git", "for-each-ref", "refs/heads/", "--format", "%(refname:short)")
	if err != nil {
		log.Fatal(err)
	}

	localBranchNames := strings.Split(localBranchNamesWithNewLine, "\n")
	fmt.Println(localBranchNames)

	if !slices.Contains(localBranchNames, baseBranchName) {
		log.Fatalf("Base branch not found: %s", baseBranchName)
	}

	mergedBranchNamesWithNewLine, err := ExecCommand("git", "branch", "--merged", baseBranchName, "--format", "%(refname:short)")
	if err != nil {
		log.Fatal(err)
	}

	mergedBranchNames := strings.Split(mergedBranchNamesWithNewLine, "\n")
	mergedBranchNames = RemoveFromSlice(mergedBranchNames, baseBranchName)

	targetBranchNames = append(targetBranchNames, mergedBranchNames...)

	for _, localBranchName := range localBranchNames {
		if localBranchName == baseBranchName {
			continue
		}

		isSquashed, err := IsSquashedBranch(baseBranchName, localBranchName)
		if err != nil {
			log.Fatal(err)
		}

		if isSquashed {
			targetBranchNames = append(targetBranchNames, localBranchName)
		}
	}

	fmt.Println(targetBranchNames)

	currentBranchName, err := ExecCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		log.Fatal(err)
	}

	for _, targetBranchName := range targetBranchNames {
		if targetBranchName == baseBranchName {
			continue
		}

		if targetBranchName == currentBranchName {
			fmt.Printf("\033[33mSkipped '%s' branch because it is current branch\033[0m\n", targetBranchName)
			continue
		}

		fmt.Printf("Target branch: %s\n", targetBranchName)
	}
}
