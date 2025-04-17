package validators

import (
	"bankService/internal/commands/commandRequest"
	"bankService/internal/helpers/consts"
	"errors"

	"gorm.io/gorm"
)

type WalletValidators struct {
	tx *gorm.DB
}

func NewWalletValidators(tx *gorm.DB) *WalletValidators {
	return &WalletValidators{
		tx: tx,
	}
}

func (v *WalletValidators) ValidateCreate(req commandRequest.WalletCreateRequest) error {
	if req.Id == nil {
		return errors.New("Id not set")
	}

	if req.Amount == nil || *req.Amount <= 0 {
		return errors.New("Amount incorrect, amount need > 0")
	}

	if req.Operation == nil || *req.Operation != consts.OperationDeposit && *req.Operation != consts.OperationWithdraw {
		return errors.New("Operation incorrect, need Deposit or Withdraw")
	}

	return nil
}

func (v *WalletValidators) ValidateGet(id *string) error {
	if id == nil {
		return errors.New("Id not set")
	}

	return nil
}
