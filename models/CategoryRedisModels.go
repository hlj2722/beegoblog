package models

import ()

///region  CategoryRedis
func AddCategoryRedis(name string) error {

	return nil
}

//删除分类的同时，包含该分类的所有文章都删除
//采用事务避免出问题
func DeleteCategoryRedis(id string) error {
	return nil

}

func GetAllCategoriesRedis(isListAll bool) ([]*Category, error) {

	return nil, nil
}

///endReigon
