package main

import (
	"bufio"
	"bytes"
	"fmt"
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
	- собрать свой пакет и опубликовать в него этот скрипт
	- при выполнении скрипта, проверять что есть обновления и предлагать обновить
	- добавить логи в каждую функцию, чтобы юзер видел что происходит
	- локализовать
	- запакопать в докер? не уверен что это нужно

	---------------
	$ branch

	$ ✅ fetch
	$ remove

	$ origin/main
	$ ✅ origin/test
	$ origin/branch

	## git fetch origin origin/test:origin/test && git checkout origin/test

	---------------
	$ branch

	$ fetch
	$ ✅ remove

	$ origin/main
	$ ✅ origin/test
	$ ✅ origin/branch

	## git branch -D origin/test origin/branch
*/

var chooseActionLabel = "Choose action: "
var chooseBranchLabel = "Choose a branch to change on: "
var chooseBranchesToRemove = "Choose branches to remove: "

// Проверка, что мы в git-репозитории
func isGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()

	if err != nil || strings.TrimSpace(string(output)) != "true" {
		return fmt.Errorf("not a git repository")
	}

	return nil
}

func fetchRemote() error {
	command := exec.Command("git", "fetch")
	_, error := command.Output()

	if error != nil {
		return error
	}

	return nil
}

// Получение всех удалённых веток
func getRemoteBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "-r")
	output, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	branches := []string{}
	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		branch := strings.TrimSpace(scanner.Text())

		if strings.Contains(branch, "->") {
			// Пропускаем алиасы типа origin/HEAD -> origin/main
			continue
		}

		branches = append(branches, branch)
	}

	return branches, nil
}

// Получение локальных веток
func getLocalBranches() ([]string, error) {
	command := exec.Command("git", "branch")
	output, error := command.Output()

	if error != nil {
		return nil, error
	}

	branches := []string{}
	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		branch := strings.TrimSpace(scanner.Text())

		if strings.Contains(branch, "->") {
			// Пропускаем алиасы типа origin/HEAD -> origin/main
			continue
		}

		branches = append(branches, branch)
	}

	return branches, nil
}

// Удаление бранчей
func removeBranches(branches []string) error {
	args := append([]string{"branch", "-D"}, branches...)
	fmt.Println("Removing branches:", branches)

	command := exec.Command("git", args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	removeError := command.Run()
	if removeError != nil {
		return removeError
	}

	return nil
}

// Проверка, есть ли локальная ветка
func hasLocalBranch(branch string) bool {
	cmd := exec.Command("git", "branch", "--list", branch)
	output, _ := cmd.Output()

	return strings.TrimSpace(string(output)) != ""
}

func main() {
	if error := isGitRepo(); error != nil {
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
		if fetchAllError := fetchRemote(); fetchAllError != nil {
			log.Fatalf("Fetch error: %v", fetchAllError)
		}

		branches, branchesError := getRemoteBranches()

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

		if hasLocalBranch(localBranch) {
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
		branches, branchesError := getLocalBranches()

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

		removeBranches(selected)
	}
}
