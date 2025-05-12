package ask

import "github.com/AlecAivazis/survey/v2"

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
