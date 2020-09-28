package service

import (
	"fmt"
	"iris/project/model"
	"xorm.io/xorm"
)

/**
 * 管理员服务
 * 标准开发模式下，每个实体提供的功能以接口标准形式定义。供控制层调用
 */

type AdminService interface {

	//通过name pwd 获取Admin
	GetAdminByAmdinNameAndPassword(username, password string) (model.Admin, bool)

	//获取管理员总数
	GetAdminCount() (int64, error)

	SaveAvatarImg(adminId int64, fileName string) bool
}

/**
 * 管理员服务实现结构体(小写，需要工厂模式)
 */
type adminService struct {
	engine *xorm.Engine
}

func NewAdminService(db *xorm.Engine) AdminService {
	return &adminService{
		engine: db,
	}
}

//adminService 得实现接口 接口的实现类 。。。。。。。在这

func (a *adminService) GetAdminByAmdinNameAndPassword(username, password string) (model.Admin, bool) {
	var admin model.Admin

	//a.engine.Where("admin_name =?","root").And("pwd=?","123").Get(&admin)
	a.engine.Where(" admin_name = ? and pwd = ? ", "root1", "123").Get(&admin)
	fmt.Println(a.engine)
	//a.engine.SQL("select * from where admin_name=?",username).Get(&admin)
	return admin, admin.AdminId != 0
}

func (a *adminService) GetAdminCount() (int64, error) {
	count, err := a.engine.Count(new(model.Admin))
	if err != nil {
		panic(err.Error())
	}
	return count, nil
}

func (ac *adminService) SaveAvatarImg(adminId int64, fileName string) bool {
	admin := model.Admin{
		Avator: fileName,
	}
	_, err := ac.engine.ID(adminId).Cols("avatar").Update(&admin)
	return err != nil
}
