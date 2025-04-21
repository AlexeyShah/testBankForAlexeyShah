package tests

import (
	"bankService/internal/commands/commandRequest"
	"bankService/internal/helpers/consts"
	"bankService/internal/logger"
	"bankService/internal/services"
	"bankService/internal/storage"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
)

func newTestWalletService(c *gin.Context) (*services.WalletService, sqlmock.Sqlmock) {
	ctx := c.Request.Context()
	db, mockDb, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	conn, err := storage.NewTestConnection(dialector)
	if err != nil {
		panic(err)
	}

	return &services.WalletService{
		BaseService: services.BaseService{
			Ctx: ctx,
			DB:  conn,
		},
		Logg: logger.Logger.WithField("env", "for test"),
	}, mockDb
}

func getTestGinContext() *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(nil)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	return ctx
}

func TestAddBalance(t *testing.T) {
	newId := uuid.New().String()
	operation := consts.OperationDeposit
	val := int64(10)
	req := &commandRequest.WalletCreateRequest{
		Id:        &newId,
		Operation: &operation,
		Amount:    &val,
	}

	service, mockDb := newTestWalletService(getTestGinContext())
	mockDb.ExpectQuery(`^SELECT id, ballance FROM wallets`).
		WithArgs(req.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "ballance"}).AddRow(nil, nil))
	mockDb.ExpectExec(`^INSERT INTO wallets`).WithArgs(req.Id, req.Amount).WillReturnResult(sqlmock.NewResult(1, 1))
	mockDb.ExpectCommit()

	err := service.Create(*req)
	assert.NoError(t, err, "Failed created")
}

func TestUpdateBalance(t *testing.T) {
	newId := uuid.New().String()
	operation := consts.OperationDeposit
	val := int64(10)
	req := &commandRequest.WalletCreateRequest{
		Id:        &newId,
		Operation: &operation,
		Amount:    &val,
	}

	service, mockDb := newTestWalletService(getTestGinContext())
	mockDb.ExpectQuery(`^SELECT id, ballance FROM wallets`).
		WithArgs(req.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "ballance"}).AddRow(newId, 0))
	mockDb.ExpectExec(`^UPDATE wallets`).WithArgs(req.Id, req.Amount).WillReturnResult(sqlmock.NewResult(1, 1))
	mockDb.ExpectCommit()

	err := service.Create(*req)
	assert.NoError(t, err, "Failed created")
}

func TestGetBalance(t *testing.T) {
	newId := uuid.New().String()
	val := int64(10)
	req := &commandRequest.WalletCreateRequest{
		Id:     &newId,
		Amount: &val,
	}

	service, mockDb := newTestWalletService(getTestGinContext())
	mockDb.ExpectQuery(`^SELECT id, ballance FROM wallets`).
		WithArgs(newId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "ballance"}).AddRow(newId, val))
	result, err := service.Get(req.Id)

	assert.NoError(t, err, "Failed get")
	assert.NotNil(t, result, "Failed get result")
	assert.Equal(t, *req.Id, *result.Result.Id, "Failed get id")
	assert.Equal(t, *req.Amount, *result.Result.Ballance, "Failed get balance")

}
