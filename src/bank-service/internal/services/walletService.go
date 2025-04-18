package services

import (
	"bankService/internal/commands/commandRequest"
	"bankService/internal/commands/commandResponse"
	"bankService/internal/logger"
	"bankService/internal/models"
	"bankService/internal/storage"
	"bankService/internal/validators"
	"sync"

	"github.com/gin-gonic/gin"
)

var mutex sync.Mutex

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

	mutex.Lock()
	err = s.db.GetTx().Exec("CALL ballanseUpdate(?, ?, ?)", req.Id, req.Amount, req.Operation).Error
	mutex.Unlock()
	if err != nil {
		return err
	}

	return nil
}

func (s *WalletService) Get(id *string) (*commandResponse.WalletResponse, error) {
	if s.autoCloseStore {
		defer s.Close()
	}
	var err error
	defer s.setRollBack(&err)

	err = validators.NewWalletValidators(s.db.GetTx()).ValidateGet(id)
	if err != nil {
		return nil, err
	}

	var entity models.Wallet
	err = s.db.GetTx().Raw("SELECT * FROM wallets w WHERE Id = ? limit 1", id).Scan(&entity).Error
	if err != nil {
		return nil, err
	}

	result := &commandResponse.WalletResponse{
		Result: commandResponse.WalletItem{
			Id:       entity.Id,
			Ballance: entity.Ballance,
		},
	}

	return result, nil
}
