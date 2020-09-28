package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"iris/project/service"
	"iris/project/util"
	"strconv"
)

//每一页最大的内容
const  MaxLimit = 50

/**
 * 用户控制器结构体：
 */
type UserController struct {
	Ctx iris.Context
	Service service.UserService
	Session sessions.Session
}

/**
 * 获取用户总数
 * type : get
 * url: "/v1/users/count
 */
func (uc *UserController) GetCount() mvc.Result  {

	total,err:=uc.Service.GetUserTotalCount()
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":util.RECODE_FAIL,
				"count":0,
			},
		}
	}
	// 正常情况的返回值

	return mvc.Response{
		Object: map[string]interface{}{
			"status":util.RECODE_OK,
			"count":total,
		},
	}

}

/**
 * 获取用户列表
 * type:get
 * url:/v1/users/list
 */
func (uc *UserController) GetList() mvc.Result  {
	offsetStr := uc.Ctx.FormValue("offset")
	limitStr := uc.Ctx.FormValue("limit")

	var offset int
	var limit int

	if offsetStr == "" || limitStr == "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":util.RECODE_FAIL,
				"type":util.RESPMSG_ERROR_USERLIST,
				"message":util.Recode2Text(util.RESPMSG_ERROR_USERLIST),
			},
		}
	}

	offset,err:= strconv.Atoi(offsetStr)
	limit,err= strconv.Atoi(limitStr)
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":util.RECODE_FAIL,
				"type":util.RESPMSG_ERROR_USERLIST,
				"message":util.Recode2Text(util.RESPMSG_ERROR_USERLIST),
			},
		}
	}

	//检测 offset 还有limit
	if offset <=0 {
		offset=0
	}
	if limit>MaxLimit{
		limit=MaxLimit
	}

	userList:=uc.Service.GetUserList(offset,limit)

	if len(userList)  == 0 {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":util.RECODE_FAIL,
				"type":util.RESPMSG_ERROR_USERLIST,
				"message":util.Recode2Text(util.RESPMSG_ERROR_USERLIST),
			},
		}
	}

	//将查询到的数据转换成前段需要的类容
	var respList []interface{}
	for _,user:=range userList{
		respList =append(respList,user.UserToRespDesc())
	}

	return mvc.Response{
		Object: &respList,
	}
}


