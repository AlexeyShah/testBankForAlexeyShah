package services

import (
	"bankService/internal/commands/commandRequest"
	"bankService/internal/commands/commandResponse"
	"bankService/internal/helpers/consts"
	"bankService/internal/logger"
	"bankService/internal/models"
	"bankService/internal/storage"
	"bankService/internal/validators"

	"github.com/gin-gonic/gin"
)

type WalletService struct {
	BaseService
}

func NewWalletService(c *gin.Context, autoCloseStore bool) *WalletService {
	ctx := c.Request.Context()
	db, err := storage.NewPostgreConnection()
	if err != nil {
		logger.Logger.Error(err)
		panic(err)
	}

	return &WalletService{
		BaseService: BaseService{
			autoCloseStore: autoCloseStore,
			ctx:            ctx,
			db:             db,
		},
	}
}

func (s *WalletService) Create(req commandRequest.WalletCreateRequest) error {
	if s.autoCloseStore {
		defer s.Close()
	}
	var err error
	defer s.setRollBack(&err)

	err = validators.NewWalletValidators(s.db.GetTx()).ValidateCreate(req)
	if err != nil {
		return err
	}

	var has bool
	err = s.db.GetTx().Raw("SELECT 1 as has FROM wallets w WHERE Id = ?", req.Id).Scan(&has).Error
	if err != nil {
		return err
	}

	var opertaionVal int64
	if *req.Operation == consts.OperationDeposit {
		if *req.Amount < 0 {
			opertaionVal = *req.Amount * -1
		} else {
			opertaionVal = *req.Amount
		}
	} else {
		if *req.Amount < 0 {
			opertaionVal = *req.Amount
		} else {
			opertaionVal = *req.Amount * -1
		}
	}

	if has {
		err = s.db.GetTx().Exec("UPDATE w SET ballance = (w.ballance + ?) FROM wallets w WHERE w.Id = ?", opertaionVal, req.Id).Error
	} else {
		err = s.db.GetTx().Exec("INSERT INTO wallets (id, ballance) VALUES (?, ?)", req.Id, opertaionVal).Error
	}

	return err
}

func (s *WalletService) Get(id *string) (*commandResponse.WalletResponse, error) {
	if s.autoCloseStore {
		defer s.Close()
	}

	err := validators.NewWalletValidators(s.db.GetTx()).ValidateGet(id)
	if err != nil {
		return nil, err
	}

	var entities []models.Wallet
	err = s.db.GetTx().Raw("SELECT * FROM wallets w WHERE Id = ?", id).Scan(&entities).Error
	if err != nil {
		return nil, err
	}

	return &commandResponse.WalletResponse{Result: entities}, err
}
