package cmd

import (
	"fmt"
	"log"
	"os/exec"
)

func Exec() {
	ParseParameters()

	output, err := exec.Command("git", "-v").CombinedOutput()
	if err != nil {
		log.Fatal("command not found: git")
	}

	fmt.Println("Hello!")
	fmt.Printf("%s", output)
}
