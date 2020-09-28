package service

import (
	"github.com/kataras/iris/v12"
	"iris/project/model"
	"xorm.io/xorm"
)

type CategoryServie interface {
	AddCategory(model *model.FoodCategory) bool
	GetCategoryByShopId(shopId int64) ([]model.FoodCategory, error)
	SaveShop(shop model.Shop) bool
	DeleteShop(shopId int64) bool

	GetAllCategory() ([]model.FoodCategory, error)
}

type categoryService struct {
	Engine *xorm.Engine
}

/**
 * @Author x
 * @Description //TODO
 * @Date 下午9:06 2020/9/27
 * @Param
 * @return
 **/
func (cs *categoryService) GetAllCategory() ([]model.FoodCategory, error) {
	var categorys []model.FoodCategory
	err := cs.Engine.Find(&categorys)
	if err != nil {
		iris.New().Logger().Error(err.Error())
		return nil, err
	}
	return categorys, nil
}

func (cs *categoryService) DeleteShop(shopId int64) bool {
	_, err := cs.Engine.Delete(shopId)
	if err != nil {
		iris.New().Logger().Error(err.Error())
	}
	return err == nil
}

func (cs *categoryService) SaveShop(shop model.Shop) bool {
	_, err := cs.Engine.Insert(&shop)
	if err != nil {
		iris.New().Logger().Error(err.Error())
	}
	return err == nil
}

func NewCategoryService(db *xorm.Engine) CategoryServie {
	return &categoryService{
		Engine: db,
	}
}

/**
 * 通过商铺ID 获取食品品类slcie
 */
func (cs *categoryService) GetCategoryByShopId(shopId int64) ([]model.FoodCategory, error) {
	categories := make([]model.FoodCategory, 0)
	err := cs.Engine.Where("restaurant_id=?", shopId).Find(&categories)
	if err != nil {
		iris.New().Logger().Error(err.Error())
	}
	iris.New().Logger().Info(categories)
	return categories, err

}

/**
 * 获取品类数量
 */
func (cs *categoryService) AddCategory(category *model.FoodCategory) bool {
	iris.New().Logger().Info(category)
	_, err := cs.Engine.Insert(&category)
	if err != nil {
		iris.New().Logger().Info(err.Error())
	}
	return err != nil
}
