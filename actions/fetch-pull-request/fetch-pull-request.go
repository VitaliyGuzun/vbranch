package fetchpullrequest

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type PullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	Url    string `json:"url"`
}

func Run() {
	cmd := exec.Command("gh", "pr", "list", "--state", "open", "--json", "number,title,state,url")
	output, err := cmd.Output()

	if err != nil {
		log.Fatal("Fetch pull requests: ", err)
	}

	var pullRequests []PullRequest

	if err := json.Unmarshal(output, &pullRequests); err != nil {
		fmt.Println("JSON parse error:", err)
		return
	}

	for _, pr := range pullRequests {
		fmt.Printf("#%d: %s (%s)\nURL: %s\n", pr.Number, pr.Title, pr.State, pr.Url)
	}
}
