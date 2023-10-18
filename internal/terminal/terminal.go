package terminal

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/MrJeffLarry/redmine-cli/internal/util"
)

type Terminal struct {
	Stdin  *os.File
	Stdout *os.File
	Stderr *os.File
}

func New(stdin *os.File, stdout *os.File, stderr *os.File) *Terminal {
	t := &Terminal{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if stdin == nil {
		t.Stdin = os.Stdin
	}

	if stdout == nil {
		t.Stdout = os.Stdout
	}

	if stderr == nil {
		t.Stderr = os.Stderr
	}

	return t
}

func (t *Terminal) Choose(label string, chooses []util.IdName) (int, string) {
	options := make([]string, len(chooses))

	for i, m := range chooses {
		options[i] = m.Name
	}

	choose := &survey.Select{
		Message: label + ":",
		Options: options,
	}

	stdio := survey.WithStdio(t.Stdin, t.Stdout, t.Stderr)

	index := 0

	if err := survey.AskOne(choose, &index, stdio); err != nil {
		if err == terminal.InterruptErr {
			os.Exit(0)
			return -1, ""
		}
		fmt.Printf("Prompt failed %v\n", err)
		return -1, ""
	}

	return chooses[index].ID, chooses[index].Name
}

func (t *Terminal) ChooseString(label string, chooses []string) (string, int) {

	choose := &survey.Select{
		Message: label + ":",
		Options: chooses,
	}

	stdio := survey.WithStdio(t.Stdin, t.Stdout, t.Stderr)

	index := 0

	if err := survey.AskOne(choose, &index, stdio); err != nil {

		if err == terminal.InterruptErr {
			os.Exit(0)
			return "", -1
		}
		fmt.Printf("Prompt failed %v\n", err)
		return "", -1
	}

	return chooses[index], index
}

func (t *Terminal) PromptPassword(label string, def string) (string, error) {
	pass := ""
	ask := &survey.Password{
		Message: label,
	}

	stdio := survey.WithStdio(t.Stdin, t.Stdout, t.Stderr)

	if err := survey.AskOne(ask, &pass, stdio); err != nil {

		if err == terminal.InterruptErr {
			os.Exit(0)
			return def, nil
		}
		fmt.Printf("Prompt failed %v\n", err)
		return def, err
	}

	return pass, nil
}

func (t *Terminal) PromptStringRequire(label string, def string) (string, error) {
	ask := &survey.Input{
		Message: label + ":",
		Default: def,
	}

	stdio := survey.WithStdio(t.Stdin, t.Stdout, t.Stderr)

	resp := ""

	if err := survey.AskOne(ask, &resp, survey.WithValidator(survey.Required), stdio); err != nil {

		if err == terminal.InterruptErr {
			os.Exit(0)
			return def, nil
		}
		fmt.Printf("Prompt failed %v\n", err)
		return def, err
	}

	return resp, nil
}

func (t *Terminal) PromptString(label string, def string) (string, error) {
	ask := &survey.Input{
		Message: label + ":",
		Default: def,
	}

	stdio := survey.WithStdio(t.Stdin, t.Stdout, t.Stderr)

	resp := ""

	if err := survey.AskOne(ask, &resp, survey.WithValidator(survey.Required), stdio); err != nil {

		if err == terminal.InterruptErr {
			os.Exit(0)
			return def, nil
		}
		fmt.Printf("Prompt failed %v\n", err)
		return def, err
	}

	return resp, nil
}

func (t *Terminal) PromptInt(label string, def int) (int, error) {
	var resp int

	ask := &survey.Input{
		Message: label + ":",
		Default: "-1",
	}

	stdio := survey.WithStdio(t.Stdin, t.Stdout, t.Stderr)

	if err := survey.AskOne(ask, &resp, survey.WithValidator(survey.Required), stdio); err != nil {

		if err == terminal.InterruptErr {
			os.Exit(0)
			return def, nil
		}
		fmt.Printf("Prompt failed %v\n", err)
		return def, err
	}

	return resp, nil
}

func (t *Terminal) Confirm(label string) bool {
	confirm := false

	prompt := &survey.Confirm{
		Message: label,
	}

	stdio := survey.WithStdio(t.Stdin, t.Stdout, t.Stderr)

	if err := survey.AskOne(prompt, &confirm, stdio); err != nil {

		if err == terminal.InterruptErr {
			os.Exit(0)
			return confirm
		}
		return confirm
	}
	return confirm
}
