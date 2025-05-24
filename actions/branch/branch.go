package branch

import (
	"fmt"
	"gh-api/actions/ask"
	"gh-api/actions/command"
	"gh-api/git"
	"strings"
)

var chooseBranchLabel = "Choose a branch to change on: "
var backLabel = "< back"

func Run() {
	localBranches, currentBranch, _ := git.GetLocalBranches()

	remoteBranches, _ := git.GetBranches(string(git.Remote))
	remoteBranches = append(remoteBranches, backLabel)

	allBranches := append(localBranches, remoteBranches...)

	// Choose a branch
	remoteBranch, _ := ask.ChooseBranch(&allBranches, &chooseBranchLabel, currentBranch)

	if remoteBranch == backLabel {
		return
	}

	localBranch := strings.TrimPrefix(remoteBranch, "origin/")

	// There is not such a branch, fetch and checkout
	if !git.HasLocalBranch(localBranch) {
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
