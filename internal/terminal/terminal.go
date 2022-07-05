package terminal

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/MrJeffLarry/redmine-cli/internal/util"
	"golang.org/x/term"
)

func input(pre string) (string, error) {
	if !term.IsTerminal(0) || !term.IsTerminal(1) {
		return "", fmt.Errorf("stdin/stdout should be terminal")
	}
	oldState, err := term.MakeRaw(0)
	if err != nil {
		return "", err
	}
	defer term.Restore(0, oldState)
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	term := term.NewTerminal(screen, "")
	term.AutoCompleteCallback = func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
		switch key {
		//		case 'T':
		//			line += "Task"
		default:
			line += string(key)
			pos++
		}
		return line, pos, true
	}
	term.SetPrompt(Green(term, pre))

	line, err := term.ReadLine()
	if err == io.EOF {
		return line, err // have to restore terminal before exit
	}
	if err != nil {
		return "", err
	}
	return line, nil
}

func WriteLine(pre string) string {
	text, err := input(pre + " ")
	if err == io.EOF {
		os.Exit(0) // need to restore before we exit
	}
	if err != nil {
		return ""
	}

	//	var screen *bytes.Buffer = new(bytes.Buffer)
	//	var output *bufio.Writer = bufio.NewWriter(os.Stdout)

	//	reader := bufio.NewReader(os.Stdin)
	//	fmt.Print(pre, ": ")
	//	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, " \n")
	return text
}

/*
// Move cursor to given position
func moveCursor(x int, y int) {
	fmt.Fprintf(screen, "\033[%d;%dH", x, y)
}

// Clear the terminal
func clearTerminal() {
	output.WriteString("\033[2J")
}

*/
func WriteLineReq(pre string, length int) string {
	for {
		value := WriteLine(pre)
		if len(value) > length {
			return value
		}
		print.Error("%s require a length of %d or more", pre, length)
	}
}

func WriteChooseIdName(pre string, chooses []util.IdName) (int64, string) {

	fmt.Printf("Choose %s\n", pre)
	for _, choose := range chooses {
		fmt.Printf("-> %s\n", choose.Name)
	}

	for {
		value := WriteLine(pre)
		for _, choose := range chooses {
			if strings.Compare(value, choose.Name) == 0 {
				return choose.ID, choose.Name
			}
		}
		print.Error("%s does not exist, please choose from list above", value)
	}
}

func WriteChooseString(pre string, chooses []string) string {
	for {
		value := WriteLine(pre)
		for _, choose := range chooses {
			if strings.Compare(value, choose) == 0 {
				return choose
			}
		}
		print.Error("%s does not exist, please choose from list above", value)
	}
}

func Confirm(text string) bool {
	for {
		writeBody := WriteLine(text + " (y/n)")
		if strings.Contains(writeBody, "y") {
			return true
		} else if strings.Contains(writeBody, "n") {
			return false
		} else {
			print.Error("%s: %s", "No valid input, valid (y=yes or n=no)", writeBody)
		}
	}
}
