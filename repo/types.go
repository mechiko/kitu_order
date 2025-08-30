package repo

import (
	"errors"
	"fmt"
	"sync"

	"github.com/mechiko/dbscan"
	"go.uber.org/zap"
)

const modError = "pkg:repo"

var Version int64

type singleMutex struct {
	mutex sync.Mutex
}

type Repository struct {
	logger  *zap.SugaredLogger
	dbs     *dbscan.Dbs
	dbMutex map[dbscan.DbInfoType]*singleMutex
	listDbs []dbscan.DbInfoType
}

// dbPath для своей БД
// func New(logcfg ILogCfg, dbPath string) (rp *Repository, err error) {
func New(logger *zap.SugaredLogger, listDbs dbscan.ListDbInfoForScan, dbPath string) (rp *Repository, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("repo panic %v", r)
		}
	}()
	if len(listDbs) == 0 {
		return nil, fmt.Errorf("список описателей бд пуст")
	}
	for tp, info := range listDbs {
		if info == nil {
			return nil, fmt.Errorf("%s in list [%v] is nil", modError, tp)
		}
	}
	dbs, err := dbscan.New(listDbs, dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s dbscan error %w", modError, err)
	}
	rp = &Repository{
		logger:  logger,
		dbs:     dbs,
		dbMutex: make(map[dbscan.DbInfoType]*singleMutex),
		listDbs: make([]dbscan.DbInfoType, 0),
	}
	for tp := range listDbs {
		dbInfo := rp.dbs.Info(tp)
		if dbInfo == nil {
			// такая ошибка не вероятна дбскан выдаст ошибку при сканировании
			// но проверить надо вдруг чего...
			err = errors.Join(err, fmt.Errorf("%s отсутствует БД %v", modError, tp))
		} else {
			rp.listDbs = append(rp.listDbs, tp)
			// создаем в мапе мьютекс
			if _, ok := rp.dbMutex[tp]; !ok {
				rp.dbMutex[tp] = &singleMutex{}
			}
		}
	}
	if err != nil {
		return nil, fmt.Errorf("%s не все бд найдены %v", modError, err)
	}
	if di := rp.dbs.Info(dbscan.Other); di != nil {
		// инициализация для Self если она есть в настройках списка доступных БД
		if err := rp.prepareSelf(); err != nil {
			return nil, fmt.Errorf("%s ошибка миграции self %w", modError, err)
		}
	}
	return rp, nil
}

// возвращаем DbInfo или nil
func (r *Repository) Info(t dbscan.DbInfoType) *dbscan.DbInfo {
	if di := r.dbs.Info(t); di != nil {
		return di
	}
	return nil
}

func (r *Repository) ListDbs() (out []dbscan.DbInfoType) {
	if r.listDbs == nil {
		return nil
	}
	out = make([]dbscan.DbInfoType, len(r.listDbs))
	copy(out, r.listDbs)
	return out
}
