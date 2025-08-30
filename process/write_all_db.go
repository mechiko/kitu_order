package process

import (
	"errors"
	"fmt"
	"kitu/reductor"
	"kitu/repo/znakdb"

	"github.com/mechiko/dbscan"
)

func (k *Krinica) WritePallets() error {
	db, err := k.repo.Lock(dbscan.TrueZnak)
	if err != nil {
		return fmt.Errorf("lock TrueZnak: %w", err)
	}
	var retErr error
	defer func() {
		if uerr := k.repo.Unlock(db); uerr != nil {
			retErr = errors.Join(retErr, fmt.Errorf("unlock TrueZnak: %w", uerr))
		}
	}()

	dbZnak, ok := db.(*znakdb.DbZnak)
	if !ok {
		return fmt.Errorf("db type wrong %T", db)
	}
	model := reductor.Instance().Model("")
	model.Inn = k.inn
	if err := dbZnak.WritePallets(k.Pallet, model); err != nil {
		return err
	}
	return retErr
}
