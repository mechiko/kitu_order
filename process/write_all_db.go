package process

import (
	"fmt"
	"kitu/reductor"
	"kitu/repo/znakdb"

	"github.com/mechiko/dbscan"
)

func (k *Krinica) WritePallets() error {
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
	model.Inn = k.inn
	if err := dbZnak.WritePallets(k.Pallet, model); err != nil {
		return err
	}
	return nil
}
