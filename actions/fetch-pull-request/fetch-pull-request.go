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

var label = "Choose pull request to fetch and checkout: "

func Run() {
	cmd := exec.Command("gh", "pr", "list", "--state", "open", "--json", "number,title,state,url,author,headRefName")
	output, err := cmd.Output()

	if err != nil {
		log.Fatal("Fetch pull requests: ", err)
	}

	var prs []shared.PullRequest

	if err := json.Unmarshal(output, &prs); err != nil {
		fmt.Println("JSON parse error:", err)
		return
	}

	ids := make([]string, len(prs))

	for i, item := range prs {
		ids[i] = strconv.Itoa(item.Number)
	}

	pullRequest, err := ask.OneWithDescription(&ids, &label, prs)

	if err != nil {
		log.Fatal(label, err)
	}

	err = command.Run("gh", "pr", "checkout", pullRequest)

	if err != nil {
		log.Fatal("Fetch branch error: ", err)
	}
}
