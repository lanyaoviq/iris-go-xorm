package service

import (
	"github.com/wonderivan/logger"
	"iris/project/model"
	"time"
	"xorm.io/xorm"
)

/**
 * 统计功能模块接口标准
 */
type StatisService interface {

	GetUserDailyCount(date string) int64
	GetOrderDailyCount(date string) int64
	GetAdminDailyCount(date string) int64


}

/**
 * 统计功能service 实现结构体
 */
type statisService struct {
	Engine *xorm.Engine
}


func NewStatisService(db *xorm.Engine)StatisService{
	return &statisService{
		Engine: db,
	}
}

func (ss *statisService) GetUserDailyCount(date string) int64  {
	startDate,err:= time.Parse("2006-01-02",date)
	if err != nil {
		logger.Error("time parse faile")
		return 0
	}
	endDate := startDate.AddDate(0,0,1)
	result,err:=ss.Engine.Where("create_time between ? and ? and del_flag=0",startDate.Format("2006-01-02 15:04:05"),endDate.Format("2006-01-02 15:04:05")).Count(model.User{})
	if err != nil{
		return 0
	}
	return result
}

func (ss *statisService) GetOrderDailyCount(date string) int64 {
	statrDate,err:= time.Parse("2006-01-02",date)
	if err != nil {
		logger.Error("time parse failed")
		return 0
	}

	endDate := statrDate.AddDate(0,0,1)

	count,err:=ss.Engine.Where("create_time between ? and ? and def_flag=0",statrDate.Format("2006-01-02 15:04:05"),endDate.Format("2006-01-02 15:04:05")).Count(model.UserOrder{})
	if err != nil {
		logger.Error("GetOrderDailyCount sql count failed")
		return 0
	}
	return count
}


func (ss *statisService) GetAdminDailyCount(date string) int64 {
	startDate,err:=time.Parse("2006-01-02",date)
	if err != nil {
		logger.Error("time parse failed")
		return 0
	}
	endDate := startDate.AddDate(0,0,1)

	count,err:=ss.Engine.Where("create_time between ? and ? and delflag=0",
		startDate.Format("2006-01-02 15:04:05"),
		endDate.Format("2006-01-02 15:04:05")).Count(model.Admin{})
	if err !=nil{
		logger.Error("GetAdminDailyCount sql count failed")
		return 0
	}
	return count
}