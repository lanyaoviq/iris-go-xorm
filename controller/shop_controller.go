package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"iris/project/service"
	"iris/project/util"
	"strconv"
)

type ShopController struct {
	Ctx     iris.Context
	Service service.ShopService
	Session *sessions.Session
}

/**
 * 获取商铺列表
 */
func (sc *ShopController) Get() mvc.Result {
	offsetStr := sc.Ctx.FormValue("offset")
	limitStr := sc.Ctx.FormValue("limit")

	if offsetStr == "" || limitStr == "" {
		offsetStr = "0"
		limitStr = "20"
	}

	offset, err := strconv.Atoi(offsetStr)
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		offset = 0
		limit = 20
	}
	shopList := sc.Service.GetShopList(offset,limit)
	if len(shopList) <= 0 {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":util.RECODE_FAIL,
				"type":util.RESPMSG_ERROR_RESTLIST,
				"message":util.Recode2Text(util.RESPMSG_ERROR_RESTLIST),
			},
		}
	}
	var respList []interface{}
	for _, shop := range shopList {
		respList =append(respList,shop.ShopToRespDesc())
	}

	return mvc.Response{
		Object: respList,
	}
}

func (sc *ShopController) GetCount() mvc.Result  {
	count :=sc.Service.GetShopCount()
	return mvc.Response{
		Object: count,
	}
}
