package main

import (
	"fmt"
	"gh-api/actions/ask"
	"gh-api/actions/branch"
	"gh-api/actions/command"
	pullrequest "gh-api/actions/pull-request"
	"gh-api/actions/remove"
	"gh-api/git"
	"log"
)

var version = "1.0.22"
var BRANCH = "branch"
var PULL_REQUEST = "pull request"
var REMOVE_BRANCH = "remove branch"
var FETCH_REMOTE = "fetch remote"

var actions = []string{BRANCH, PULL_REQUEST, REMOVE_BRANCH, FETCH_REMOTE}

func main() {
	fmt.Printf("v.%v\n\n", version)

	err := git.Is()

	if err != nil {
		log.Fatal("Git is not inited", err)
	}

	for {
		action, err := ask.One(&actions, "Choose action:")

		if err != nil {
			log.Fatal("Choose action:", err)
		}

		if action == BRANCH {
			branch.Run()
		} else if action == PULL_REQUEST {
			pullrequest.Run()
		} else if action == REMOVE_BRANCH {
			remove.Run()
		} else if action == FETCH_REMOTE {
			command.Run("git", "remote", "prune", "origin")
			git.FetchRemote()
		}

		fmt.Println()
		action = ""
	}
}
