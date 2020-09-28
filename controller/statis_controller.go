package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"iris/project/service"
	"iris/project/util"
	"strings"
)

type StatisController struct {
	Ctx iris.Context
	//统计功能的服务实现接口
	Serveice service.StatisService
	//session  //查询到的数据缓存到session中
	Session *sessions.Session
}

const (
	USERMODEL  = "user_"
	ORDERMODEL = "order_"
	ADMINMODEL = "admin_"
)

/**
 * 解析统计功能路由请求
 * type:get
 * url: /statis/user/2020-09-25/count
 */
func (sc *StatisController) GetCount() mvc.Result {
	path := sc.Ctx.Path()

	var pathSlice []string
	if path != "" {
		pathSlice = strings.Split(path, "/")
	}
	// 不符合请求格式
	if len(pathSlice) != 5 {
		return mvc.Response{
			Object: map[string]interface{}{
				"status": util.RECODE_FAIL,
				"count":  0,
			},
		}
	}
	//去掉路径最前面一段
	pathSlice = pathSlice[1:]
	model := pathSlice[1]
	date := pathSlice[2]

	var result int64
	switch model {
	case "user":
		//先从缓存中拿数据
		userResult := sc.Session.Get(USERMODEL+date)
		if userResult != nil {
			userResult = userResult.(float64)
			return mvc.Response{
				Object: map[string]interface{}{
					"status":util.RECODE_OK,
					"count":userResult,
				},
			}
		}else{
			iris.New().Logger().Error(date)
			result = sc.Serveice.GetUserDailyCount(date)
			//放到缓存里面
			sc.Session.Set(USERMODEL+date,result)
		}
	case "order":
		orderResult:= sc.Session.Get(ORDERMODEL+date)
		if orderResult != nil {
			 orderResult = orderResult.(float64)
			 return mvc.Response{
			 	Object: map[string]interface{}{
			 		"status":util.RECODE_OK,
			 		"count":orderResult,
				},
			 }
		}else{
			iris.New().Logger().Error(date)
			result = sc.Serveice.GetOrderDailyCount(date)
			sc.Session.Set(ORDERMODEL+date,result)
		}

	case "admin":
		adminResult:=sc.Session.Get(ADMINMODEL+date)
		if adminResult != nil{
			adminResult = adminResult.(float64)
			return mvc.Response{
				Object: map[string]interface{}{
					"status":util.RECODE_OK,
					"count":adminResult,
				},
			}
		}else{
			iris.New().Logger().Error(date)
			result = sc.Serveice.GetAdminDailyCount(date)
			sc.Session.Set(ADMINMODEL+date,result)
		}
	}

	return mvc.Response{
		Object: map[string]interface{}{
			"status": util.RECODE_OK,
			"count":  result,
		},
	}

}
