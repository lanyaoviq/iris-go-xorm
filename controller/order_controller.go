package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"iris/project/service"
	"iris/project/util"
)

type OrderController struct {
	Ctx iris.Context
	Service service.OrderService
	Session *sessions.Session
}

/**
 * 查询订单记录总数
 */
func (oc *OrderController) GetCount() mvc.Result  {
	iris.New().Logger().Info("查询订单记录总数")
	count,err:=oc.Service.GetCount()
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":util.RECODE_FAIL,
				"count":0,
			},
		}
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"status":util.RECODE_OK,
			"count":count,
		},
	}

}