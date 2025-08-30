package repo

import (
	"errors"
	"fmt"
	"kitu/domain"
	"kitu/repo/configdb"
	"kitu/repo/znakdb"

	"github.com/mechiko/dbscan"
)

// if err is nil then must after Lock launch UnLock
// всегда или открывает базу и проверяет объект или возвращает ошибку
func (r *Repository) Lock(t dbscan.DbInfoType) (domain.RepoDB, error) {
	r.logger.Infof("repo Lock %v", t)
	info := r.dbs.Info(t)
	if info == nil {
		return nil, fmt.Errorf("repo lock dbinfo is nil for %v", t)
	}
	mu, ok := r.dbMutex[t]
	if ok {
		mu.mutex.Lock()
	} else {
		return nil, fmt.Errorf("repo lock not present mutex %v", t)
	}
	switch t {
	case dbscan.Config:
		db, err := configdb.New(info)
		if err != nil {
			mu.mutex.Unlock()
			return nil, fmt.Errorf("repo lock open %v error %w", dbscan.Config, err)
		}
		return db, nil
	case dbscan.TrueZnak:
		db, err := znakdb.New(info)
		if err != nil {
			mu.mutex.Unlock()
			return nil, fmt.Errorf("repo lock open %v error %w", dbscan.TrueZnak, err)
		}
		return db, nil
	default:
		mu.mutex.Unlock()
		return nil, fmt.Errorf("repo lock not present type mutex %v", t)
	}
}

func (r *Repository) Unlock(db domain.RepoDB) error {
	if db == nil {
		return fmt.Errorf("repo: unlock error db for unlock is nil")
	}
	errClose := db.Close()
	r.logger.Infof("repo UnLock %v", db.InfoType())
	mu, ok := r.dbMutex[db.InfoType()]
	if ok {
		mu.mutex.Unlock()
	} else {
		errUnlock := fmt.Errorf("%s unlock not present mutex %v", modError, dbscan.Other)
		return errors.Join(errClose, errUnlock)
	}
	return errClose
}
