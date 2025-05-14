package fetch

import (
	"gh-api/actions/ask"
	"gh-api/actions/command"
	"gh-api/utilities"
	"log"
	"strings"
)

var chooseBranchLabel = "Choose a branch to change on: "

func Run() {
	// prune remote branches locally
	err := command.Run("git", "remote", "prune", "origin")

	if err != nil {
		log.Fatalf("Prune error: %v", err)
	}

	if err := utilities.FetchRemote(); err != nil {
		log.Fatalf("Fetch error: %v", err)
	}

	remoteBranches, err := utilities.GetRemoteBranches()

	if err != nil {
		log.Fatal("Fetch remote branches:", err)
	} else if len(remoteBranches) == 0 {
		log.Fatal("No branches on origin")
	}

	// Choose a branch
	remoteBranch, err := ask.One(&remoteBranches, &chooseBranchLabel)

	if err != nil {
		log.Fatal(chooseBranchLabel, err)
	}

	localBranch := strings.TrimPrefix(remoteBranch, "origin/")

	if utilities.HasLocalBranch(localBranch) {
		// TODO
		// If a branch exists already, ask
		// "Try to refetch? (Remove && Fetch) || Rebase"

		// Chosen branch exists already localy
		err := command.Run("git", "checkout", localBranch)

		if err != nil {
			log.Fatal("Checkout error: ", err)
		}
	} else {
		// There is not such a branch, fetch and checkout
		err := command.Run("git", "fetch", "origin", localBranch+":"+localBranch)

		if err != nil {
			log.Fatal("Fetch branch error: ", err)
		}

		err1 := command.Run("git", "checkout", localBranch)

		if err1 != nil {
			log.Fatal("Checkout error: ", err1)
		}
	}
}
