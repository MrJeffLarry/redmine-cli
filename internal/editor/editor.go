package editor

import (
	"os"
	"os/exec"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
)

func StartEdit(body string) string {
	return editor("nano", body)
}

func StartView(body string) {
	viewer("less", body)
}

func createFile(editor string, arg []string, body string) error {
	var err error

	cmd := exec.Command(editor, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Start(); err != nil {
		print.Error("%s", err)
		return err
	}

	if err = cmd.Wait(); err != nil {
		print.Error("%s", err)
		return err
	}
	return nil
}

func viewer(viewer string, body string) {
	var path string
	var err error

	if path, err = config.CreateTmpFile(body); err != nil {
		print.Error("%s", err)
		return
	}

	if err = createFile(viewer, []string{"-f", path}, body); err != nil {
		print.Error("%s", err)
		return
	}

	if err = os.Remove(path); err != nil {
		print.Error("%s", err)
	}
}

func editor(editor, body string) string {
	var path string
	var data []byte
	var err error

	if path, err = config.CreateTmpFile(body); err != nil {
		print.Error("%s", err)
		return ""
	}

	if err = createFile(editor, []string{path}, body); err != nil {
		print.Error("%s", err)
		return ""
	}

	if data, err = os.ReadFile(path); err != nil {
		print.Error("%s", err)
		return ""
	}

	if err = os.Remove(path); err != nil {
		print.Error("%s", err)
	}
	return string(data)
}
