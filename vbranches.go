package main

import (
	"fmt"
	"gh-api/utilities"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

/*
	TODO:
	- добавить тесты к каждой функции
	- при выполнении скрипта, проверять что есть обновления и предлагать обновить
	- добавить логи в каждую функцию, чтобы юзер видел что происходит
*/

var chooseActionLabel = "Choose action: "
var chooseBranchLabel = "Choose a branch to change on: "
var chooseBranchesToRemove = "Choose branches to remove: "

func main() {
	if error := utilities.IsGitRepo(); error != nil {
		log.Fatal("Git is not inited", error)
	}

	// Choose action: fetch | remove
	var action string
	actions := []string{"fetch", "remove"}

	actionsPrompt := &survey.Select{
		Message: chooseActionLabel,
		Options: actions,
	}

	actionError := survey.AskOne(actionsPrompt, &action)

	if actionError != nil {
		log.Fatalf(chooseActionLabel, actionError)
	}

	if action == "fetch" {
		if fetchAllError := utilities.FetchRemote(); fetchAllError != nil {
			log.Fatalf("Fetch error: %v", fetchAllError)
		}

		branches, branchesError := utilities.GetRemoteBranches()

		if branchesError != nil {
			log.Fatal("Fetch remote branches:", branchesError)
		} else if len(branches) == 0 {
			log.Fatal("No branches on origin")
		}

		// Choose a branch
		var selected string

		prompt := &survey.Select{
			Message: chooseBranchLabel,
			Options: branches,
		}

		branchError := survey.AskOne(prompt, &selected)

		if branchError != nil {
			log.Fatal(chooseBranchLabel, branchError)
		}

		localBranch := strings.TrimPrefix(selected, "origin/")

		if utilities.HasLocalBranch(localBranch) {
			// TODO
			// If a branch exists already, ask
			// "Try to refetch? (Remove && Fetch) || Rebase"

			// Chosen branch exists already localy
			checkoutCmd := exec.Command("git", "checkout", localBranch)
			checkoutCmd.Stdout = os.Stdout
			checkoutCmd.Stderr = os.Stderr
			checkoutError := checkoutCmd.Run()

			if checkoutError != nil {
				log.Fatal("Checkout error: ", checkoutError)
			}
		} else {
			// There is not such a branch, fetch and checkout
			fetchCmd := exec.Command("git", "fetch", "origin", localBranch+":"+localBranch)
			fetchCmd.Stdout = os.Stdout
			fetchCmd.Stderr = os.Stderr
			fetchError := fetchCmd.Run()

			if fetchError != nil {
				log.Fatal("Fetch branch error: ", fetchError)
			}

			checkoutCmd := exec.Command("git", "checkout", localBranch)
			checkoutCmd.Stdout = os.Stdout
			checkoutCmd.Stderr = os.Stderr
			checkoutError := checkoutCmd.Run()

			if checkoutError != nil {
				log.Fatal("Checkout error: ", checkoutError)
			}
		}
	} else if action == "remove" {
		branches, currentBranch, branchesError := utilities.GetLocalBranches()

		if branchesError != nil {
			log.Fatal("Error for local branches: ", branchesError)
		} else if len(branches) == 0 {
			log.Fatal("There are no local branches to remove")
		} else if len(branches) == 1 {
			fmt.Println("🔴 Error:")
			fmt.Printf("   You can't remove branch \"%v\" because there is no other branch to switch to.\n", branches[0])
			fmt.Printf("   If you still want to remove \"%v\", create another branch before to switch to.\n", branches[0])
			return
		}

		removeBranches := []string{}

		branchesPrompt := &survey.MultiSelect{
			Message: chooseBranchesToRemove,
			Options: branches,
		}

		branchError := survey.AskOne(branchesPrompt, &removeBranches)

		if branchError != nil {
			log.Fatal(chooseBranchesToRemove, branchError)
		}

		// If user selected the current branch for removing, we have to checkout to another branch
		if utilities.ShouldChangeBranch(branches, currentBranch) {
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
}
