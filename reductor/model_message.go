package reductor

import (
	"fmt"
	"kitu/domain"
)

type Model struct {
	StartNumberSSCC int
	PerPallet       int
	Inn             string
	PrefixSSCC      string
	Entirely        bool
	Name            string
	Gtin            string
	LastSSCC        int
	Date            string
	ProductionDate  string
	Order           int64
	StatusKM        bool
}

type ModelList map[string]Model

// формат сообщения
type Message struct {
	Sender string
	Page   string
	Model  Model
}

func (m *Model) Sync(cfg domain.Apper) error {
	cfg.SaveOptions("entirely", m.Entirely)
	cfg.SaveOptions("statuskm", m.StatusKM)
	// cfg.SaveOptions("inn", m.Inn)
	cfg.SaveOptions("ssccprefix", m.PrefixSSCC)
	cfg.SaveOptions("ssccstartnumber", m.StartNumberSSCC)
	cfg.SaveOptions("perpallet", m.PerPallet)
	return cfg.SaveAllOptions()
}

func (m *Model) Read(cfg domain.Apper) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	m.Entirely = cfg.Options().Entirely
	m.StatusKM = cfg.Options().StatusKM
	// m.Inn = cfg.GetKeyString("inn")
	m.PrefixSSCC = cfg.Options().SsccPrefix
	m.StartNumberSSCC = cfg.Options().SsccStartNumber
	m.PerPallet = cfg.Options().PerPallet
	return nil
}
