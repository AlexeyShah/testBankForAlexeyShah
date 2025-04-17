package storage

import (
	"bankService/internal/helpers/consts"
	"errors"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreConnector struct {
	tx *gorm.DB
}

func NewPostgreConnection() (*PostgreConnector, error) {
	connectionString := os.Getenv(consts.CloudConnectionPostgre)

	if len(connectionString) == 0 {
		return nil, errors.New("connection string invalid")
	}

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &PostgreConnector{
		tx: tx,
	}, nil
}

func (s *PostgreConnector) Close(rollback bool) error {
	if s.tx != nil {
		if rollback {
			err := s.tx.Rollback().Error
			if err != nil {
				return err
			}
		} else {
			err := s.tx.Commit().Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *PostgreConnector) GetTx() *gorm.DB {
	return s.tx
}
