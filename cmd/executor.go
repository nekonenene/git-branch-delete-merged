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

	localBranchNamesWithNewLine, err := ExecCommand("git", "for-each-ref", "refs/heads/", "--format=%(refname:short)")
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

		ancestorCommitObjHash, err := ExecCommand("git", "merge-base", baseBranchName, branchName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ancestorCommitObjHash)

		rootTreeObjHash, err := ExecCommand("git", "rev-parse", fmt.Sprintf("%s^{tree}", branchName))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(rootTreeObjHash)

		tmpCommitObjHash, err := ExecCommand("git", "commit-tree", rootTreeObjHash, "-p", ancestorCommitObjHash, "-m", "Temporary commit")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(tmpCommitObjHash)

		cherryResult, err := ExecCommand("git", "cherry", baseBranchName, tmpCommitObjHash)
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
