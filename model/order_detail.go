package model

type OrderDetail struct {
	UserOrder `xorm:"extends"`
	Address   `xorm:"extends"`
	Shop      `xorm:"extends"`
}
