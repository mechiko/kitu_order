package znakdb

import (
	"fmt"
	"kitu/reductor"
	"sort"
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
	Version                 string `db:"version"`
	State                   string `db:"state"`
	Status                  string `db:"status"`
	OrderId                 string `db:"order_id"`
	Archive                 int    `db:"archive"`
	Json                    string `db:"json"`
}

type AggregationCode struct {
	Id                     int64  `db:"id,omitempty"`
	IdOrderMarkAggregation int    `db:"id_order_mark_aggregation"`
	SerialNumber           string `db:"serial_number"`
	Code                   string `db:"code"`
	Status                 string `db:"status"`
}

// opts map[inn] [order] [count]
func (z *DbZnak) WritePallets(pallets map[string][]*utility.CisInfo, model reductor.Model) (err error) {
	opts := map[string]string{
		"inn":   model.Inn,
		"order": fmt.Sprintf("%d", model.Order),
		"count": fmt.Sprintf("%d", model.PerPallet),
	}
	sess := z.Sess()
	err = sess.Tx(func(tx db.Session) error {
		keys2 := make([]string, 0, len(pallets))
		for k := range pallets {
			keys2 = append(keys2, k)
		}
		sort.Strings(keys2)
		for _, pallet := range keys2 {
			cis := pallets[pallet]
			opts["count"] = fmt.Sprintf("%d", len(cis))
			if err := z.writePallet(tx, pallet, cis, opts); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (z *DbZnak) writePallet(tx db.Session, pallet string, cis []*utility.CisInfo, opts map[string]string) (err error) {
	agr := &Aggregation{
		CreateDate:              time.Now().Local().Format("2006.01.02 15:04:05"),
		Inn:                     opts["inn"],
		OrderId:                 opts["order"],
		UnitSerialNumber:        pallet,
		AggregationUnitCapacity: opts["count"],
		AggregatedItemsCount:    opts["count"],
		AggregationType:         "Упаковка",
		Version:                 "1",
		State:                   "Создан",
		Status:                  "Не проведён",
		Archive:                 0,
	}
	if err := tx.Collection("order_mark_aggregation").InsertReturning(agr); err != nil {
		return err
	} else {
		for _, c := range cis {
			code := &AggregationCode{
				IdOrderMarkAggregation: int(agr.Id),
				SerialNumber:           c.Serial,
				Code:                   c.Code,
				Status:                 "Импортирован",
			}
			if _, err := tx.Collection("order_mark_aggregation_codes").Insert(code); err != nil {
				return err
			}
		}
	}
	return nil
}
