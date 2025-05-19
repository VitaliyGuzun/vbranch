package branch

import (
	"fmt"
	"gh-api/actions/ask"
	"gh-api/actions/command"
	"gh-api/utilities"
	"strings"
)

var chooseBranchLabel = "Choose a branch to change on: "
var backLabel = "< back"

func Run() {
	remoteBranches, _ := utilities.GetRemoteBranches()
	remoteBranches = append(remoteBranches, backLabel)

	// Choose a branch
	remoteBranch, _ := ask.One(&remoteBranches, &chooseBranchLabel)

	if remoteBranch == backLabel {
		return
	}

	localBranch := strings.TrimPrefix(remoteBranch, "origin/")

	// There is not such a branch, fetch and checkout
	if !utilities.HasLocalBranch(localBranch) {
		command.Run("git", "fetch", "origin", localBranch+":"+localBranch)
		fmt.Println("\n---------------")
		command.Run("git", "checkout", localBranch)
		fmt.Println("---------------")
	} else {
		fmt.Println("\n---------------")
		command.Run("git", "checkout", localBranch)
		command.Run("git", "pull", "--rebase", "origin", localBranch)
		fmt.Println("\n---------------")
	}
}
