package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func Exec() {
	ParseParameters()

	baseBranchName := params.BaseBranchName
	var targetBranchNames []string

	err := exec.Command("git", "version").Run()
	if err != nil {
		log.Fatal("Command not found: git")
	}

	localBranchNamesWithNewLine, err := ExecCommand("git", "for-each-ref", "refs/heads/", "--format", "%(refname:short)")
	if err != nil {
		log.Fatal(err)
	}

	localBranchNames := strings.Split(localBranchNamesWithNewLine, "\n")

	if !slices.Contains(localBranchNames, baseBranchName) {
		log.Fatalf("Base branch not found: %s", baseBranchName)
	}

	fmt.Println("Searching merged branches...")

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

	currentBranchName, err := ExecCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		log.Fatal(err)
	}

	if len(targetBranchNames) == 0 {
		fmt.Printf("\033[33mThere is no branch which merged into '%s'\033[0m\n", baseBranchName)
		return
	} else {
		fmt.Printf("Found \033[32m%d\033[0m merged branches: %s\n", len(targetBranchNames), targetBranchNames)
	}

	for _, targetBranchName := range targetBranchNames {
		if targetBranchName == baseBranchName {
			continue
		}

		if targetBranchName == currentBranchName {
			fmt.Printf("\033[33mSkipped '%s' branch because it is current branch\033[0m\n", targetBranchName)
			continue
		}

		deleteBranchPrompt(targetBranchName, params.AllYesFlag)
	}
}

// Confirm whether to delete the branch
func deleteBranchPrompt(targetBranchName string, yesFlag bool) {
	loopEndFlag := false

	for !loopEndFlag {
		var response string

		if yesFlag {
			response = "yes"
		} else {
			fmt.Printf("\nAre you sure to delete \033[33m'%s'\033[0m branch? [y|n|l|d|q|help]: ", targetBranchName)
			fmt.Scanln(&response)
		}

		switch response {
		case "y", "yes":
			latestCommitId, err := ExecCommand("git", "rev-parse", targetBranchName)
			if err != nil {
				log.Fatal(err)
			}

			_, err = ExecCommand("git", "branch", "-D", targetBranchName)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("\033[32mDeleted '%s' branch\033[0m\n", targetBranchName)
			fmt.Printf("You can recreate this branch with `git branch %s %s`\n", targetBranchName, latestCommitId)
			loopEndFlag = true
		case "n", "no":
			fmt.Println("Skipped")
			loopEndFlag = true
		case "l", "log":
			err := DelegateCommand("git", "log", targetBranchName, "-100") // Show only 100 logs to avoid broken pipe error
			if err != nil {
				log.Fatal(err)
			}
		case "d", "diff":
			err := DelegateCommand("git", "show", targetBranchName, "-v")
			if err != nil {
				log.Fatal(err)
			}
		case "q", "quit":
			fmt.Println("Suspends processing")
			os.Exit(1)
		case "h", "help":
			fmt.Println("y: Yes, delete the branch")
			fmt.Println("n: No, skip deleting")
			fmt.Println("l: Show git logs of the branch")
			fmt.Println("d: Show the latest commit of the branch and its diff")
			fmt.Println("q: Quit immediately")
			fmt.Println("h: Display this help")
		}
	}
}
