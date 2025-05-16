package remove

import (
	"fmt"
	"gh-api/actions/ask"
	"gh-api/utilities"
	"log"
)

var chooseBranchesToRemove = "Choose branches to remove: "

func Run() {
	branches, currentBranch, err := utilities.GetLocalBranches()

	if err != nil {
		log.Fatal("Error for local branches: ", err)
	} else if len(branches) == 0 {
		log.Fatal("There are no local branches to remove")
	} else if len(branches) == 1 {
		fmt.Println("🔴 Error:")
		fmt.Printf("   You the only branch you have is \"%v\".\n", branches[0])
		fmt.Println("   You can't remove it because there is no other branch to switch to.")
		fmt.Printf("   If you still want to remove \"%v\", create another branch before to switch to.\n", branches[0])
		return
	}

	removeBranches, err := ask.Many(&branches, &chooseBranchesToRemove)

	if err != nil {
		log.Fatal(chooseBranchesToRemove, err)
	}

	// If user selected the current branch for removing, we have to checkout to another branch
	if utilities.ShouldChangeBranch(removeBranches, currentBranch) {
		var checkoutBranch string

		// Go through all branches searching for a branch that is not in the branches to remove
		for _, branch := range branches {
			shouldSkip := false

			for _, removeBranche := range removeBranches {
				if removeBranche == branch {
					shouldSkip = true
					break
				}
			}

			if !shouldSkip {
				checkoutBranch = branch
				break
			}
		}

		// If there is a branch to checkout: do checkout before removing
		if checkoutBranch != "" {
			utilities.Checkout(checkoutBranch)
		} else {
			fmt.Println("🔴 Error:")
			fmt.Printf("   Looks like you choose all branches to remove\n")
			fmt.Printf("   You can't remove all branches. There must be at least one branch left.\n")
			fmt.Printf("   Please, leave at least one branch not selected.\n")
		}
	}

	utilities.RemoveBranches(removeBranches)
}
