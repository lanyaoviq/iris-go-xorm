package service

import (
	"github.com/kataras/iris/v12"
	"iris/project/model"
	"time"
	"xorm.io/xorm"
)

/**
 * 用户模块service 接口
 */
type UserService interface {
	//用户日增长统计数据
	GetUserDailyStatisCount(datetime string) int64

	//获取用户总数
	GetUserTotalCount() (int64,error)

	//用户列表
	GetUserList(offset, limit int) []*model.User
}

/**
 * userService 实现 struct
 */
type userService struct {
	Engine *xorm.Engine
}

/**
 * 工厂模式，实例化userService Struct
 */
func NewUserService(db *xorm.Engine) UserService {
	return &userService{
		Engine: db,
	}
}

func (us *userService) GetUserDailyStatisCount(datetime string) int64 {
	date, err := time.Parse("2006-01-02 15:04:05", datetime)
	if err != nil {
		return 0
	}
	count, err := us.Engine.Where("create_time = ?", date).Count(model.User{})
	if err != nil {
		iris.New().Logger().Error(err.Error())
		panic(err.Error())
		return 0
	}
	return count
}

func (us *userService) GetUserTotalCount() (int64, error) {

	count, err := us.Engine.Where("del_flag=0").Count(model.User{})
	if err != nil {
		iris.New().Logger().Error(err.Error())
		panic(err.Error())
		return 0, err
	}
	return count, nil
}

func (us *userService) GetUserList(offset, limit int) []*model.User {
	var users []*model.User
	err:= us.Engine.Where("del_flag=?",0).Limit(limit,offset).Find(&users)
	if err != nil {
		iris.New().Logger().Error(err.Error())
		panic(err.Error())
		return nil
	}
	return users
}
