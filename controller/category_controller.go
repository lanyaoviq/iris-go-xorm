package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"iris/project/model"
	"iris/project/service"
	"iris/project/util"
	"strconv"
)

type CategoryController struct {
	Ctx     iris.Context
	Service service.CategoryServie
	Session sessions.Session
}

type CategoryEntity struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	RestaurantId string `json:"restaurant_id"`
}

func (cc *CategoryController) BeforeActivation(a mvc.BeforeActivation) {
	//通过商铺Id获取对应的食品种类
	a.Handle("GET", "/getcategory/{shopId}", "GetCategoryByShopId")

	//获取全部的食品种类
	a.Handle("GET", "/v2/restaurant/category", "GetAllCategory")

	//添加商铺记录
	a.Handle("POST", "/addShop", "PostAddShop")

	//删除商铺记录
	a.Handle("DELETE", "/restaurant/{restaurant_id}", "DeleteRestaurant")

	//删除食品记录
	a.Handle("DELETE", "/v2/food/{food_id}", "DeleteFood")

	//获取某个商铺的信息
	a.Handle("GET", "/restaurant/{restaurant_id}", "GetRestaurantInfo")
}

/**
 * url：/shopping/v2/restaurant/category
 * type：get
 * desc：获取所有食品种类供添加商铺时进行添加
 */
func (cc *CategoryController) GetAllCategory() mvc.Result {
	categorys, err := cc.Service.GetAllCategory()
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  util.RECODE_FAIL,
				"message": util.RESPMSG_ERROR_CATEGORIES,
			},
		}
	}
	return mvc.Response{
		Object: &categorys,
	}
}

/**
 * 删除商户记录
 *
 */
func (cc *CategoryController) DeleteRestaurant() mvc.Result {
	restaurant_id := cc.Ctx.Params().Get("restaurant_id")
	shopId, err := strconv.Atoi(restaurant_id)
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  util.RECODE_FAIL,
				"type":    util.RESPMSG_HASNOACCESS,
				"message": util.Recode2Text(util.RESPMSG_HASNOACCESS),
			},
		}
	}
	delete := cc.Service.DeleteShop(int64(shopId))
	if !delete {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  util.RECODE_FAIL,
				"type":    util.RESPMSG_HASNOACCESS,
				"message": util.Recode2Text(util.RESPMSG_HASNOACCESS),
			},
		}
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"status":  util.RECODE_OK,
			"type":    util.RESPMSG_SUCCESS_DELETESHOP,
			"message": util.Recode2Text(util.RESPMSG_SUCCESS_DELETESHOP),
		},
	}

}

/**
 * url：/shopping/getcategory/1
 * type：get
 * desc：根据商铺Id获取对应的商铺的食品种类列表信息
 */
func (cc *CategoryController) GetCategoryByShopId() mvc.Result {
	shopIdStr := cc.Ctx.Params().Get("shopId")
	if shopIdStr == "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  util.RECODE_FAIL,
				"message": util.Recode2Text(util.RESPMSG_ERROR_CATEGORIES),
			},
		}
	}
	shopId, err := strconv.Atoi(shopIdStr)
	if err != nil {
		iris.New().Logger().Error(err.Error())
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  util.RECODE_FAIL,
				"message": util.Recode2Text(util.RESPMSG_ERROR_CATEGORIES),
			},
		}
	}

	categorys, err := cc.Service.GetCategoryByShopId(int64(shopId))
	if err != nil {
		iris.New().Logger().Error(err.Error())
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  util.RECODE_FAIL,
				"message": util.Recode2Text(util.RESPMSG_ERROR_CATEGORIES),
			},
		}
	}

	return mvc.Response{
		Object: map[string]interface{}{
			"status":        util.RECODE_OK,
			"category_list": &categorys,
		},
	}

}

/**
 * url：/shopping/addcategory
 * type：post
 * desc：添加食品种类记录
 */
func (cc *CategoryController) PostAddcategory() mvc.Result {
	var categoryEntity CategoryEntity
	cc.Ctx.ReadJSON(&categoryEntity)

	if categoryEntity.Name == "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  util.RECODE_FAIL,
				"message": util.Recode2Text(util.RESPMSG_ERROR_CATEGORYADD),
			},
		}
	}
	restaurant_id, _ := strconv.Atoi(categoryEntity.RestaurantId)
	foodCategory := &model.FoodCategory{
		CategoryName:     categoryEntity.Name,
		CategoryDesc:     categoryEntity.Description,
		RestaurantId:     int64(restaurant_id),
		Level:            1,
		ParentCategoryId: 0,
	}
	saveSucc := cc.Service.AddCategory(foodCategory)
	if !saveSucc {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  util.RECODE_FAIL,
				"message": util.Recode2Text(util.RESPMSG_ERROR_CATEGORYADD),
			},
		}
	}

	return mvc.Response{
		Object: map[string]interface{}{
			"status":  util.RECODE_OK,
			"message": util.Recode2Text(util.RESPMSG_SUCCESS_CATEGORYADD),
		},
	}
}

/**
 * 添加商铺方法
 * url：/shopping/addShop
 * type：Post
 * desc：添加商铺数据记录
 */
func (cc *CategoryController) PostAddShop() mvc.Result {
	shop := new(model.Shop)
	err := cc.Ctx.ReadJSON(&shop)
	if err != nil {
		cc.Ctx.Request()
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  util.RECODE_FAIL,
				"message": util.Recode2Text(util.RESPMSG_FAIL_ADDREST),
			},
		}
	}
	res := cc.Service.SaveShop(*shop)
	if !res {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  util.RECODE_FAIL,
				"message": util.Recode2Text(util.RESPMSG_FAIL_ADDREST),
			},
		}
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"status":     util.RECODE_OK,
			"message":    util.Recode2Text(util.RESPMSG_FAIL_ADDREST),
			"shopDetail": shop,
		},
	}

}
