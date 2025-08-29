package process

import (
	"fmt"
	"kitu/reductor"
	"kitu/repo/znakdb"

	"github.com/mechiko/dbscan"
	"github.com/mechiko/utility"
)

func (k *Krinica) ReadOrder() (err error) {
	db, err := k.repo.Lock(dbscan.TrueZnak)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer k.repo.Unlock(db)

	dbZnak, ok := db.(*znakdb.DbZnak)
	if !ok {
		return fmt.Errorf("db type wrong %T", db)
	}

	model := reductor.Instance().Model("")
	if order, err := dbZnak.FindOrder(model.Order); err != nil {
		return err
	} else {
		model.Name = order.ProductName
		model.Gtin = order.Gtin
	}
	if production, err := dbZnak.FindOrderProductionDate(model.Order); err == nil {
		model.ProductionDate = production
	}
	var arr znakdb.SliceOrderSerialNumbers
	if model.StatusKM {
		if arr, err = dbZnak.OrderSerialNumbersApply(model.Order); err != nil {
			return err
		}
	} else {
		if arr, err = dbZnak.OrderSerialNumbers(model.Order); err != nil {
			return err
		}
	}
	for _, cis := range arr {
		item, err := utility.ParseCisInfo(cis.Code)
		if err != nil {
			return fmt.Errorf("parse cis cis:%s %w", cis.Code, err)
		}
		k.Cis = append(k.Cis, item)
	}
	_ = reductor.Instance().SetModel("", model)
	return nil
}
