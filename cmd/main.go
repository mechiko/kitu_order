package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"kitu/app"
	"kitu/config"
	"kitu/gui"
	"kitu/process"
	"kitu/reductor"
	"kitu/repo"
	"kitu/zaplog"

	"github.com/mechiko/dbscan"
	"github.com/mechiko/utility"
	"go.uber.org/zap"
)

var order = flag.Int64("order", 0, "")

// если home true то папка создается локально
var home = flag.Bool("home", false, "")

var fileExe string
var dir string

const StartNumber = 7
const CountByPallet = 24

func init() {
	flag.Parse()
	fileExe = os.Args[0]
	var err error
	dir, err = filepath.Abs(filepath.Dir(fileExe))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get absolute path: %v\n", err)
		os.Exit(1)
	}
	if err := os.Chdir(dir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to change directory: %v\n", err)
		os.Exit(1)
	}
}

func errMessageExit(loger *zap.SugaredLogger, title string, err error) {
	if loger != nil {
		loger.Errorf("%s %v", title, err)
	}
	utility.MessageBox(title, err.Error())
	os.Exit(-1)
}

func main() {

	cfg, err := config.New("", *home)
	if err != nil {
		errMessageExit(nil, "ошибка конфигурации", err)
	}

	var logsOutConfig = map[string][]string{
		"logger": {"stdout", filepath.Join(cfg.LogPath(), config.Name)},
	}
	zl, err := zaplog.New(logsOutConfig, true)
	if err != nil {
		errMessageExit(nil, "ошибка создания логера", err)
	}

	lg, err := zl.GetLogger("logger")
	if err != nil {
		errMessageExit(nil, "ошибка получения логера", err)
	}
	loger := lg.Sugar()
	loger.Debug("zaplog started")
	loger.Infof("mode = %s", config.Mode)
	if cfg.Warning() != "" {
		loger.Infof("pkg:config warning %s", cfg.Warning())
	}

	// создаем приложение с опциями из конфига и логером основным
	app := app.New(cfg, loger, dir)
	// бд основные находятся в текущем каталоге если не переопределено в настройках
	app.SetDefaultDbPath("")

	// инициализируем пути необходимые приложению
	app.CreatePath()

	// инициализируем REPO
	// TODO изменить получение путей из конфига
	listDbs := make(dbscan.ListDbInfoForScan)
	listDbs[dbscan.Config] = &dbscan.DbInfo{}
	listDbs[dbscan.TrueZnak] = &dbscan.DbInfo{}

	repoStart, err := repo.New(app.Logger(), listDbs, app.DefaultDbPath())
	if err != nil {
		errMessageExit(loger, "Ошибки запуска репозитория", err)
	}
	err = app.SetRepo(repoStart)
	if err != nil {
		errMessageExit(loger, "Ошибки установки в app репозитория", err)
	}

	k, err := process.New(app, repoStart)
	if err != nil {
		errMessageExit(loger, "Ошибки установки в app репозитория", err)
	}

	model := reductor.Model{}
	model.Read(app)
	model.Order = *order
	if model.StartNumberSSCC <= 0 {
		model.StartNumberSSCC = 1
		if err := model.Sync(app); err != nil {
			errMessageExit(loger, "Ошибки записи файла конфигурации", err)
		}
	}

	// создаем редуктор с новой моделью
	reductor.New(model, app.Logger())
	gui.New(k, app).Run()
}
