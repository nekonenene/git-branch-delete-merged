package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

type parameters struct {
	BaseBranchName string
	ShowVersion    bool
}

var params parameters
var Version = "" // Overwrite when building

func ParseParameters() {
	flag.StringVar(&params.BaseBranchName, "base-branch", "", "[Required] Base branch name (e.g. main, develop)")
	flag.BoolVar(&params.ShowVersion, "version", false, "[Opiton] Show version")
	flag.BoolVar(&params.ShowVersion, "v", false, "[Opiton] Shorthand of -version")
	flag.Parse()

	if params.ShowVersion {
		fmt.Println(getVersionString())
		os.Exit(0)
	}
	if params.BaseBranchName == "" {
		log.Fatalln("-base-branch is required")
	}
}

func getVersionString() string {
	// For downloading a binary from GitHub Releases
	if Version != "" {
		return Version
	}

	// For `go install`
	i, ok := debug.ReadBuildInfo()
	if ok {
		return i.Main.Version
	}

	return "Development version"
}
