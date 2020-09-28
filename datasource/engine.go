package datasource

import (
	_ "github.com/go-sql-driver/mysql"
	"iris/project/model"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"
)


func main() {
	engine, err := xorm.NewEngine("mysql", "root:1234a@/go?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	engine.SetMapper(names.GonicMapper{}) //设置映射规则，智能的。 name+驼峰
	engine.Sync2(
		new(model.Admin),
		new(model.City),
		new(model.AdminPermission),
		new(model.User),
		new(model.UserOrder),
		new(model.Permission),
		new(model.Address),
		new(model.OrderStatus),
		new(model.Shop),
		new(model.Food),
		new(model.FoodCategory),
		)
}


func NewMySqlEngine() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", "root:1234a@/go?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	engine.Sync2(
		new(model.Admin),
		new(model.City),
		new(model.AdminPermission),
		new(model.User),
		new(model.UserOrder),
		new(model.Permission),
		new(model.Address),
		new(model.OrderStatus),
		new(model.Shop),
		new(model.Food),
		new(model.FoodCategory),
	)

	engine.ShowSQL(true)
	engine.SetMaxOpenConns(10)
	engine.SetLogLevel(log.LOG_DEBUG)

	return engine
}