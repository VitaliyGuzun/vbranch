package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// Проверка, что мы в git-репозитории
func checkGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()

	if err != nil || strings.TrimSpace(string(output)) != "true" {
		return fmt.Errorf("not a git repository")
	}

	return nil
}

func main() {
	error := checkGitRepo()

	if error != nil {
		log.Fatalf("ERROR: %v", error)
	}

	log.Println("Все ок")
}
