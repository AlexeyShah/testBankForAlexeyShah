package storage

import (
	"bankService/internal/models"
)

type Migrator struct {
	tx         *PostgreConnector
	isRollback *bool
}

func NewMigrator(tx *PostgreConnector) *Migrator {
	return &Migrator{
		tx: tx,
	}
}

func (m *Migrator) Magrate() error {
	defer m.close()
	var err error
	defer m.setRollBack(&err)

	err = m.tx.GetTx().AutoMigrate(&models.Wallet{})
	if err != nil {
		return err
	}

	return nil
}

func (s *Migrator) close() {
	s.tx.Close(*s.isRollback)
}

func (s *Migrator) setRollBack(err *error) {
	boolVal := false
	if err != nil && *err != nil && len((*err).Error()) > 0 {
		boolVal = true
		s.isRollback = &boolVal
	} else {
		s.isRollback = &boolVal
	}
}
