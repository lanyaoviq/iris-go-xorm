package config

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	AppName string `json:"app_name"`
	Port string `json:"port"`
	StaticPath string `json:"static_path"`
	Mode string `json:"mode"`
}

var ServConfig AppConfig


//初始化服务器配置
func InitConfig() *AppConfig{
	file,err:=os.Open("/Users/x/goPro/iris/project/config.json")
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	decoder :=json.NewDecoder(file)
	conf :=AppConfig{}
	decoder.Decode(&conf)
	if err != nil {
		panic(err.Error())
	}
	return &conf
}


