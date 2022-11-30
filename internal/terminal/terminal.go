package terminal

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/MrJeffLarry/redmine-cli/internal/util"
)

func Choose(label string, chooses []util.IdName) (int, string) {
	options := make([]string, len(chooses))

	for i, m := range chooses {
		options[i] = m.Name
	}

	choose := &survey.Select{
		Message: label + ":",
		Options: options,
	}

	index := 0

	if err := survey.AskOne(choose, &index); err != nil {
		if err == terminal.InterruptErr {
			os.Exit(0)
			return -1, ""
		}
		fmt.Printf("Prompt failed %v\n", err)
		return -1, ""
	}

	return chooses[index].ID, chooses[index].Name
}

func ChooseString(label string, chooses []string) (string, int) {

	choose := &survey.Select{
		Message: label + ":",
		Options: chooses,
	}

	index := 0

	if err := survey.AskOne(choose, &index); err != nil {
		if err == terminal.InterruptErr {
			os.Exit(0)
			return "", -1
		}
		fmt.Printf("Prompt failed %v\n", err)
		return "", -1
	}

	return chooses[index], index
}

func PromptPassword(label string, def string) (string, error) {
	pass := ""
	ask := &survey.Password{
		Message: label,
	}

	if err := survey.AskOne(ask, &pass); err != nil {
		if err == terminal.InterruptErr {
			os.Exit(0)
			return def, nil
		}
		fmt.Printf("Prompt failed %v\n", err)
		return def, err
	}

	return pass, nil
}

func PromptStringRequire(label string, def string) (string, error) {
	ask := &survey.Input{
		Message: label + ":",
		Default: def,
	}

	resp := ""

	if err := survey.AskOne(ask, &resp, survey.WithValidator(survey.Required)); err != nil {
		if err == terminal.InterruptErr {
			os.Exit(0)
			return def, nil
		}
		fmt.Printf("Prompt failed %v\n", err)
		return def, err
	}

	return resp, nil
}

func PromptString(label string, def string) (string, error) {
	ask := &survey.Input{
		Message: label + ":",
		Default: def,
	}

	resp := ""

	if err := survey.AskOne(ask, &resp, survey.WithValidator(survey.Required)); err != nil {
		if err == terminal.InterruptErr {
			os.Exit(0)
			return def, nil
		}
		fmt.Printf("Prompt failed %v\n", err)
		return def, err
	}

	return resp, nil
}

func PromptInt(label string, def int) (int, error) {
	var resp int

	ask := &survey.Input{
		Message: label + ":",
		Default: "-1",
	}

	if err := survey.AskOne(ask, &resp, survey.WithValidator(survey.Required)); err != nil {
		if err == terminal.InterruptErr {
			os.Exit(0)
			return def, nil
		}
		fmt.Printf("Prompt failed %v\n", err)
		return def, err
	}

	return resp, nil
}

func Confirm(label string) bool {
	confirm := false

	prompt := &survey.Confirm{
		Message: label,
	}

	if err := survey.AskOne(prompt, &confirm); err != nil {
		if err == terminal.InterruptErr {
			os.Exit(0)
			return confirm
		}
		return confirm
	}
	return confirm
}
