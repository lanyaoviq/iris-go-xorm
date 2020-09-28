package service

import (
	"github.com/wonderivan/logger"
	"iris/project/model"
	"xorm.io/xorm"
)

/**
 * shop service
 */
type ShopService interface {
	GetShopList(offset,limit int) []*model.Shop

	GetShopCount() int64

}

/**
 * shop 实现结构体
 */
type shopService struct{
	Engine *xorm.Engine
}


func NewShopServie(db *xorm.Engine) ShopService  {
	return &shopService{
		Engine: db,
	}
}


func (ss *shopService) GetShopCount() int64{
	count,err:=ss.Engine.Where("dele =0 ").Count(model.Shop{})
	if err != nil {
		logger.Error("GetShopCount sql count failed")
		return 0
	}
	return count
}

func (ss *shopService) GetShopList(offset,limit int) []*model.Shop {
	var shops []*model.Shop
	err:=ss.Engine.Where("dele=0").Limit(limit,offset).Find(&shops)
	if err != nil {
		logger.Error("GetShopList sql failed")
		return nil
	}
	return shops
}