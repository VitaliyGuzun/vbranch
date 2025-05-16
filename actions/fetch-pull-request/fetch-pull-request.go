package fetchpullrequest

import (
	"encoding/json"
	"fmt"
	"gh-api/actions/ask"
	"gh-api/actions/command"
	"gh-api/actions/shared"
	"log"
	"os/exec"
	"strconv"
)

var choosePullRequestLabel = "Choose pull request to fetch and checkout: "

func Run() {
	cmd := exec.Command("gh", "pr", "list", "--state", "open", "--json", "number,title,state,url,author,headRefName")
	output, err := cmd.Output()

	if err != nil {
		log.Fatal("Fetch pull requests: ", err)
	}

	var pullRequests []shared.PullRequest

	if err := json.Unmarshal(output, &pullRequests); err != nil {
		fmt.Println("JSON parse error:", err)
		return
	}

	fmt.Println("pullRequests: ", pullRequests)

	pullRequestTitle := make([]string, len(pullRequests))

	for i, item := range pullRequests {
		pullRequestTitle[i] = strconv.Itoa(item.Number)
	}

	pullRequest, err := ask.OneWithDescription(&pullRequestTitle, &choosePullRequestLabel, pullRequests)

	if err != nil {
		log.Fatal(choosePullRequestLabel, err)
	}

	// Здесь надо подставлять не title а номер бранча
	// There is not such a branch, fetch and checkout
	err = command.Run("gh", "pr", "checkout", pullRequest)

	if err != nil {
		log.Fatal("Fetch branch error: ", err)
	}
}
