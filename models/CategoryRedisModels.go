package models

import (
	"beegoblog/tools"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"strings"
	"time"
)

///region  CategoryRedis
func AddCategoryRedis(name string) error {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()
	//判断是否已经存在该分类
	categoryKeys, err := redis.Values(conn.Do("HKEYS", "category"))
	titleValues := make([]string, 0)
	for _, categoryKey := range categoryKeys {
		categoryKeyStr := string(categoryKey.([]byte))
		if strings.Contains(categoryKeyStr, "_Title") {
			titleValue, _ := conn.Do("HGET", "category", categoryKeyStr+"_Title")
			titleValues = append(titleValues, string(titleValue.([]byte)))
		}
	}

	for i := 0; i < len(titleValues); i++ {
		if titleValues[i] == name {
			return nil
		}

	}

	//新增一个分类
	guid := tools.GetGuid() //得到GUID
	timeNow := time.Now()
	conn.Do("HMSET", "category",
		guid+"_Id", guid,
		guid+"_Title", name,
		guid+"_Views", 0,
		guid+"_TopicCount", 0,
		guid+"_Created", timeNow,
		guid+"_Updated", timeNow)
	return nil
}

//TODO:删除分类后相应文章的删除
func DeleteCategoryRedis(id string) error {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()

	//删除Category
	conn.Do("HDEL", "category",
		id+"_Id",
		id+"_Title",
		id+"_Views",
		id+"_TopicCount",
		id+"_Created",
		id+"_Updated")
	return nil
}

//TODO:isListAll为false的情况完善
func GetAllCategoriesRedis(isListAll bool) (categories []*Category, err error) {
	_ = isListAll
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()

	categoryKeys, err := redis.Values(conn.Do("HKEYS", "category"))
	idValues := make([]string, 0)
	for _, categoryKey := range categoryKeys {
		categoryKeyStr := string(categoryKey.([]byte))
		if strings.Contains(categoryKeyStr, "_Id") {
			idValue, _ := conn.Do("HGET", "category", categoryKeyStr+"_Id")
			idValues = append(idValues, string(idValue.([]byte)))
		}
	}
	for _, idValue := range idValues {
		category := new(Category)
		Id, _ := conn.Do("HGET", "category", idValue+"_Id")
		Title, _ := conn.Do("HGET", "category", idValue+"_Title")
		Views, _ := conn.Do("HGET", "category", idValue+"_Views")
		TopicCount, _ := conn.Do("HGET", "category", idValue+"_TopicCount")
		//Created, _ := conn.Do("HGET", "category", idValue+"_Created")
		//Updated, _ := conn.Do("HGET", "category", idValue+"_Updated")

		category.Id = int64(Id.(int64))
		category.Title = string(Title.([]byte))
		category.Views = int64(Views.(int64))
		category.TopicCount = int64(TopicCount.(int64))
		//category.Created =
		//category.Updated =
		categories = append(categories, category)
	}

	return
}

///endReigon
