package gui

import (
	_ "embed"
	"kitu/domain"
	"kitu/process"
	"kitu/reductor"
	"time"

	tk "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

const (
	logRotate = 1000
	tick      = 10 * time.Millisecond
)

//go:embed 192.png
var ico []byte

type GuiApp struct {
	domain.Apper
	icon           *tk.Img
	logText        *tk.TextWidget
	buttonFrame    *tk.TFrameWidget
	inputFrame     *tk.TFrameWidget
	logFrame       *tk.TFrameWidget
	startButton    *tk.TButtonWidget
	exitButton     *tk.TButtonWidget
	configButton   *tk.TButtonWidget
	logCh          chan LogMsg
	stateFinish    chan struct{}
	stateStart     chan struct{}
	logClear       chan struct{}
	stateIsProcess chan bool

	order           *tk.TEntryWidget
	perPalet        *tk.TEntryWidget
	startNumberSscc *tk.TEntryWidget
	krinica         *process.Krinica
	yscroll         *tk.Window
	entirely        *tk.CheckbuttonWidget
	statusKM        *tk.CheckbuttonWidget
	isProces        bool
}

func New(k *process.Krinica, app domain.Apper) *GuiApp {
	a := &GuiApp{
		Apper:   app,
		krinica: k,
	}
	a.logCh = make(chan LogMsg, 10)
	a.stateFinish = make(chan struct{})
	a.stateStart = make(chan struct{})
	a.icon = tk.NewPhoto(tk.Data(ico))
	a.logClear = make(chan struct{})
	a.stateIsProcess = make(chan bool, 2)

	tk.App.IconPhoto(a.icon)
	tk.ErrorMode = tk.CollectErrors
	tk.App.WmTitle("Агрегация по номеру заказа КМ")
	tk.WmProtocol(tk.App, "WM_DELETE_WINDOW", a.onQuitApp)
	if err := tk.ActivateTheme("azure light"); err != nil {
		a.Logger().Errorf("gui theme %s", err.Error())
	}
	tk.InitializeExtension("autoscroll")

	tk.NewTicker(tick, a.tick)
	// tk.NewTicker(tack, a.tack)

	model := reductor.Instance().Model("")
	a.makeWidgets(&model)
	a.makeLayout()
	a.makeBindings()
	return a
}

func (a *GuiApp) Run() {
	tk.App.Center()
	tk.WmDeiconify(tk.App)
	tk.App.Wait()
}

func (a *GuiApp) onQuitApp() {
	if a.isProces {
		a.logg("", "выход из программы ограничен, запущена обработка")
		return
	}
	tk.Destroy(tk.App)
}
