package entityHandlers

import (
	"bankService/internal/commands/commandRequest"
	"bankService/internal/commands/commandResponse"
	"bankService/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
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
func (*WalletHandler) Create(c *gin.Context) {
	var req commandRequest.WalletCreateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, &commandResponse.ErrorReponse{
			Error: err.Error(),
		})
		return
	}

	err = services.NewWalletService(c, true).Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, &commandResponse.ErrorReponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, nil)
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
	id := c.Param("id")
	res, err := services.NewWalletService(c, true).Get(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &commandResponse.ErrorReponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
