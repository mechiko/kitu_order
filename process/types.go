package process

import (
	"fmt"
	"kitu/domain"
	"kitu/repo/configdb"

	"github.com/mechiko/dbscan"
	"github.com/mechiko/utility"
)

// const startSSCC = "1462709225" // gs1 rus id zapivkom для памяти запивком

type Krinica struct {
	domain.Apper
	repo     domain.Repo
	inn      string
	Sscc     []string
	Cis      []*utility.CisInfo
	Pallet   map[string][]*utility.CisInfo
	warnings []string
	errors   []string
}

func New(app domain.Apper, repo domain.Repo) (*Krinica, error) {
	db, err := repo.Lock(dbscan.Config)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer repo.Unlock(db)
	dbCfg, ok := db.(*configdb.DbConfig)
	if !ok {
		return nil, fmt.Errorf("db type wrong %T", db)
	}
	inn, err := dbCfg.Key("inn")
	if err != nil {
		return nil, fmt.Errorf("find inn error %w", err)
	}
	k := &Krinica{
		Apper:    app,
		repo:     repo,
		inn:      inn,
		Pallet:   make(map[string][]*utility.CisInfo),
		warnings: make([]string, 0),
		errors:   make([]string, 0),
		Sscc:     make([]string, 0),
		Cis:      make([]*utility.CisInfo, 0),
	}
	return k, nil
}

func (k *Krinica) AddWarn(warn string) {
	k.warnings = append(k.warnings, warn)
}

func (k *Krinica) Warnings() []string {
	return k.warnings
}

func (k *Krinica) AddError(err string) {
	k.errors = append(k.errors, err)
}

func (k *Krinica) Errors() []string {
	return k.errors
}

func (k *Krinica) ResetPalletMap() {
	for key := range k.Pallet {
		delete(k.Pallet, key)
	}
}

func (k *Krinica) Reset() {
	k.ResetPalletMap()
	k.Cis = make([]*utility.CisInfo, 0)
	k.Sscc = make([]string, 0)
	k.errors = make([]string, 0)
	k.warnings = make([]string, 0)
}
