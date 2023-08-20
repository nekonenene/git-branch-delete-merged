package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Remove matched string from slice
func RemoveFromSlice(slice []string, targetString string) []string {
	var result []string

	for _, item := range slice {
		if item != targetString {
			result = append(result, item)
		}
	}

	return result
}

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
		return stdout.String(), fmt.Errorf("%s command failed:\n\n%s", name, stderr.String())
	}

	trimmedString := strings.TrimRight(stdout.String(), "\n") // Remove newlines from end of string

	return trimmedString, nil
}

// Display the command as usual that requires user operation, such as `less` or `top`
func DelegateCommand(name string, arg ...string) error {
	var stderr bytes.Buffer

	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s command failed:\n\n%s", name, stderr.String())
	}

	return nil
}

// Check if branch is squash and merged
func IsSquashedBranch(baseBranchName string, branchName string) (bool, error) {
	ancestorCommitObjHash, err := ExecCommand("git", "merge-base", baseBranchName, branchName)
	if err != nil {
		return false, err
	}

	rootTreeObjHash, err := ExecCommand("git", "rev-parse", fmt.Sprintf("%s^{tree}", branchName))
	if err != nil {
		return false, err
	}

	tmpCommitObjHash, err := ExecCommand("git", "commit-tree", rootTreeObjHash, "-p", ancestorCommitObjHash, "-m", "Temporary commit")
	if err != nil {
		return false, err
	}

	cherryResult, err := ExecCommand("git", "cherry", baseBranchName, tmpCommitObjHash)
	if err != nil {
		return false, err
	}

	if strings.HasPrefix(cherryResult, "- ") {
		return true, nil
	} else {
		return false, nil
	}
}
