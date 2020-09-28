package datasource

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
	"iris/project/config"
)

/**
 * 返回Redis实例
 */
func NewRedis() *redis.Database  {
	var database *redis.Database


	//项目配置
	proConfig := config.InitConfig()
	if proConfig !=nil {
		iris.New().Logger().Info("redis.go + hello")
		rd:=proConfig.Redis
		iris.New().Logger().Info(rd)
		database = redis.New(redis.Config{
			Network: rd.NetWork,
			Addr: rd.Addr,
			Password: rd.Password,
			Database: "",
			MaxActive: 10,
			Timeout: redis.DefaultRedisTimeout,
			Prefix: rd.Prefix,
			Driver: redis.Redigo(),
		})
	}else{
		iris.New().Logger().Info("redis.go + error")
	}
	return database
}

