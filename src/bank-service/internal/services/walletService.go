package services

import (
	"bankService/internal/commands/commandRequest"
	"bankService/internal/commands/commandResponse"
	"bankService/internal/helpers/consts"
	"bankService/internal/models"
	"bankService/internal/storage"
	"bankService/internal/validators"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type WalletService struct {
	BaseService
	Logg *logrus.Entry
}

func NewWalletService(c *gin.Context, autoCloseStore bool) *WalletService {
	logg := c.Keys["logg"].(*logrus.Entry)
	ctx := c.Request.Context()
	db, err := storage.NewPostgreConnection()
	if err != nil {
		logg.Error(err)
		panic(err)
	}

	return &WalletService{
		BaseService: BaseService{
			AutoCloseStore: autoCloseStore,
			Ctx:            ctx,
			DB:             db,
		},
		Logg: logg,
	}
}

func (s *WalletService) Create(req commandRequest.WalletCreateRequest) error {
	if s.AutoCloseStore {
		defer s.Close()
	}
	var err error
	defer s.setRollBack(&err)

	err = validators.NewWalletValidators(s.DB.GetTx()).ValidateCreate(req)
	if err != nil {
		return err
	}

	var entity *models.Wallet
	err = s.DB.GetTx().Raw("SELECT id, ballance FROM wallets w WHERE id = $1 limit 1", req.Id).Scan(&entity).Error
	if err != nil {
		return err
	}

	has := entity != nil && entity.Id != nil && len(*entity.Id) > 0
	increment := int64(0)
	if has {
		increment = *entity.Ballance
	}

	if *req.Operation == consts.OperationDeposit {
		increment = increment + *req.Amount
	} else if *req.Operation == consts.OperationWithdraw {
		if increment-*req.Amount > 0 {
			increment = increment - *req.Amount
		} else {
			return errors.New("impossible switch to credit")
		}
	}

	if has {
		err = s.DB.GetTx().Exec("UPDATE wallets SET ballance = $2 WHERE id = $1", req.Id, increment).Error
	} else {
		err = s.DB.GetTx().Exec("INSERT INTO wallets (id, ballance) VALUES($1, $2)", req.Id, increment).Error
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *WalletService) Get(id *string) (*commandResponse.WalletResponse, error) {
	if s.AutoCloseStore {
		defer s.Close()
	}
	var err error
	defer s.setRollBack(&err)

	err = validators.NewWalletValidators(s.DB.GetTx()).ValidateGet(id)
	if err != nil {
		return nil, err
	}

	var entity models.Wallet
	err = s.DB.GetTx().Raw("SELECT id, ballance FROM wallets w WHERE id = ? limit 1", id).Scan(&entity).Error
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

func (s *WalletService) GetAll() (*commandResponse.WalletAllResponse, error) {
	if s.AutoCloseStore {
		defer s.Close()
	}
	var err error
	defer s.setRollBack(&err)

	var entity []models.Wallet
	err = s.DB.GetTx().Raw("SELECT * FROM wallets w").Scan(&entity).Error
	if err != nil {
		return nil, err
	}

	result := &commandResponse.WalletAllResponse{}
	for _, val := range entity {
		result.Result = append(result.Result, commandResponse.WalletItem{
			Id:       val.Id,
			Ballance: val.Ballance,
		})
	}
	return result, nil
}
