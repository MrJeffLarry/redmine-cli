package terminal

import "golang.org/x/term"

func Blue(t *term.Terminal, text string) string {
	return string(t.Escape.Blue) + text + string(t.Escape.Reset)
}

func Green(t *term.Terminal, text string) string {
	return string(t.Escape.Green) + text + string(t.Escape.Reset)
}

func White(t *term.Terminal, text string) string {
	return string(t.Escape.White) + text + string(t.Escape.Reset)
}

func Black(t *term.Terminal, text string) string {
	return string(t.Escape.Black) + text + string(t.Escape.Reset)
}

func Cyan(t *term.Terminal, text string) string {
	return string(t.Escape.Cyan) + text + string(t.Escape.Reset)
}

func Magenta(t *term.Terminal, text string) string {
	return string(t.Escape.Magenta) + text + string(t.Escape.Reset)
}

func Red(t *term.Terminal, text string) string {
	return string(t.Escape.Red) + text + string(t.Escape.Reset)
}

func Yellow(t *term.Terminal, text string) string {
	return string(t.Escape.Yellow) + text + string(t.Escape.Reset)
}
