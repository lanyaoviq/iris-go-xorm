package model

import (
	"iris/project/util"
	"time"
)

type User struct {
	Id           int64     `xorm:"pk autoincr" json:"id"`
	UserName     string    `xorm:"varchar(12)" json:"user_name"`
	RegisterTime time.Time `json:"register_time" json:"register_time"`
	Mobile       string    `xorm:"varchar(11)" json:"mobile"`
	IsActive     int64     `json:"is_active"`
	Balance      int64     `json:"balance"`
	Avatar       string    `xorm:"varchar(255)" json:"avatar"`
	Pwd          string    `json:"password""`
	DelFlag      int64     `json:"del_flag"`
	CityName     string    `xorm:"varchar(24)" json:"city_name"`
	City         *City     `xorm:"- <- ->"`
}

/**
 * 将数据库结果格式组装成 request 需要的json 格式  map
 */
func (u *User) UserToRespDesc() interface{} {
	respInfo := map[string]interface{}{
		"id":           u.Id,
		"user_id":      u.Id,
		"username":     u.UserName,
		"city":         u.CityName,
		"registe_time": util.FormatDatetime(u.RegisterTime),
		"moblie":       u.Mobile,
		"is_active":    u.IsActive,
		"balance":      u.Balance,
		"avatar":       u.Avatar,
	}
	return respInfo
}
