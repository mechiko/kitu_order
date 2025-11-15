package dconfig

import (
	"strings"

	tk "modernc.org/tk9.0"
)

type ConfigDialogData struct {
	Ok         bool
	Inn        string
	PrefixSSCC string
}

type ConfigDialog struct {
	data *ConfigDialogData
	win  *tk.ToplevelWidget

	inn        *tk.TEntryWidget
	prefixSSCC *tk.TEntryWidget

	buttonFrame  *tk.TFrameWidget
	inputFrame   *tk.TFrameWidget
	okButton     *tk.TButtonWidget
	cancelButton *tk.TButtonWidget
}

func NewConfigDialog(data *ConfigDialogData) *ConfigDialog {
	dlg := &ConfigDialog{data: data}
	dlg.win = tk.App.Toplevel()
	dlg.win.WmTitle("Config")
	// tk.WmAttributes(dlg.win, tk.Type("dialog")) // TODO
	tk.WmProtocol(dlg.win.Window, tk.WM_DELETE_WINDOW, dlg.onCancel)

	dlg.makeWidgets()
	dlg.makeLayout()
	return dlg
}

func (me *ConfigDialog) onOk() {
	me.data.Ok = true
	me.data.PrefixSSCC = strings.Trim(me.prefixSSCC.Textvariable(), " ")
	tk.Destroy(me.win)
}

func (me *ConfigDialog) onCancel() {
	tk.Destroy(me.win)
}

func (me *ConfigDialog) ShowModal() {
	me.win.Raise(tk.App)
	tk.Focus(me.win)
	// tk.Focus(me.percentSpinbox)
	tk.GrabSet(me.win)
	me.win.Center().Wait()
}
