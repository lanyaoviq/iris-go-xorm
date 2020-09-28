package model

import "time"

//定义管理员结构体
type Admin struct {
	AdminId int64 `xorm:"pk autoincr" json:"id"`
	AdminName string `xorm:"varchar(32)" json:"admin_name"`
	CreateTime time.Time `xorm:"DateTime" json:"create_time"`
	Status int64 `xorm:"default 0" json:"status"`
	Avator string `xorm:"varchar(255)" json:"avator"`
	Pwd string `xorm:"varchar(255)" json:"pwd"`
	CityName string `xorm:"varchar(12)" json:"city_name" `
	CityId int64 `xorm:"index" json:"city_id"`
	City *City `xorm:"- <- ->"`
}


/**
 * 从admin数据库实体转换为前端请求的response的json格式, map类型
 */
func (this *Admin) AdminToRespDesc() interface{}{
	respDesc := map[string]interface{}{
		"user_name":this.AdminName,
		"id":this.AdminId,
		"create_time":this.CreateTime,
		"status":this.Status,
		"avatar":this.Avator,
		"city":this.CityName,
		"admin":"管理员",
	}
	return  respDesc
}
