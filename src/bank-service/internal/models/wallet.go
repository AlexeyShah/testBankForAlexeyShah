package models

type Wallet struct {
	Id       *string `json:"id,omitempty" gorm:"not null;column:id"`
	Ballance *int64  `json:"ballance,omitempty" gorm:"not null;column:ballance"`
}
