package process

import (
	"errors"
	"fmt"
	"kitu/reductor"
	"kitu/repo/znakdb"

	"github.com/mechiko/dbscan"
	"github.com/mechiko/utility"
	"github.com/upper/db/v4"
)

// start начальный номер SSCC палетты
// count количество в одной палетте
func (k *Krinica) GeneratePalletOrder() error {
	dbRepo, err := k.repo.Lock(dbscan.TrueZnak)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer k.repo.Unlock(dbRepo)

	dbZnak, ok := dbRepo.(*znakdb.DbZnak)
	if !ok {
		return fmt.Errorf("db type wrong %T", dbRepo)
	}

	indexPallet := 0
	model := reductor.Instance().Model("")
	startSSCC := model.StartNumberSSCC
	lastSSCC := model.StartNumberSSCC
	for {
		cis := make([]*utility.CisInfo, 0)
		if model.Entirely {
			cis = append(cis, k.Cis...)
		} else {
			cis = nextRecords(k.Cis, indexPallet, model.PerPallet)
		}
		if len(cis) == 0 {
			// больше нет км
			// выходим без расчета номера палеты и последняя палета так и останется последней сгенерированной
			break
		}
		lastSSCC = startSSCC + indexPallet
		pallet, err := k.GenerateSSCC(lastSSCC, model.PrefixSSCC)
		if err != nil {
			return fmt.Errorf("generate sscc error %w", err)
		}
		// ищем такую палетту в БД
		plt, err := dbZnak.FindPallet(pallet)
		if !errors.Is(err, db.ErrNoMoreRows) {
			return fmt.Errorf("поиск палеты %s ошибка БД", pallet)
		}
		if plt != nil {
			return fmt.Errorf("паллета %s уже использована в БД", pallet)
		}

		if _, ok := k.Pallet[pallet]; ok {
			return fmt.Errorf("паллета %s уже сгенерирована прежде в обработке", pallet)
		}
		k.Pallet[pallet] = cis
		if model.Entirely {
			break
		}
		if len(cis) < model.PerPallet {
			// последняя не полная палетта
			break
		}
		indexPallet++
	}
	model.LastSSCC = lastSSCC
	_ = reductor.Instance().SetModel("", model)
	return nil
}

// получить следующие км
// i номер группы по count штук
// если размер массива меньше count значит последний
// елси размер массива 0 значит больше нет
func (k *Krinica) nextRecords(i int, count int) (out []*utility.CisInfo) {
	lenCis := len(k.Cis)
	out = make([]*utility.CisInfo, 0)
	first := i * count // первая км в цикле 24 шт
	for i := 0; i < count; i++ {
		index := i + first
		if (index + 1) > lenCis {
			return out
		}
		out = append(out, k.Cis[index])
	}
	return out
}

// index from 0 startIndex 0
// nextRecords returns a batch of records starting from startIndex
// Returns empty slice when no more records are available
func nextRecords(arr []*utility.CisInfo, index int, count int) []*utility.CisInfo {
	startIndex := index * count
	if startIndex >= len(arr) {
		return []*utility.CisInfo{}
	}
	endIndex := startIndex + count
	// если последний индекс больше длины массива укорачиваем до размера массива
	if endIndex > len(arr) {
		endIndex = len(arr)
	}
	return arr[startIndex:endIndex]
}
