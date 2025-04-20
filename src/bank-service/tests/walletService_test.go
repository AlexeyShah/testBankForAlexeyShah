package tests

import (
	"bankService/internal/commands/commandRequest"
	"bankService/internal/helpers/consts"
	"bankService/internal/logger"
	"bankService/internal/services"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	sqliteDriver "modernc.org/sqlite"
)

func newTestWalletService(c *gin.Context) *services.WalletService {
	ctx := c.Request.Context()
	db, err := gorm.Open(sqliteDriver.Open(":memory:"), &gorm.Config{})
	if err != nil {
		logger.Logger.Error(err)
		panic(err)
	}

	return &services.WalletService{
		BaseService: services.BaseService{
			Ctx: ctx,
			DB:  db,
		},
	}
}

func getTestGinContext() *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(nil)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	return ctx
}

func TestBalance(t *testing.T) {
	newId := uuid.New().String()
	operation := consts.OperationDeposit
	val := int64(10)
	req := &commandRequest.WalletCreateRequest{
		Id:        &newId,
		Operation: &operation,
		Amount:    &val,
	}
	service := newTestWalletService(getTestGinContext())
	err := service.Create(*req)

	assert.NoError(t, err, "Failed created")

	result, err := service.Get(req.Id)

	assert.NoError(t, err, "Failed get")
	assert.NotNil(t, result, "Failed get result")
	assert.Equal(t, *req.Id, *result.Result.Id, "Failed get id")
	assert.Equal(t, *req.Amount, *result.Result.Ballance, "Failed get balance")

}
