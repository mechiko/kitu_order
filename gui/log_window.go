package gui

import (
	"fmt"
	"strconv"

	tk "modernc.org/tk9.0"
)

type LogMsg struct {
	Error bool
	Msg   string
}

func (a *GuiApp) logg(s, e string) {
	blue := "color1"
	red := "color2"
	if s != "" {
		s += "\n"
	}
	if e != "" {
		e += "\n"
	}
	a.logText.Configure(tk.State("normal"))
	a.logText.Insert(tk.END, s, blue, e, red)

	lines := a.logText.Count(tk.Lines(), "1.0", "end")
	if len(lines) == 1 {
		if n, err := strconv.ParseInt(lines[0], 10, 32); err == nil {
			if n > logRotate {
				a.logText.Delete("1.0", fmt.Sprintf("%v.0", n-logRotate))
			}
		}
	}
	a.logText.See("end")
	a.logText.Configure(tk.State("disabled"))
}

// вызывать из gorutine
// из основного потока вызывать только как go
func (a *GuiApp) SendError(s string) {
	// non-blocking send; drop if buffer full
	msg := LogMsg{
		Error: true,
		Msg:   s,
	}
	a.Logger().Error(s)
	select {
	case a.logCh <- msg:
	default:
		// drop
	}
}

// вызывать из gorutine
// из основного потока вызывать только как go
func (a *GuiApp) SendLog(s string) {
	// non-blocking send; drop if buffer full
	msg := LogMsg{
		Error: false,
		Msg:   s,
	}
	select {
	case a.logCh <- msg:
	default:
		// drop
	}
}
