package storage

import (
	"gorm.io/gorm"
)

func NewTestConnection(d gorm.Dialector) (*PostgreConnector, error) {

	db, err := gorm.Open(d, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &PostgreConnector{
		tx: db,
	}, nil
}
