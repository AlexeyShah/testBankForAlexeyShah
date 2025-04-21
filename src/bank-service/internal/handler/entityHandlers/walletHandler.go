package entityHandlers

import (
	"bankService/internal/commands/commandRequest"
	"bankService/internal/commands/commandResponse"
	"bankService/internal/services"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type WalletHandler struct {
	mu sync.Mutex
}

func NewWalletHandler() *WalletHandler {
	return &WalletHandler{}
}

// @Summary create
// @Tags Wallet
// @Accept  json
// @Produce  json
// @Param WalletCreateRequest body commandRequest.WalletCreateRequest false "WalletCreateRequest"
// @Success 201 {object} nil
// @Failure 400 {object} commandResponse.ErrorReponse "Bad Request"
// @Router /wallet [post]
func (h *WalletHandler) Create(c *gin.Context) {
	logg := c.Keys["logg"].(*logrus.Entry)
	logg.Debug("Start Create")
	var req commandRequest.WalletCreateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logg.Error(err)
		c.JSON(http.StatusBadRequest, &commandResponse.ErrorReponse{
			Error: err.Error(),
		})
		return
	}

	h.mu.Lock()
	logg.Debug("loks")
	defer func() {
		h.mu.Unlock()
		logg.Debug("unloks")
	}()

	err = services.NewWalletService(c, true).Create(req)
	if err != nil {
		logg.Error(err)
		c.JSON(http.StatusBadRequest, &commandResponse.ErrorReponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, nil)
	logg.Debug("Created")
}

// @Summary get
// @Tags Wallet
// @Accept  json
// @Produce  json
// @Param id path string true "Id"
// @Success 200 {object} commandResponse.WalletResponse
// @Failure 400 {object} commandResponse.ErrorReponse "Bad Request"
// @Router /wallets/{id} [get]
func (*WalletHandler) Get(c *gin.Context) {
	logg := c.Keys["logg"].(*logrus.Entry)
	logg.Debug("Start Get")
	id := c.Param("id")
	res, err := services.NewWalletService(c, true).Get(&id)
	if err != nil {
		logg.Error(err)
		c.JSON(http.StatusBadRequest, &commandResponse.ErrorReponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary get all
// @Tags Wallet
// @Accept  json
// @Produce  json
// @Success 200 {object} commandResponse.WalletAllResponse
// @Failure 400 {object} commandResponse.ErrorReponse "Bad Request"
// @Router /wallets [get]
func (*WalletHandler) GetAll(c *gin.Context) {
	logg := c.Keys["logg"].(*logrus.Entry)
	logg.Debug("Start GetAll")
	res, err := services.NewWalletService(c, true).GetAll()
	if err != nil {
		logg.Error(err)
		c.JSON(http.StatusBadRequest, &commandResponse.ErrorReponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
