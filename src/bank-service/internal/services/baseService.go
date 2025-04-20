package services

import (
	"bankService/internal/storage"
	"context"
)

type BaseService struct {
	Ctx            context.Context
	DB             *storage.PostgreConnector
	AutoCloseStore bool
	IsRollback     *bool
}

func (s *BaseService) Close() {
	s.DB.Close(*s.IsRollback)
}

func (s *BaseService) setRollBack(err *error) {
	boolVal := false
	if err != nil && *err != nil && len((*err).Error()) > 0 {
		boolVal = true
		s.IsRollback = &boolVal
	} else {
		s.IsRollback = &boolVal
	}
}
