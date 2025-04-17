package commandRequest

type WalletCreateRequest struct {
	Id        *string `form:"id,omitempty" json:"id,omitempty"`
	Operation *string `form:"operation,omitempty" json:"operation,omitempty"`
	Amount    *int64  `form:"amount,omitempty" json:"amount,omitempty"`
}
