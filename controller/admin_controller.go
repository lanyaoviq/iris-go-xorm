package controller

import (
	"encoding/base64"
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/wonderivan/logger"
	"iris/project/model"
	"iris/project/service"
	"iris/project/util"
)

/**
 * 管理员控制器
 */
type AdminController struct {
	//iris自动为每个请求绑定的上下文对象
	Ctx iris.Context

	//admin功能实体
	Service service.AdminService

	//Session 对象 ，事务
	Session *sessions.Session
}

const (
	ADMINTABLENAME = "admin"
	ADMIN          = "admin"
)

type AdminLogin struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

/**
 * 获取用户信息
 */
func (ac *AdminController) GetInfo() mvc.Result {
	//从session 中获取信息
	userStr := ac.Session.Get(ADMIN)
	//session 为空
	if userStr == nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  util.RECODE_UNLOGIN,
				"type":    util.EEROR_UNLOGIN,
				"message": util.Recode2Text(util.EEROR_UNLOGIN),
			},
		}
	}
	//解析数据到admin
	var admin model.Admin
 	var userByte []byte
	// 一言难尽，
	userByte,err := base64.StdEncoding.DecodeString(userStr.(string))
	err = json.Unmarshal(userByte,&admin)
	////解析失败
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  util.RECODE_UNLOGIN,
				"type":    util.EEROR_UNLOGIN,
				"message": util.Recode2Text(util.EEROR_UNLOGIN),
			},
		}
	}

	//解析成功
	return mvc.Response{
		Object: map[string]interface{}{
			"status": util.RECODE_OK,
			"data":   admin.AdminToRespDesc(),
		},
	}

}

/**
 * 管理员登陆功能
 *  接口 admin/login
 */
func (ac *AdminController) PostLogin(c iris.Context) mvc.Result {
	iris.New().Logger().Info("amdin login")

	var adminLogin AdminLogin
	ac.Ctx.ReadJSON(&adminLogin)

	//数据校验
	if adminLogin.UserName == "" || adminLogin.Password == "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  "0",
				"success": "登陆失败",
				"message": "用户名或者密码为空，请重新填写",
			},
		}
	}

	//查库对比
	admin, exist := ac.Service.GetAdminByAmdinNameAndPassword(adminLogin.UserName, adminLogin.Password)

	if !exist {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  "0",
				"success": "登陆失败",
				"message": "用户名或者密码错误，请重新的登陆",
			},
		}
	}
	// 管理员存在
	userByte,_:=json.Marshal(admin)
	logger.Error(userByte, "----userByte login")
	//存入session
	ac.Session.Set(ADMIN, userByte)
	return mvc.Response{
		Object: map[string]interface{}{
			"status":  "1",
			"success": "登陆成功",
			"message": "管理员登陆成功",
		},
	}
}

/**
 * 管理员退出功能
 * 请求类 GEt
 * 请求url :admin/singout
 */
func (a *AdminController) GetSingout() mvc.Result {
	a.Session.Delete(ADMIN)
	return mvc.Response{
		Object: map[string]interface{}{
			"status":  util.RECODE_OK,
			"success": util.Recode2Text(util.RESPMSG_SIGNOUT),
		},
	}
}

func (a *AdminController) GetCount() mvc.Result {
	adminCount := a.Session.Get("adminCount")
	if adminCount != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status": util.RECODE_OK,
				"count":  adminCount,
			},
		}
	} else {
		count, err := a.Service.GetAdminCount()
		if err != nil {
			return mvc.Response{
				Object: map[string]interface{}{
					"status": util.RECODE_FAIL,
					"count":  0,
				},
			}
		}
		a.Session.Set("adminCount", count)
		return mvc.Response{
			Object: map[string]interface{}{
				"status": util.RECODE_OK,
				"count":  count,
			},
		}
	}
}

func (ac *AdminController) GetTest() mvc.Result {
	ac.Session.Set("key", "阿斯顿发发沙发沙发")
	return mvc.Response{
		Object: map[string]interface{}{
			"value": ac.Session.Get("key"),
		},
	}
}
