package gui

import (
	"fmt"
	"kitu/reductor"

	tk "modernc.org/tk9.0"
	"modernc.org/tk9.0/extensions/autoscroll"
)

func (a *GuiApp) makeWidgets(model *reductor.Model) {
	a.makeLog()
	a.makeInputs(model)
	a.makeButtons()
}

func (a *GuiApp) makeLog() {
	a.logFrame = tk.TFrame()
	a.yscroll = autoscroll.Autoscroll(tk.TScrollbar(tk.Command(func(e *tk.Event) { e.Yview(a.logText) })).Window)
	a.logText = a.logFrame.Text(
		// tk.Font(tk., 10),
		// tk.Font(tk.NewFont(tk.Family("Consolas"))),
		tk.Yscrollcommand(func(e *tk.Event) { e.ScrollSet(a.yscroll) }),
		tk.Setgrid(true),
		tk.Wrap(tk.WORD),
		tk.Padx("2m"),
		tk.Pady("2m"),
		tk.Height(10),
	)
	tag := fmt.Sprintf("t%v", 1)
	a.logText.TagConfigure(tag, tk.Foreground(tk.Red))
	a.logText.TagConfigure("bgstipple", tk.Background(tk.Black), tk.Borderwidth(0), tk.Bgstipple(tk.Gray12))
	a.logText.TagConfigure("big", tk.Font("helvetica", 12, "bold"))
	a.logText.TagConfigure("bold", tk.Font("helvetica", 10, "bold", "italic"))
	a.logText.TagConfigure("center", tk.Justify("center"))
	a.logText.TagConfigure("color1", tk.Foreground(tk.Blue))
	a.logText.TagConfigure("color2", tk.Foreground(tk.Red))
	a.logText.TagConfigure("margins", tk.Lmargin1("12m"), tk.Lmargin2("6m"), tk.Rmargin("10m"))
	a.logText.TagConfigure("overstrike", tk.Overstrike(1))
	a.logText.TagConfigure("raised", tk.Relief("raised"), tk.Borderwidth(1))
	a.logText.TagConfigure("right", tk.Justify("right"))
	a.logText.TagConfigure("spacing", tk.Spacing1("10p"), tk.Spacing2("2p"), tk.Lmargin1("12m"), tk.Lmargin2("6m"), tk.Rmargin("10m"))
	a.logText.TagConfigure("sub", tk.Offset("-2p"), tk.Font("helvetica", 8))
	a.logText.TagConfigure("sunken", tk.Relief("sunken"), tk.Borderwidth(1))
	a.logText.TagConfigure("super", tk.Offset("4p"), tk.Font("helvetica", 8))
	a.logText.TagConfigure("tiny", tk.Font("times", 8, "bold"))
	a.logText.TagConfigure("underline", tk.Underline(1))
	a.logText.TagConfigure("verybig", tk.Font(tk.CourierFont(), 22, "bold"))
}

func (a *GuiApp) makeInputs(model *reductor.Model) {
	a.inputFrame = tk.TFrame()
	a.entirely = a.inputFrame.Checkbutton(tk.Txt("агрегировать весь заказ в одну упаковку"), tk.Command(func() {
		// val := a.entirely.Variable()
		// a.logg(fmt.Sprintf("val checkbutton %s %T", val, val), "")
		model.Entirely = a.entirely.Variable() == "1"
		reductor.Instance().SetModel("", *model)
		if err := model.Sync(a); err != nil {
			a.Logger().Errorf("gui entirely Checkbutton синхронизация модели %v", err)
		}
	}))
	if model.Entirely {
		a.entirely.Configure(tk.Variable("1"))
	} else {
		a.entirely.Configure(tk.Variable("0"))
	}

	a.statusKM = a.inputFrame.Checkbutton(tk.Txt("агрегировать только нанесенные КМ"), tk.Command(func() {
		model.StatusKM = a.statusKM.Variable() == "1"
		reductor.Instance().SetModel("", *model)
		if err := model.Sync(a); err != nil {
			a.Logger().Errorf("gui statusKM Checkbutton синхронизация модели %v", err)
		}
	}))
	if model.StatusKM {
		a.statusKM.Configure(tk.Variable("1"))
	} else {
		a.statusKM.Configure(tk.Variable("0"))
	}

	a.order = a.inputFrame.TEntry(tk.Textvariable(fmt.Sprintf("%d", model.Order)))
	a.perPalet = a.inputFrame.TEntry(tk.Textvariable(fmt.Sprintf("%d", model.PerPallet)))
	a.startNumberSscc = a.inputFrame.TEntry(tk.Textvariable(fmt.Sprintf("%d", model.StartNumberSSCC)))
}

func (a *GuiApp) makeButtons() {
	a.buttonFrame = tk.TFrame()
	a.exitButton = a.buttonFrame.TExit(tk.Txt("Выход"))
	a.startButton = a.buttonFrame.TButton(tk.Txt("Пуск"), tk.Command(func() {
		if a.isProces {
			return
		}
		a.startButton.Configure(tk.State("disabled"))
		a.exitButton.Configure(tk.State("disabled"))
		go func() {
			defer func() {
				if r := recover(); r != nil {
					a.Logger().Error("panic go func")
				}
			}()
			a.generate()
			a.Logger().Info("exit go generate")
		}()

	}))
	a.configButton = a.buttonFrame.TButton(tk.Txt("Config…"),
		tk.Command(a.onConfig))

}
