package main

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"io"
	"io/ioutil"
	"iris/project/config"
	"iris/project/controller"
	"iris/project/datasource"
	"iris/project/model"
	"iris/project/service"
	"iris/project/util"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	app := newApp()

	//应用App设置
	configation(app)

	////路由设置
	mvcHandle(app)
	//
	conf := config.InitConfig()
	app.Run(
		iris.Addr(":"+conf.Port),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}

func newApp() *iris.Application {
	app := iris.New()

	//设定应用图标
	app.Favicon("./static/favicons/favicon.ico")

	//设置日志级别
	app.Logger().SetLevel("debug")

	//注册静态资源
	app.HandleDir("/static", "./static")
	app.HandleDir("/manage/static", "./static")
	app.HandleDir("/img", "./static/img")

	//注册视图文件
	app.RegisterView(iris.HTML("./static", ".html"))
	app.Get("/", func(c context.Context) {
		c.View("index.html")
	})

	return app
}

/**
 * @Author x
 * @Description //跨域
 * @Date 下午9:30 2020/9/27
 * @Param
 * @return
 **/
func Cors(app *iris.Application) {
	cros := func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Method() == iris.MethodOptions {
			ctx.Header("Access-Control-Methods", "POST, PUT, PATCH, DELETE")
			ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Content-Type")
			ctx.Header("Access-Control-Max-Age", "86400")
			ctx.StatusCode(iris.StatusNoContent)
			return
		}
		ctx.Next()
	}

	app.UseGlobal(cros)

}

/**
* MVC 架构设置
 */
func mvcHandle(app *iris.Application) {
	//启用session
	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookie",
		Expires: 24 * time.Hour,
	})
	//redis
	redis := datasource.NewRedis()
	////把session 的同步位置为redis
	sessManager.UseDatabase(redis)

	//mysal 引擎实例化
	engine := datasource.NewMySqlEngine()

	//管理员模块功能
	adminService := service.NewAdminService(engine)
	admin := mvc.New(app.Party("/admin"))
	admin.Register( //注册后会放到依赖里面
		adminService,
		sessManager.Start,
	)
	admin.Handle(new(controller.AdminController))

	//管理员模块功能
	statisService := service.NewStatisService(engine)
	statis := mvc.New(app.Party("/statis/{model}/{date}/"))
	statis.Register(
		statisService,
		sessManager.Start,
	)
	statis.Handle(new(controller.StatisController))

	//用户路由
	useService := service.NewUserService(engine)
	user := mvc.New(app.Party("/v1/users"))
	user.Register(
		useService,
		sessManager.Start,
	)
	user.Handle(new(controller.UserController))

	//订单模块
	orderService := service.NewOrderService(engine)
	order := mvc.New(app.Party("/bos/orders/"))
	order.Register(
		orderService,
		sessManager.Start,
	)
	order.Handle(new(controller.OrderController)) //控制器

	//商铺模块
	shopService := service.NewShopServie(engine)
	shop := mvc.New(app.Party("/shopping/restaurants/"))
	shop.Register(
		shopService,
		sessManager.Start,
	)
	shop.Handle(new(controller.ShopController))

	//食品类别
	categoryService := service.NewCategoryService(engine)
	category := mvc.New(app.Party("/shopping/"))
	category.Register(
		categoryService,
	)
	category.Handle(new(controller.CategoryController))

	//食品

	//图片上传
	app.Post("/v1/addimg/{model}", func(c context.Context) {
		model := c.Params().Get("model")
		file, info, err := c.FormFile("file")
		if err != nil {
			iris.New().Logger().Error(err.Error())
			c.JSON(iris.Map{
				"status":  util.RECODE_FAIL,
				"type":    util.RESPMSG_ERROR_PICTUREADD,
				"failure": util.Recode2Text(util.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}
		defer file.Close()
		fname := info.Filename
		isExist, err := util.PathExists("./uploads/" + model)
		if !isExist {
			err := os.Mkdir("./uploads/"+model, 0777)
			if err != nil {
				c.JSON(iris.Map{
					"satus":   util.RECODE_FAIL,
					"type":    util.RESPMSG_ERROR_PICTUREADD,
					"failure": util.Recode2Text(util.RESPMSG_ERROR_PICTUREADD),
				})
				return
			}
		}
		out, err := os.OpenFile("./uploads/"+model+"/"+fname, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			iris.New().Logger().Error(err.Error())
			c.JSON(iris.Map{
				"status":  util.RECODE_FAIL,
				"type":    util.RESPMSG_ERROR_PICTUREADD,
				"failure": util.Recode2Text(util.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			iris.New().Logger().Error(err.Error())
			c.JSON(iris.Map{
				"status":  util.RECODE_FAIL,
				"type":    util.RESPMSG_ERROR_PICTUREADD,
				"failure": util.Recode2Text(util.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}

		//上传成功
		c.JSON(iris.Map{
			"status":     util.RECODE_OK,
			"image_path": fname,
		})

	})

	//地址poi 检索
	app.Get("/v1/pois?{poi}", func(c context.Context) {
		path := c.Request().URL.String()
		rs, err := http.Get("https://elm.cangdu.org" + path)
		if err != nil {
			c.JSON(iris.Map{
				"satatus": util.RECODE_FAIL,
				"type":    util.RESPMSG_ERROR_SEARCHADDRESS,
				"message": util.Recode2Text(util.RESPMSG_ERROR_SEARCHADDRESS),
			})
			return
		}

		// 请求成功
		body, err := ioutil.ReadAll(rs.Body)
		var searchList []*model.PoiSearch
		json.Unmarshal(body, &searchList)
		c.JSON(&searchList)

	})

	// 文件上传
	app.Post("/admin/update/avatar/{adminId}", func(c context.Context) {
		adminId := c.Params().Get("adminId")
		iris.New().Logger().Info(adminId)
		file, info, err := c.FormFile("file")
		if err != nil {
			iris.New().Logger().Error(err.Error())
			c.JSON(iris.Map{
				"status":  util.RECODE_FAIL,
				"type":    util.RESPMSG_ERROR_PICTUREADD,
				"failure": util.Recode2Text(util.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}
		defer file.Close()
		fname := info.Filename
		out, err := os.OpenFile("./uploads/"+fname, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			iris.New().Logger().Error(err.Error())
			c.JSON(iris.Map{
				"status":  util.RECODE_FAIL,
				"type":    util.RESPMSG_ERROR_PICTUREADD,
				"failure": util.Recode2Text(util.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}
		iris.New().Logger().Info("文件路径：" + out.Name())
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			c.JSON(iris.Map{
				"status":  util.RECODE_FAIL,
				"type":    util.RESPMSG_ERROR_PICTUREADD,
				"failure": util.Recode2Text(util.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}

		intAdminId, _ := strconv.Atoi(adminId)
		adminService.SaveAvatarImg(int64(intAdminId), fname)
		c.JSON(iris.Map{
			"status":     util.RECODE_OK,
			"image_path": file,
		})
	})

}

/**
 * 项目设置
 */
func configation(app *iris.Application) {
	//配置 字符编码
	app.Configure(iris.WithConfiguration(iris.Configuration{
		Charset: "UTF-8",
	}))

	//错误配置
	//未发现错误 404
	app.OnErrorCode(iris.StatusNotFound, func(c context.Context) {
		c.JSON(iris.Map{
			"errmsg": iris.StatusNotFound,
			"msg":    "not found",
			"data":   iris.Map{},
		})
	})

	app.OnErrorCode(iris.StatusInternalServerError, func(c context.Context) {
		c.JSON(iris.Map{
			"errmsg": iris.StatusInternalServerError,
			"msg":    "interal error",
			"data":   iris.Map{},
		})

	})
}
