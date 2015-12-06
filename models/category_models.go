package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"time"
)

///region  Category
func AddCategory(name string) error {
	o := orm.NewOrm()

	cate := &Category{
		Title:   name,
		Created: time.Now(),
		Updated: time.Now(),
	}

	// 查询数据
	qs := o.QueryTable("category")
	err := qs.Filter("title", name).One(cate)
	if err == nil {
		return err
	}

	// 插入数据
	_, err = o.Insert(cate)
	if err != nil {
		return err
	}

	return nil
}

//删除分类的同时，包含该分类的所有文章都删除
//采用事务避免出问题
func DeleteCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	cateDel := &Category{Id: cid}
	cateTemp := new(Category)

	err = o.QueryTable("category").Filter("id", cid).One(cateTemp)
	if err != nil {
		return err
	}

	//事务中删除分类和该分类的所有文章
	if errB := o.Begin(); errB == nil {
		if _, errDelTopic := o.QueryTable("topic").Filter("Category", cateTemp.Title).Delete(); errDelTopic == nil {

			if _, errDelCate := o.Delete(cateDel); errDelCate == nil {
				if errC := o.Commit(); errC == nil {
					return nil
				} else {
					if errR := o.Rollback(); errR == nil {
						return errC
					} else {
						return errR
					}
				}
			} else {
				return errDelCate
			}
		} else {
			return errDelTopic
		}
	} else {
		return errB
	}

}

func GetAllCategories(isListAll bool) ([]*Category, error) {
	o := orm.NewOrm()
	cates := make([]*Category, 0)
	qs := o.QueryTable("category")
	var err error
	if isListAll {
		_, err = qs.All(&cates) //过滤得到文章数量TopicCount大于0的分类

	} else {
		_, err = qs.Filter("TopicCount__gt", 0).All(&cates) //过滤得到文章数量TopicCount大于0的分类
	}
	return cates, err
}

///endReigon
