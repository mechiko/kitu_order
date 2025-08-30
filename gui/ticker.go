package gui

import (
	"fmt"
	"kitu/reductor"

	tk "modernc.org/tk9.0"
)

func (a *GuiApp) tick() {
	select {
	case s := <-a.logCh:
		if s.Error {
			a.logg("", s.Msg)
		} else {
			a.logg(s.Msg, "")
		}
	case <-a.logClear:
		a.logText.Configure(tk.State("normal"))
		a.logText.Delete("1.0", tk.END)
		a.logText.Configure(tk.State("disabled"))
	case <-a.stateStart:
		// состояние начала
		a.startButton.Configure(tk.State("enabled"))
		a.exitButton.Configure(tk.State("enabled"))
	case <-a.stateFinish:
		// состояние после пуск
		model := reductor.Instance().Model("")
		a.startNumberSscc.Configure(tk.Textvariable(fmt.Sprintf("%d", model.StartNumberSSCC)))
		a.startButton.Configure(tk.State("enabled"))
		a.exitButton.Configure(tk.State("enabled"))
	case a.isProces = <-a.stateIsProcess:
		if a.isProces {
			a.startButton.Configure(tk.State("disabled"))
		}
	default:
	}
}
