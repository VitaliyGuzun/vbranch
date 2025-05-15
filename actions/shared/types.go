package shared

type PullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	Url    string `json:"url"`
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
}
