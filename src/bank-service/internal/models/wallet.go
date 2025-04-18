package models

type Wallet struct {
	Id       *string `gorm:"not null;column:id;type:text;primaryKey"`
	Ballance *int64  `gorm:"not null;column:ballance;type:bigint"`
}

func (*Wallet) TableName() string {
	return "wallets"
}
