package services

import (
	"bankService/internal/storage"
	"context"
)

type BaseService struct {
	ctx            context.Context
	db             *storage.PostgreConnector
	autoCloseStore bool
	isRollback     *bool
}

func (s *BaseService) Close() {
	s.db.Close(*s.isRollback)
}

func (s *BaseService) setRollBack(err *error) {
	boolVal := false
	if err != nil && *err != nil && len((*err).Error()) > 0 {
		boolVal = true
		s.isRollback = &boolVal
	} else {
		s.isRollback = &boolVal
	}
}
