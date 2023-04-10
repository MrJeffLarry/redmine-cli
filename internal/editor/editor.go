package editor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
)

func StartEdit(exec, body string) string {
	edit := "nano"

	if exec != "" {
		edit = exec
	} else if runtime.GOOS == "windows" {
		edit = "notepad"
	} else if g := os.Getenv("GIT_EDITOR"); g != "" {
		edit = g
	} else if v := os.Getenv("VISUAL"); v != "" {
		edit = v
	} else if e := os.Getenv("EDITOR"); e != "" {
		edit = e
	}
	return editor(edit, body)
}

func StartPage(exec, body string) {
	var path string
	var err error
	var args []string

	view := ""

	if exec == "" {
		fmt.Println(body)
		return
	} else if exec != "" {
		view = exec
	} else if g := os.Getenv("LESS"); g != "" {
		view = "less"
		args = append(args, "-f")
	} else if runtime.GOOS == "windows" {
		view = "notepad"
	} else {
		// did not found any viewer
		fmt.Println(body)
		return
	}

	if path, err = config.CreateTmpFile(body); err != nil {
		print.Error(err.Error())
		return
	}

	args = append(args, path)

	if err = createFile(view, args, body); err != nil {
		print.Error(err.Error())
		return
	}

	if err = os.Remove(path); err != nil {
		print.Error(err.Error())
	}
}

func createFile(editor string, arg []string, body string) error {
	var err error

	cmd := exec.Command(editor, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Start(); err != nil {
		print.Error(err.Error())
		return err
	}

	if err = cmd.Wait(); err != nil {
		print.Error(err.Error())
		return err
	}
	return nil
}

func editor(editor, body string) string {
	var path string
	var data []byte
	var err error

	if path, err = config.CreateTmpFile(body); err != nil {
		print.Error(err.Error())
		return ""
	}

	if err = createFile(editor, []string{path}, body); err != nil {
		print.Error(err.Error())
		return ""
	}

	if data, err = os.ReadFile(path); err != nil {
		print.Error(err.Error())
		return ""
	}

	if err = os.Remove(path); err != nil {
		print.Error(err.Error())
	}
	return string(data)
}
