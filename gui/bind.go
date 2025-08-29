package gui

import (
	"kitu/reductor"
	"strconv"

	tk "modernc.org/tk9.0"
)

func (a *GuiApp) makeBindings() {
	tk.Bind(tk.App, "<Escape>", tk.Command(a.onQuitApp))
	tk.Bind(a.order, "<KeyRelease>", tk.Command(func(e *tk.Event) {
		order, err := strconv.ParseInt(a.order.Textvariable(), 10, 64)
		if err != nil {
			a.order.Configure(tk.Textvariable(""))
		}
		model := reductor.Instance().Model("")
		model.Order = order
		reductor.Instance().SetModel("", model)
		a.Logger().Debugf("order %d", order)
	}))
	tk.Bind(a.startNumberSscc, "<KeyRelease>", tk.Command(func(e *tk.Event) {
		start, err := strconv.ParseInt(a.startNumberSscc.Textvariable(), 10, 64)
		if err != nil {
			a.startNumberSscc.Configure(tk.Textvariable(""))
		}
		model := reductor.Instance().Model("")
		model.StartNumberSSCC = int(start)
		reductor.Instance().SetModel("", model)
		a.Logger().Debugf("startNumberSscc %d", start)
	}))
	tk.Bind(a.perPalet, "<KeyRelease>", tk.Command(func(e *tk.Event) {
		per, err := strconv.ParseInt(a.perPalet.Textvariable(), 10, 64)
		if err != nil {
			a.perPalet.Configure(tk.Textvariable(""))
		}
		model := reductor.Instance().Model("")
		model.PerPallet = int(per)
		reductor.Instance().SetModel("", model)
		a.Logger().Debugf("perPalet %d", per)
	}))
}
