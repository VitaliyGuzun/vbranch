package ask

import (
	"gh-api/actions/shared"

	"github.com/AlecAivazis/survey/v2"
)

// Hello
// hello2

func One(options *[]string, label *string) (string, error) {
	var value string

	actionsPrompt := &survey.Select{
		Message: *label,
		Options: *options,
	}

	err := survey.AskOne(actionsPrompt, &value)

	if err != nil {
		return "", err
	}

	return value, nil
}

func Many(options *[]string, label *string) ([]string, error) {
	var value []string

	actionsPrompt := &survey.MultiSelect{
		Message: *label,
		Options: *options,
	}

	err := survey.AskOne(actionsPrompt, &value)

	if err != nil {
		return nil, err
	}

	return value, nil
}

func OneMultiline() {
	text := ""
	prompt := &survey.Multiline{
		Message: "ping\nHello",
	}
	survey.AskOne(prompt, &text)
}

func OneWithDescription(options *[]string, label *string, meta []shared.PullRequest) (string, error) {
	var value string

	actionsPrompt := &survey.Select{
		Message: *label,
		Options: *options,
		Description: func(value string, index int) string {
			return "by " + meta[index].Author.Login + ", branch: " + meta[index].Branch
		},
	}

	err := survey.AskOne(actionsPrompt, &value)

	if err != nil {
		return "", err
	}

	return value, nil
}
