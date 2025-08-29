package dconfig

import tk "modernc.org/tk9.0"

func (me *ConfigDialog) makeLayout() {
	me.layoutInputs()
	tk.Grid(me.inputFrame, tk.Row(1), tk.Column(0), tk.Pady(5), tk.Sticky(tk.WE))
	me.layoutButton()
	tk.Grid(me.buttonFrame, tk.Row(2), tk.Column(0), tk.Sticky(tk.WE))

	// tk.GridColumnConfigure(tk.App, 0, tk.Weight(1))
	// tk.GridRowConfigure(tk.App, 0, tk.Weight(1))
	// tk.App.Configure(tk.Padx(5), tk.Pady(5))
}

func (me *ConfigDialog) layoutInputs() {
	tk.Grid(me.inputFrame.Label(tk.Txt("Префикс упаковки:")), tk.Row(2), tk.Column(0), tk.Sticky(tk.W))
	tk.Grid(me.prefixSSCC, tk.Row(2), tk.Column(1), tk.Sticky(tk.W))
	// tk.Grid(me.inputFrame.Label(tk.Txt("ИНН:")), tk.Row(3), tk.Column(0), tk.Sticky(tk.W))
	// tk.Grid(me.inn, tk.Row(3), tk.Column(1), tk.Sticky(tk.W))

	// tk.GridColumnConfigure(me.inputFrame, 1, tk.Weight(2))
	// tk.GridRowConfigure(me.inputFrame, 0, tk.Weight(2))
}

func (me *ConfigDialog) layoutButton() {
	opts := tk.Opts{tk.Padx(3), tk.Pady(3)}
	tk.Grid(me.buttonFrame, tk.Row(1), tk.Column(0), tk.Columnspan(2), opts)
	tk.Grid(me.okButton, tk.Row(0), tk.Column(0), tk.Sticky(tk.E), opts)
	tk.Grid(me.cancelButton, tk.Row(0), tk.Column(1), tk.Sticky(tk.E), opts)
	tk.GridColumnConfigure(me.win, 1, tk.Weight(1))
}
