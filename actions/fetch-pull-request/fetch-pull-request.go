package fetchpullrequest

import (
	"encoding/json"
	"fmt"
	"gh-api/actions/ask"
	"gh-api/actions/shared"
	"log"
	"os/exec"
)

var choosePullRequestLabel = "Choose pull request to fetch and checkout: "

func Run() {
	cmd := exec.Command("gh", "pr", "list", "--state", "open", "--json", "number,title,state,url,author")
	output, err := cmd.Output()

	if err != nil {
		log.Fatal("Fetch pull requests: ", err)
	}

	var pullRequests []shared.PullRequest

	if err := json.Unmarshal(output, &pullRequests); err != nil {
		fmt.Println("JSON parse error:", err)
		return
	}

	pullRequestTitle := make([]string, len(pullRequests))

	for i, item := range pullRequests {
		pullRequestTitle[i] = item.Title
	}

	pullRequest, err := ask.OneWithDescription(&pullRequestTitle, &choosePullRequestLabel, pullRequests)

	if err != nil {
		log.Fatal(choosePullRequestLabel, err)
	}

	fmt.Println("pullRequest: ", pullRequest)
}
