package models

import (
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"strings"
	"time"
)

///region TopicRedis
func AddTopicRedis(title, category, lable, content, attachment, author string) error {
	// 处理标签
	lable = "$" + strings.Join(strings.Split(lable, " "), "#$") + "#"

	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()

	//新增一个文章
	guid, _ := conn.Do("HINCRBY", "topic", "guid", 1) //生成Guid,并保存到键topic的guid域
	guidStr := string(guid.([]byte))
	timeNow := time.Now()
	conn.Do("HMSET", "topic",
		guidStr+"_Id", guidStr,
		guidStr+"_Title", title,
		guidStr+"_Category", category,
		guidStr+"_Lables", lable,
		guidStr+"_Content", content,
		guidStr+"_Attachment", attachment,
		guidStr+"_Views", 0,
		guidStr+"_Author", author,
		guidStr+"_ReplyTime", timeNow,
		guidStr+"_ReplyCount", 0,
		guidStr+"_Created", timeNow.Format("2006-01-02 15:04:05"),
		guidStr+"_Updated", timeNow.Format("2006-01-02 15:04:05"))

	//更新分类统计
	categoryKeys, err := redis.Values(conn.Do("HKEYS", "category"))
	if err != nil {
		return err
	}

	existsCategory := false

	for _, categoryKey := range categoryKeys {
		categoryKeyStr := string(categoryKey.([]byte))
		if strings.Contains(categoryKeyStr, "_Title") {
			titleValue, _ := conn.Do("HGET", "category", categoryKeyStr)
			titleValueStr := string(titleValue.([]byte))

			if titleValueStr == category {
				existsCategory = true
				topicCountKeyStr := strings.TrimRight(categoryKeyStr, "_Title") + "_TopicCount"
				topicCount, _ := conn.Do("HINCRBY", "category", topicCountKeyStr, 1)
				if int(topicCount.(int)) < 1 {
					conn.Do("HSET", "category", topicCountKeyStr, 1)
				}

			}

		}
	}

	if !existsCategory {
		//新增一个分类
		guid, _ := conn.Do("HINCRBY", "category", "guid", 1) //生成Guid,并保存到键category的guid域
		guidStr := string(guid.([]byte))
		timeNow := time.Now()
		conn.Do("HMSET", "category",
			guidStr+"_Id", guidStr,
			guidStr+"_Title", category,
			guidStr+"_Views", 0,
			guidStr+"_TopicCount", 0,
			guidStr+"_Created", timeNow.Format("2006-01-02 15:04:05"),
			guidStr+"_Updated", timeNow.Format("2006-01-02 15:04:05"))

	}

	return nil
}

func GetTopicRedis(tid string) (*Topic, error) {

	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()

	topic := new(Topic)
	Id, _ := conn.Do("HGET", "topic", tid+"_Id")
	Title, _ := conn.Do("HGET", "topic", tid+"_Title")
	Category, _ := conn.Do("HGET", "topic", tid+"_Category")
	Lables, _ := conn.Do("HGET", "topic", tid+"_Lables")
	Content, _ := conn.Do("HGET", "topic", tid+"_Content")
	Attachment, _ := conn.Do("HGET", "topic", tid+"_Attachment")
	Views, _ := conn.Do("HINCRBY", "topic", tid+"_Views", 1)
	Author, _ := conn.Do("HGET", "topic", tid+"__Author")
	ReplyTime, _ := conn.Do("HGET", "topic", tid+"_ReplyTime")
	ReplyCount, _ := conn.Do("HGET", "topic", tid+"_ReplyCount")
	Created, _ := conn.Do("HGET", "topic", tid+"_Created")
	Updated, _ := conn.Do("HGET", "topic", tid+"_Updated")

	topic.Id = int64(Id.(int64))
	topic.Title = string(Title.([]byte))
	topic.Category = string(Category.([]byte))
	topic.Lables = strings.Replace(strings.Replace(string(Lables.([]byte)), "#", " ", -1), "$", "", -1)
	topic.Content = string(Content.([]byte))
	topic.Attachment = string(Attachment.([]byte))
	topic.Views = int64(Views.(int64))
	topic.Author = string(Author.([]byte))
	topic.ReplyTime, _ = time.Parse("2006-01-02 15:04:05", string(ReplyTime.([]byte)))
	topic.ReplyCount = int64(ReplyCount.(int64))
	topic.Created, _ = time.Parse("2006-01-02 15:04:05", string(Created.([]byte)))
	topic.Updated, _ = time.Parse("2006-01-02 15:04:05", string(Updated.([]byte)))

	return topic, nil
}

func ModifyTopicRedis(tid, title, category, lable, content, attachment string) error {

	return AddTopicRedis(title, category, lable, content, attachment, "")
}

func DeleteTopicRedis(tid string) error {

	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()

	//暂存删除的分类
	category, _ := conn.Do("HGET", "topic", tid+"_Category")

	//删除分类
	conn.Do("HDEL", "topic",
		tid+"_Id",
		tid+"_Title",
		tid+"_Category",
		tid+"_Lables",
		tid+"_Content",
		tid+"_Attachment",
		tid+"_Views",
		tid+"_Author",
		tid+"_ReplyTime",
		tid+"_ReplyCount",
		tid+"_Created",
		tid+"_Updated")

	topicKeys, err := redis.Values(conn.Do("HKEYS", "topic"))
	if err != nil {
		return err
	}

	for _, topicKey := range topicKeys {
		topicKeyStr := string(topicKey.([]byte))
		if strings.Contains(topicKeyStr, "_Category") {
			categoryValue, _ := conn.Do("HGET", "topic", topicKeyStr)
			categoryValueStr := string(categoryValue.([]byte))
			if categoryValueStr == category {
				conn.Do("HINCRBY", "topic", strings.TrimRight(topicKeyStr, "_Category")+"_TopicCount", -1)
			}

		}
	}

	return nil
}

//TODO:记录排序
func GetAllTopicsRedis(category, lable string, isDesc bool) (topics []*Topic, err error) {

	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()

	topicKeys, err := redis.Values(conn.Do("HKEYS", "topic"))
	if err != nil {
		return nil, err
	}

	for _, topicKey := range topicKeys {
		topicKeyStr := string(topicKey.([]byte))
		if strings.Contains(topicKeyStr, "_Id") {
			idValue, _ := conn.Do("HGET", "topic", topicKeyStr)
			idValueStr := string(idValue.([]byte))

			topic := new(Topic)
			Id, _ := conn.Do("HGET", "topic", idValueStr+"_Id")
			Title, _ := conn.Do("HGET", "topic", idValueStr+"_Title")
			Category, _ := conn.Do("HGET", "topic", idValueStr+"_Category")
			Lables, _ := conn.Do("HGET", "topic", idValueStr+"_Lables")
			Content, _ := conn.Do("HGET", "topic", idValueStr+"_Content")
			Attachment, _ := conn.Do("HGET", "topic", idValueStr+"_Attachment")
			Views, _ := conn.Do("HGET", "topic", idValueStr+"_Views")
			Author, _ := conn.Do("HGET", "topic", idValueStr+"__Author")
			ReplyTime, _ := conn.Do("HGET", "topic", idValueStr+"_ReplyTime")
			ReplyCount, _ := conn.Do("HGET", "topic", idValueStr+"_ReplyCount")
			Created, _ := conn.Do("HGET", "topic", idValueStr+"_Created")
			Updated, _ := conn.Do("HGET", "topic", idValueStr+"_Updated")

			topic.Id = int64(Id.(int64))
			topic.Title = string(Title.([]byte))
			topic.Category = string(Category.([]byte))
			topic.Lables = string(Lables.([]byte))
			topic.Content = string(Content.([]byte))
			topic.Attachment = string(Attachment.([]byte))
			topic.Views = int64(Views.(int64))
			topic.Author = string(Author.([]byte))
			topic.ReplyTime, _ = time.Parse("2006-01-02 15:04:05", string(ReplyTime.([]byte)))
			topic.ReplyCount = int64(ReplyCount.(int64))
			topic.Created, _ = time.Parse("2006-01-02 15:04:05", string(Created.([]byte)))
			topic.Updated, _ = time.Parse("2006-01-02 15:04:05", string(Updated.([]byte)))

			if isDesc {
				if len(category) > 0 {
					if topic.Category == category {
						topics = append(topics, topic)
					}
				}
				if len(lable) > 0 {
					lable = strings.Trim(lable, " ")

					if strings.Contains(topic.Lables, "$"+lable+"#") {
						topics = append(topics, topic)
					}
				}

			} else {
				topics = append(topics, topic)
			}

		}
	}

	return

}

///endRegion
