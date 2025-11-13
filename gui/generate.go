package gui

import (
	"fmt"
	"kitu/process/protocol"
	"kitu/reductor"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mechiko/utility"
)

func (a *GuiApp) generate() {
	defer func() {
		a.Logger().Debug("defer generate()")
		a.isProces = false
		// по завершению обработки в БД кнопка Пуск запрещена
		a.stateFinish <- struct{}{}
	}()
	a.isProces = true
	model := reductor.Instance().Model("")
	if order, err := strconv.ParseInt(a.order.Textvariable(), 10, 64); err != nil {
		a.SendError(fmt.Sprintf("ошибка номера заказа: %s", err.Error()))
		return
	} else {
		if order <= 0 {
			a.SendError("ошибка: номер заказа должен быть > 0")
			return
		}
		model.Order = order
	}
	if perPalet, err := strconv.ParseInt(a.perPalet.Textvariable(), 10, 64); err != nil {
		a.SendError(fmt.Sprintf("ошибка количества в палете: %s", err.Error()))
		return
	} else {
		model.PerPallet = int(perPalet)
		if perPalet <= 0 {
			a.SendError("ошибка: количество в палете должно быть > 0")
			return
		}
		// записываем модель
		reductor.Instance().SetModel("", model)
	}
	if startNumber, err := strconv.ParseInt(a.startNumberSscc.Textvariable(), 10, 64); err != nil {
		a.SendError(fmt.Sprintf("ошибка начального номера палеты: %s", err.Error()))
		return
	} else {
		model.StartNumberSSCC = int(startNumber)
		if startNumber <= 0 {
			a.SendError("ошибка: начальный номер палеты должен быть > 0")
			return
		}
		// записываем модель
		reductor.Instance().SetModel("", model)
	}

	// сброс
	a.krinica.Reset()
	if err := a.krinica.ReadOrder(); err != nil {
		a.SendError(fmt.Sprintf("ошибка получения КМ из заказа: %s", err.Error()))
		return
	}
	if len(a.krinica.Cis) == 0 {
		a.SendError(fmt.Sprintf("нет КМ из заказа %d", model.Order))
		return
	}
	a.SendLog(fmt.Sprintf("для заказа %d взято %d КМ", model.Order, len(a.krinica.Cis)))
	if err := a.krinica.GeneratePalletOrder(); err != nil {
		a.SendError(fmt.Sprintf("ошибка генерации палет: %s", err.Error()))
		return
	}
	a.SendLog(fmt.Sprintf("сгенерировано %d палет по %d шт.", len(a.krinica.Pallet), model.PerPallet))

	if err := a.krinica.WritePaletsForce(); err != nil {
		a.SendError(fmt.Sprintf("ошибка записи в БД: %s", err.Error()))
		return
	}
	a.SendLog("агрегация записана в базу А3")

	model = reductor.Instance().Model("")
	// запишем значения следующей SSCC в конфиг
	model.StartNumberSSCC = model.LastSSCC + 1
	// model.StartNumberSSCC = model.LastSSCC
	if err := model.Sync(a); err != nil {
		a.SendError(fmt.Sprintf("ошибка model.Sync конфигурации: %s", err.Error()))
		return
	}
	model.Date = time.Now().Format("2006-01-02")
	// запоминаем модель уже когда палеты внесены в бд
	reductor.Instance().SetModel("", model)

	// если ошибка формирования протокола, то обращаться ко мне...
	if fileTxt, err := protocol.PrintKrinicaProtocol(a.krinica); err != nil {
		a.SendError(fmt.Sprintf("ошибка формирования протокола: %s", err.Error()))
		return
	} else {
		fileName := fmt.Sprintf("pallet_%s_%s.html", time.Now().Format("2006-01-02"), utility.String(6))
		dir := filepath.Join(utility.UserHomeDir(), a.LogPath(), a.Options().Output)
		if filepath.IsAbs(a.Options().Output) {
			dir = a.Options().Output
		}
		url := filepath.Join(dir, fileName)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			a.SendError(fmt.Sprintf("ошибка создания каталога %q: %s", dir, err.Error()))
			return
		}
		if err := os.WriteFile(url, fileTxt, 0o644); err != nil {
			a.SendError(fmt.Sprintf("ошибка записи файла %q: %s", url, err.Error()))
			return
		}
		if err := OpenFile(url); err != nil {
			a.SendError(fmt.Sprintf("ошибка открытия файла: %s", err.Error()))
			return
		}
	}
}
