package commandResponse

import "bankService/internal/models"

type WalletResponse struct {
	Result []models.Wallet `form:"result,omitempty" json:"result,omitempty"`
}

type ErrorReponse struct {
	Error string `form:"error,omitempty" json:"error,omitempty"`
}
