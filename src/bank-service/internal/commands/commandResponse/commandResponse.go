package commandResponse

type WalletAllResponse struct {
	Result []WalletItem `json:"result,omitempty"`
}

type WalletResponse struct {
	Result WalletItem `json:"result,omitempty"`
}

type WalletItem struct {
	Id       *string `json:"id,required"`
	Ballance *int64  `json:"ballance,required"`
}

type ErrorReponse struct {
	Error string `json:"error,required"`
}
