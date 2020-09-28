package util

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/wonderivan/logger"
	"os"
	"reflect"
	"time"
)

//根据Json格式设置object对象
func SetObjByJson(obj interface{}, data map[string]interface{}) error {
	//遍历map
	for key, value := range data {
		if err := setField(obj, key, value); err != nil {
			logger.Error("SetObjByJson set field fail.")
			return err
		}
	}
	return nil
}

//设置结构体中的变量  通过反射给结构体变量赋值，也就是构造，或者java--set方法
func setField(obj interface{}, name string, value interface{}) error {
	structData := reflect.ValueOf(obj).Type().Elem()
	fmt.Println(structData, "---> utils.go-->setField")
	fieldValue, result := structData.FieldByName(name)
	if !result {
		logger.Error("No such field ", name)
		return fmt.Errorf("No such filed %s", name)
	}

	//结构体中变量的类型
	fieldType := fieldValue.Type
	//参数的值
	val := reflect.ValueOf(value)
	//参数的类型
	valTypeStr := val.Type().String()

	//结构体重变量的类型
	fieldTypeStr := fieldType.String()
	// float64 to int
	if valTypeStr == "float64" && fieldTypeStr == "int" {
		val = val.Convert(fieldType)
	}

	//类型必须匹配
	if fieldType != val.Type() {
		return fmt.Errorf("value type %s didn't match obj field type %s ", valTypeStr, fieldTypeStr)
	}

	reflect.ValueOf(obj).Elem().FieldByName(name).Set(val)
	return nil

}

func LogInfo(app *iris.Application, v ...interface{}) {
	app.Logger().Info(v)
}

func LogError(app *iris.Application, v ...interface{}) {
	app.Logger().Error(v)
}

func LogDebug(app *iris.Application, v ...interface{}) {
	app.Logger().Debug(v)
}

/**
 * 格式化数据
 */
func FormatDatetime(time time.Time) string {
	return time.Format("2006-01-02 03:04:05")
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
