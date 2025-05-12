package main

import (
	"gh-api/actions/ask"
	"gh-api/actions/fetch"
	"gh-api/actions/remove"
	"gh-api/utilities"
	"log"
)

/*
TODO:
- добавить тесты к каждой функции
- при выполнении скрипта, проверять что есть обновления и предлагать обновить
- добавить логи в каждую функцию, чтобы юзер видел что происходит
*/

var chooseActionLabel = "Choose action: "

func main() {
	if error := utilities.IsGitRepo(); error != nil {
		log.Fatal("Git is not inited", error)
	}

	actions := []string{"fetch", "remove"}
	action, err := ask.One(&actions, &chooseActionLabel)

	if err != nil {
		log.Fatal(chooseActionLabel, err)
	}

	if action == "fetch" {
		fetch.Run()
	} else if action == "remove" {
		remove.Run()
	}
}
