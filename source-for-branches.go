// package main

// import (
// 	"bufio"
// 	"bytes"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/exec"
// 	"strings"

// 	"github.com/AlecAivazis/survey/v2"
// )

// // Проверка, что мы в git-репозитории
// func checkGitRepo() error {
//     cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
//     output, err := cmd.Output()
//     if err != nil || strings.TrimSpace(string(output)) != "true" {
//         return fmt.Errorf("not a git repository")
//     }
//     return nil
// }

// // Получение всех удалённых веток
// func getRemoteBranches() ([]string, error) {
//     cmd := exec.Command("git", "branch", "-r")
//     output, err := cmd.Output()
//     if err != nil {
//         return nil, err
//     }
//     branches := []string{}
//     scanner := bufio.NewScanner(bytes.NewReader(output))
//     for scanner.Scan() {
//         branch := strings.TrimSpace(scanner.Text())
//         if strings.Contains(branch, "->") {
//             // Пропускаем алиасы типа origin/HEAD -> origin/main
//             continue
//         }
//         branches = append(branches, branch)
//     }
//     return branches, nil
// }

// // Проверка, есть ли локальная ветка
// func hasLocalBranch(branch string) bool {
//     branch = strings.TrimPrefix(branch, "origin/")
//     cmd := exec.Command("git", "branch", "--list", branch)
//     output, _ := cmd.Output()
//     return strings.TrimSpace(string(output)) != ""
// }

// func main() {
//     // 1. Проверяем git-репозиторий
//     if err := checkGitRepo(); err != nil {
//         log.Fatalf("Ошибка: %v", err)
//     }

//     // 2. fetch --all
//     fmt.Println("Скачиваем все ветки с origin...")
//     fetchCmd := exec.Command("git", "fetch", "--all")
//     fetchCmd.Stdout = os.Stdout
//     fetchCmd.Stderr = os.Stderr
//     if err := fetchCmd.Run(); err != nil {
//         log.Fatalf("Ошибка fetch: %v", err)
//     }

//     // 3. Получаем список веток
//     branches, err := getRemoteBranches()
//     if err != nil {
//         log.Fatalf("Ошибка получения веток: %v", err)
//     }
//     if len(branches) == 0 {
//         log.Fatalf("Нет веток на origin")
//     }

//     // 4. Выбор ветки
//     var selected string
//     prompt := &survey.Select{
//         Message: "Выбери ветку для переключения:",
//         Options: branches,
//     }
//     if err := survey.AskOne(prompt, &selected); err != nil {
//         log.Fatalf("Ошибка выбора: %v", err)
//     }

//     // 5. Переключение на выбранную ветку
//     localBranch := strings.TrimPrefix(selected, "origin/")
//     if hasLocalBranch(selected) {
//         // Ветка уже есть локально
//         checkoutCmd := exec.Command("git", "checkout", localBranch)
//         checkoutCmd.Stdout = os.Stdout
//         checkoutCmd.Stderr = os.Stderr
//         if err := checkoutCmd.Run(); err != nil {
//             log.Fatalf("Ошибка checkout: %v", err)
//         }
//     } else {
//         // Ветки нет локально, создаём её отслеживая origin
//         checkoutCmd := exec.Command("git", "checkout", "-b", localBranch, selected)
//         checkoutCmd.Stdout = os.Stdout
//         checkoutCmd.Stderr = os.Stderr
//         if err := checkoutCmd.Run(); err != nil {
//             log.Fatalf("Ошибка создания ветки: %v", err)
//         }
//     }

//     fmt.Printf("Переключено на ветку %s\n", localBranch)
// }
