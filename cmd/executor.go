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
	var targetBranchNames []string

	err := exec.Command("git", "-v").Run()
	if err != nil {
		log.Fatal("Command not found: git")
	}

	localBranchNamesWithNewLine, err := execCommand("git", "for-each-ref", "refs/heads/", "--format=%(refname:short)")
	if err != nil {
		log.Fatal(err)
	}

	localBranchNames := strings.Split(localBranchNamesWithNewLine, "\n")
	fmt.Println(localBranchNames)

	if !slices.Contains(localBranchNames, baseBranchName) {
		log.Fatalf("Base branch not found: %s", baseBranchName)
	}

	for _, branchName := range localBranchNames {
		if branchName == baseBranchName {
			continue
		}

		fmt.Println(branchName)

		ancestorCommitObjHash, err := execCommand("git", "merge-base", baseBranchName, branchName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ancestorCommitObjHash)

		rootTreeObjHash, err := execCommand("git", "rev-parse", fmt.Sprintf("%s^{tree}", branchName))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(rootTreeObjHash)

		tmpCommitObjHash, err := execCommand("git", "commit-tree", rootTreeObjHash, "-p", ancestorCommitObjHash, "-m", "Temporary commit")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(tmpCommitObjHash)

		cherryResult, err := execCommand("git", "cherry", baseBranchName, tmpCommitObjHash)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(cherryResult)

		if strings.HasPrefix(cherryResult, "- ") {
			targetBranchNames = append(targetBranchNames, branchName)
		}
	}

	fmt.Println(targetBranchNames)
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
