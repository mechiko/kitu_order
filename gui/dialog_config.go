package gui

import (
	"kitu/gui/dconfig"
	"kitu/reductor"
)

func (a *GuiApp) onConfig() {
	model := reductor.Instance().Model("")
	data := dconfig.ConfigDialogData{
		Inn:        model.Inn,
		PrefixSSCC: model.PrefixSSCC,
	}
	dlg := dconfig.NewConfigDialog(&data)
	dlg.ShowModal()
	if data.Ok {
		model.Inn = data.Inn
		model.PrefixSSCC = data.PrefixSSCC
		if err := model.Sync(a); err != nil {
			a.Logger().Errorf("диалог onConfig синхронизация модели %v", err)
		}
	}
	reductor.Instance().SetModel("", model)
}
