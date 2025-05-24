package branch

import (
	"fmt"
	"gh-api/actions/ask"
	"gh-api/actions/command"
	"gh-api/git"
	"strings"
)

var chooseBranchLabel = "Choose a branch to change on: "
var createLabel = "CREATE NEW BRANCH"
var backLabel = "< back"

func Run() {
	localBranches, currentBranch, _ := git.GetLocalBranches()

	remoteBranches, _ := git.GetBranches(string(git.Remote))

	options := append([]string{createLabel}, localBranches...)
	options = append(options, remoteBranches...)
	options = append(options, backLabel)

	// Choose a branch
	answer, _ := ask.ChooseBranch(&options, &chooseBranchLabel, currentBranch)

	if answer == createLabel {
		var branch string
		fmt.Print("Type a name for a new branch: ")
		fmt.Scan(&branch)
		git.CreateBranchAndCheckout(branch)
		return
	}

	if answer == backLabel {
		return
	}

	branch := strings.TrimPrefix(answer, "origin/")

	// There is not such a branch, fetch and checkout
	if !git.HasLocalBranch(branch) {
		command.Run("git", "fetch", "origin", branch+":"+branch)
		fmt.Println("\n---------------")
		command.Run("git", "checkout", branch)
		fmt.Println("---------------")
	} else {
		fmt.Println("\n---------------")
		command.Run("git", "checkout", branch)
		command.Run("git", "pull", "--rebase", "origin", branch)
		fmt.Println("\n---------------")
	}
}
