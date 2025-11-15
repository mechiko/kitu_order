package znakdb

import (
	"fmt"
	"kitu/reductor"
	"time"

	"github.com/mechiko/utility"
	"github.com/upper/db/v4"
)

type Aggregation struct {
	Id                      int64  `db:"id,omitempty"`
	CreateDate              string `db:"create_date"`
	Inn                     string `db:"inn"`
	UnitSerialNumber        string `db:"unit_serial_number"`
	AggregationUnitCapacity string `db:"aggregation_unit_capacity"`
	AggregatedItemsCount    string `db:"aggregated_items_count"`
	AggregationType         string `db:"aggregation_type"`
	Gtin                    string `db:"gtin"`
	Note                    string `db:"note"`
	Version                 string `db:"version"`
	State                   string `db:"state"`
	Status                  string `db:"status"`
	OrderId                 string `db:"order_id"`
	Archive                 int    `db:"archive"`
	Json                    string `db:"json"`
}

type AggregationCode struct {
	Id                      int64  `db:"id,omitempty"`
	IdOrderMarkAggregation  int    `db:"id_order_mark_aggregation"`
	SerialNumber            string `db:"serial_number"`
	Code                    string `db:"code"`
	UnitSerialNumber        string `db:"unit_serial_number"`
	AggregationUnitCapacity string `db:"aggregation_unit_capacity"`
	AggregatedItemsCount    string `db:"aggregated_items_count"`
	Status                  string `db:"status"`
}

func (z *DbZnak) WritePallets(pallets map[string][]*utility.CisInfo, model reductor.Model) (err error) {
	opts := map[string]string{
		"inn":   model.Inn,
		"order": fmt.Sprintf("%d", model.Order),
		"count": fmt.Sprintf("%d", model.PerPallet),
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	sess := z.Sess()
	defer func() {
		if err != nil {
			if errClose := sess.Close(); errClose != nil {
				err = fmt.Errorf("%w%w", errClose, err)
			}
		} else {
			err = sess.Close()
		}
	}()
	err = sess.Tx(func(tx db.Session) error {
		agr := &Aggregation{
			CreateDate:      time.Now().Local().Format("2006.01.02 15:04:05"),
			Inn:             model.Inn,
			OrderId:         fmt.Sprintf("%d", model.Order),
			Note:            "",
			Gtin:            model.Gtin,
			AggregationType: "Упаковка",
			Version:         "1",
			State:           "Создан",
			Status:          "Не проведён",
			Archive:         0,
		}
		if err := tx.Collection("order_mark_aggregation").InsertReturning(agr); err != nil {
			return err
		} else {
			for _, paletProducer := range pallets {
				cis := k.Palet[paletProducer]
				sscc := k.Sscc[paletProducer]
				if err := k.writePaletForce(tx, agr, sscc, cis); err != nil {
					return err
				}
			}
			return nil
		}
	})
	return nil
}

func (z *DbZnak) writePaletForce(tx db.Session, agr *Aggregation, palet string, cis []*utility.CisInfo) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic %v", r)
		}
	}()
	for i := range cis {
		cis := &AggregationCode{
			IdOrderMarkAggregation:  int(agr.Id),
			SerialNumber:            cis[i].Serial,
			Code:                    cis[i].Cis,
			Status:                  "Импортирован",
			UnitSerialNumber:        palet,
			AggregationUnitCapacity: fmt.Sprintf("%d", len(cis)),
			AggregatedItemsCount:    fmt.Sprintf("%d", len(cis)),
		}
		if _, err := tx.Collection("order_mark_aggregation_codes").Insert(cis); err != nil {
			return err
		}
	}
	return nil
}
