package handler

import (
	"bankService/internal/handler/entityHandlers"
	"bankService/internal/helpers/middlewhere"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func RegisterRoutes(r *gin.Engine) {
	r.Use(middlewhere.LoggingMiddleware())
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		walletHandler := entityHandlers.NewWalletHandler()
		wallet := v1.Group("/wallet")
		{
			wallet.POST("", walletHandler.Create)
		}

		wallets := v1.Group("/wallets")
		{
			wallets.GET("/:id", walletHandler.Get)
		}
	}
}
