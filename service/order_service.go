package service

import (
	"github.com/kataras/iris/v12"
	"iris/project/model"
	"xorm.io/xorm"
)

type OrderService interface {
	GetCount()(int64,error)
	GetOrderList(offset,limit int) []model.OrderDetail
}

type orderServcie struct {
	Engine *xorm.Engine
}

func NewOrderService(db *xorm.Engine) OrderService{
	return &orderServcie{
		Engine: db,
	}
}

/**
 * 获取订单总数量
 */
func (os *orderServcie) GetCount()(int64,error)  {
	count,err:=os.Engine.Where("del_flag =?",0).Count(new(model.UserOrder))
	if err != nil {
		return 0, err
	}
	return count,nil
}

/**
 * 获取订单列表
 */
func (os *orderServcie) GetOrderList(offset,limit int) []model.OrderDetail  {
	 orderList := make([]model.OrderDetail,0)
	 err:=os.Engine.Table("user_order").
	 	Join("INNER","order_status","order_status.id=user_order.order_status_id").
	 	Join("INNER","user","user.id=user_order.user_id").
	 	Join("INNER","shop","shop.shop_id=user_order.shop_id").
	 	Join("INNER","address","address.address_id=user_order.address_id").
	 	Find(&orderList)
	 if err != nil{
	 	iris.New().Logger().Error(err.Error())
	 	panic(err.Error())
	 	return nil
	 }
	 iris.New().Logger().Info(orderList)
	 return orderList
}

