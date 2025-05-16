package utilities

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func IsGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()

	if err != nil || strings.TrimSpace(string(output)) != "true" {
		return fmt.Errorf("not a git repository")
	}

	return nil
}

func FetchRemote() error {
	command := exec.Command("git", "fetch")
	_, error := command.Output()

	if error != nil {
		return error
	}

	return nil
}

// Получение всех удалённых веток
func GetRemoteBranches() ([]string, error) {
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

func GetLocalBranches() ([]string, string, error) {
	command := exec.Command("git", "branch")
	output, error := command.Output()
	var currentBranch string

	if error != nil {
		return nil, "", error
	}

	branches := []string{}
	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		branch := strings.TrimSpace(scanner.Text())

		if strings.Contains(branch, "* ") {
			branch = strings.Replace(branch, "* ", "", 1)
			currentBranch = branch
		}

		// Skip aliaces origin/HEAD -> origin/main
		if strings.Contains(branch, "->") {
			continue
		}

		branches = append(branches, branch)
	}

	return branches, currentBranch, nil
}

// Удаление бранчей
func RemoveBranches(branches []string) error {
	args := append([]string{"branch", "-D"}, branches...)

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
func HasLocalBranch(branch string) bool {
	cmd := exec.Command("git", "branch", "--list", branch)
	output, _ := cmd.Output()

	return strings.TrimSpace(string(output)) != ""
}

func Contains(array []string, element string) bool {
	for _, current := range array {
		if current == element {
			return true
		}
	}

	return false
}

func Checkout(branch string) {
	checkoutCommand := exec.Command("git", "checkout", branch)
	checkoutCommand.Stdout = os.Stdout
	checkoutCommand.Stderr = os.Stderr

	checkoutError := checkoutCommand.Run()

	if checkoutError != nil {
		fmt.Println("🔴 Error:")
		fmt.Println(checkoutError)
	}
}

// Case: The current branch was selected to be removed
// We want to find another branch to checkout to be able to remove the current branch
// Search for an element that:
// exists in the first array and does not exist in the second

var mainBranch = "main"
var masterBranch = "master"

func GetCheckoutBranch(array *[]string, array2 *[]string) string {
	var branchesForCheckout []string

	for _, branch := range *array {
		if !slices.Contains(*array2, branch) {
			branchesForCheckout = append(branchesForCheckout, branch)
		}
	}

	// If "main" or "master" branches are available for checkout
	// Try to return them first
	for _, branch := range branchesForCheckout {
		if branch == mainBranch || branch == masterBranch {
			return branch
		}
	}

	if len(branchesForCheckout) == 0 {
		return ""
	}

	return branchesForCheckout[0]
}
