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
	- разбить файл на маленькие файлы
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
		branches, branchesError := utilities.GetLocalBranches()

		if branchesError != nil {
			log.Fatal("Error for local branches: ", branchesError)
		} else if len(branches) == 0 {
			log.Fatalf("There are no local branches to remove")
		}

		selected := []string{}

		branchesPrompt := &survey.MultiSelect{
			Message: chooseBranchesToRemove,
			Options: branches,
		}

		branchError := survey.AskOne(branchesPrompt, &selected)

		if branchError != nil {
			log.Fatal(chooseBranchesToRemove, branchError)
		}

		fmt.Println(selected)

		utilities.RemoveBranches(selected)
	}
}
