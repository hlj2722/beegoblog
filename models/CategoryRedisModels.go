package models

import (
	_ "beegoblog/tools"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"strconv"
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
	if err != nil {
		return err
	}
	for _, categoryKey := range categoryKeys {
		categoryKeyStr := string(categoryKey.([]byte))
		if strings.Contains(categoryKeyStr, "_Title") {
			titleValue, _ := conn.Do("HGET", "category", categoryKeyStr)
			if string(titleValue.([]byte)) == name {
				return nil
			}
		}
	}
	//新增一个分类
	guid, _ := conn.Do("HINCRBY", "category", "guid", 1) //生成Guid,并保存到键category的guid域
	guidStr := strconv.FormatInt(int64(guid.(int64)), 10)

	timeNow := time.Now()
	conn.Do("HMSET", "category",
		guidStr+"_Id", guidStr,
		guidStr+"_Title", name,
		guidStr+"_Views", 0,
		guidStr+"_TopicCount", 0,
		guidStr+"_Created", timeNow.Format("2006-01-02 15:04:05"),
		guidStr+"_Updated", timeNow.Format("2006-01-02 15:04:05"))
	return nil
}

func DeleteCategoryRedis(id string) error {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()
	//暂存删除的分类
	category, _ := conn.Do("HGET", "category", id+"_Title")
	categoryStr := string(category.([]byte))
	//删除Category
	conn.Do("HDEL", "category", id+"_Id")
	conn.Do("HDEL", "category", id+"_Title")
	conn.Do("HDEL", "category", id+"_Views")
	conn.Do("HDEL", "category", id+"_TopicCount")
	conn.Do("HDEL", "category", id+"_Created")
	conn.Do("HDEL", "category", id+"_Updated")

	//删除分类下的所有文章
	topicKeys, err := redis.Values(conn.Do("HKEYS", "topic"))
	if err != nil {
		return err
	}

	for _, topicKey := range topicKeys {
		topicKeyStr := string(topicKey.([]byte))
		if strings.Contains(topicKeyStr, "_Category") {
			categoryValue, _ := conn.Do("HGET", "topic", topicKeyStr)
			categoryValueStr := string(categoryValue.([]byte))

			if categoryValueStr == categoryStr {
				idStr := strings.TrimRight(topicKeyStr, "_Category")

				conn.Do("HDEL", "topic", idStr+"_Id")
				conn.Do("HDEL", "topic", idStr+"_Title")
				conn.Do("HDEL", "topic", idStr+"_Category")
				conn.Do("HDEL", "topic", idStr+"_Lables")
				conn.Do("HDEL", "topic", idStr+"_Content")
				conn.Do("HDEL", "topic", idStr+"_Attachment")
				conn.Do("HDEL", "topic", idStr+"_Views")
				conn.Do("HDEL", "topic", idStr+"_Author")
				conn.Do("HDEL", "topic", idStr+"_ReplyTime")
				conn.Do("HDEL", "topic", idStr+"_ReplyCount")
				conn.Do("HDEL", "topic", idStr+"_Created")
				conn.Do("HDEL", "topic", idStr+"_Updated")
			}

		}
	}

	return nil
}

func GetAllCategoriesRedis(isListAll bool) (categories []*Category, err error) {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()

	categoryKeys, err := redis.Values(conn.Do("HKEYS", "category"))
	if err != nil {
		return nil, err
	}

	for _, categoryKey := range categoryKeys {
		categoryKeyStr := string(categoryKey.([]byte))
		if strings.Contains(categoryKeyStr, "_Id") {
			idValue, _ := conn.Do("HGET", "category", categoryKeyStr)
			idValueStr := string(idValue.([]byte))

			category := new(Category)
			Id, _ := conn.Do("HGET", "category", idValueStr+"_Id")
			Title, _ := conn.Do("HGET", "category", idValueStr+"_Title")
			Views, _ := conn.Do("HGET", "category", idValueStr+"_Views")
			TopicCount, _ := conn.Do("HGET", "category", idValueStr+"_TopicCount")
			Created, _ := conn.Do("HGET", "category", idValueStr+"_Created")
			Updated, _ := conn.Do("HGET", "category", idValueStr+"_Updated")

			category.Id, _ = strconv.ParseInt(string(Id.([]byte)), 10, 0)
			category.Title = string(Title.([]byte))
			category.Views, _ = strconv.ParseInt(string(Views.([]byte)), 10, 0)
			category.TopicCount, _ = strconv.ParseInt(string(TopicCount.([]byte)), 10, 0)
			category.Created, _ = time.Parse("2006-01-02 15:04:05", string(Created.([]byte)))
			category.Updated, _ = time.Parse("2006-01-02 15:04:05", string(Updated.([]byte)))

			if isListAll {
				categories = append(categories, category)

			} else if category.TopicCount > 0 {
				categories = append(categories, category)
			}

		}
	}

	return
}

///endReigon
