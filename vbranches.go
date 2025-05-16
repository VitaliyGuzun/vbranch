package main

import (
	"fmt"
	"gh-api/actions/ask"
	"gh-api/actions/branch"
	pullrequest "gh-api/actions/pull-request"
	"gh-api/actions/remove"
	"gh-api/utilities"
	"log"
)

/*
TODO:
- залупить экшены. После выполнения команды, надо возвращаться обратно к списку команд
- возвращать после удаления на main / master
- добавить тесты к каждой функции
- при выполнении скрипта, проверять что есть обновления и предлагать обновить
- добавить логи в каждую функцию, чтобы юзер видел что происходит
*/

var chooseActionLabel = "Choose action:"
var BRANCH = "branch"
var PULL_REQUEST = "pull request"
var REMOVE_BRANCH = "remove branch"

func main() {
	if error := utilities.IsGitRepo(); error != nil {
		log.Fatal("Git is not inited", error)
	}

	actions := []string{BRANCH, PULL_REQUEST, REMOVE_BRANCH}

	for {
		action, err := ask.One(&actions, &chooseActionLabel)

		if err != nil {
			log.Fatal(chooseActionLabel, err)
		}

		if action == BRANCH {
			branch.Run()
		} else if action == PULL_REQUEST {
			pullrequest.Run()
		} else if action == REMOVE_BRANCH {
			remove.Run()
		}

		fmt.Println()
		action = ""
	}
}
