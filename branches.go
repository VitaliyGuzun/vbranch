package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// Проверка, что мы в git-репозитории
func isGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()

	if err != nil || strings.TrimSpace(string(output)) != "true" {
		return fmt.Errorf("not a git repository")
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

func main() {
	error := isGitRepo()

	if error != nil {
		log.Fatalf("ERROR: %v", error)
	}

	branches, branchesError := getRemoteBranches()

	if branchesError != nil {
		log.Fatalf("Ошибка получения веток: %v", branchesError)
	} else if len(branches) == 0 {
		log.Fatalf("Нет веток на origin")
	}

	fmt.Println(branches)

	// 4. Выбор ветки
	var selected string

	prompt := &survey.Select{
		Message: "Выбери ветку для переключения:",
		Options: branches,
	}

	branchError := survey.AskOne(prompt, &selected)

	if branchError != nil {
		log.Fatalf("Ошибка выбора: %v", branchError)
	}

	log.Println("Все ок")
}
