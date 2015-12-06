package models

import (
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"os"
	"path"
	"strconv"
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

	timeNow := time.Now()

	//新增一个文章
	guid, _ := conn.Do("HINCRBY", "topic", "guid", 1) //生成Guid,并保存到键topic的guid域
	guidStr := strconv.FormatInt(int64(guid.(int64)), 10)

	conn.Do("HMSET", "topic",
		guidStr+"_Id", guidStr,
		guidStr+"_Title", title,
		guidStr+"_Category", category,
		guidStr+"_Lables", lable,
		guidStr+"_Content", content,
		guidStr+"_Attachment", attachment,
		guidStr+"_Views", 0,
		guidStr+"_Author", author,
		guidStr+"_ReplyTime", timeNow.Format("2006-01-02 15:04:05"),
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
				topicCountInt := int64(topicCount.(int64))
				if topicCountInt < 1 {
					conn.Do("HSET", "category", topicCountKeyStr, 1)
				}
				conn.Do("HSET", "category", strings.TrimRight(categoryKeyStr, "_Title")+"_Updated",
					timeNow.Format("2006-01-02 15:04:05"))

			}

		}
	}

	if !existsCategory {
		//新增一个分类
		guid, _ := conn.Do("HINCRBY", "category", "guid", 1) //生成Guid,并保存到键category的guid域
		guidStr := strconv.FormatInt(int64(guid.(int64)), 10)

		conn.Do("HMSET", "category",
			guidStr+"_Id", guidStr,
			guidStr+"_Title", category,
			guidStr+"_Views", 0,
			guidStr+"_TopicCount", 1,
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
	Author, _ := conn.Do("HGET", "topic", tid+"_Author")
	ReplyTime, _ := conn.Do("HGET", "topic", tid+"_ReplyTime")
	ReplyCount, _ := conn.Do("HGET", "topic", tid+"_ReplyCount")
	Created, _ := conn.Do("HGET", "topic", tid+"_Created")
	Updated, _ := conn.Do("HGET", "topic", tid+"_Updated")

	topic.Id, _ = strconv.ParseInt(string(Id.([]byte)), 10, 0)
	topic.Title = string(Title.([]byte))
	topic.Category = string(Category.([]byte))
	topic.Lables = strings.Replace(strings.Replace(string(Lables.([]byte)), "#", " ", -1), "$", "", -1)
	topic.Content = string(Content.([]byte))
	topic.Attachment = string(Attachment.([]byte))
	topic.Views = int64(Views.(int64))
	topic.Author = string(Author.([]byte))
	topic.ReplyTime, _ = time.Parse("2006-01-02 15:04:05", string(ReplyTime.([]byte)))
	topic.ReplyCount, _ = strconv.ParseInt(string(ReplyCount.([]byte)), 10, 0)
	topic.Created, _ = time.Parse("2006-01-02 15:04:05", string(Created.([]byte)))
	topic.Updated, _ = time.Parse("2006-01-02 15:04:05", string(Updated.([]byte)))

	return topic, nil
}

func ModifyTopicRedis(tid, title, category, lable, content, attachment string) error {

	lable = "$" + strings.Join(strings.Split(lable, " "), "#$") + "#"

	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()

	oldCate, _ := conn.Do("HGET", "topic", tid+"_Category")
	oldAttach, _ := conn.Do("HGET", "topic", tid+"_Attachment")
	oldCateStr := string(oldCate.([]byte))
	oldAttachStr := string(oldAttach.([]byte))

	//更新Topic
	timeNow := time.Now()
	conn.Do("HMSET", "topic",
		tid+"_Title", title,
		tid+"_Category", category,
		tid+"_Lables", lable,
		tid+"_Content", content,
		tid+"_Attachment", attachment,
		tid+"_Updated", timeNow.Format("2006-01-02 15:04:05"))

	// 更新分类统计
	if oldCateStr != category {
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
					topicCountInt := int64(topicCount.(int64))
					if topicCountInt < 1 {
						conn.Do("HSET", "category", topicCountKeyStr, 1)
					}

				}

				if titleValueStr == oldCateStr {
					topicCountKeyStr := strings.TrimRight(categoryKeyStr, "_Title") + "_TopicCount"
					topicCount, _ := conn.Do("HINCRBY", "category", topicCountKeyStr, -1)
					topicCountInt := int64(topicCount.(int64))
					if topicCountInt < 0 {
						conn.Do("HSET", "category", topicCountKeyStr, 0)
					}

				}

				conn.Do("HSET", "category", strings.TrimRight(categoryKeyStr, "_Title")+"_Updated",
					timeNow.Format("2006-01-02 15:04:05"))

			}
		}

		if !existsCategory {
			//新增一个分类
			guid, _ := conn.Do("HINCRBY", "category", "guid", 1) //生成Guid,并保存到键category的guid域
			guidStr := strconv.FormatInt(int64(guid.(int64)), 10)
			conn.Do("HMSET", "category",
				guidStr+"_Id", guidStr,
				guidStr+"_Title", category,
				guidStr+"_Views", 0,
				guidStr+"_TopicCount", 1,
				guidStr+"_Created", timeNow.Format("2006-01-02 15:04:05"),
				guidStr+"_Updated", timeNow.Format("2006-01-02 15:04:05"))

		}
	}

	// 删除旧的附件
	if oldAttachStr != attachment {
		os.Remove(path.Join("attachment", oldAttachStr))
	}

	return nil
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
	categoryStr := string(category.([]byte))

	//删除分类
	conn.Do("HDEL", "topic", tid+"_Id")
	conn.Do("HDEL", "topic", tid+"_Title")
	conn.Do("HDEL", "topic", tid+"_Category")
	conn.Do("HDEL", "topic", tid+"_Lables")
	conn.Do("HDEL", "topic", tid+"_Content")
	conn.Do("HDEL", "topic", tid+"_Attachment")
	conn.Do("HDEL", "topic", tid+"_Views")
	conn.Do("HDEL", "topic", tid+"_Author")
	conn.Do("HDEL", "topic", tid+"_ReplyTime")
	conn.Do("HDEL", "topic", tid+"_ReplyCount")
	conn.Do("HDEL", "topic", tid+"_Created")
	conn.Do("HDEL", "topic", tid+"_Updated")

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
			Author, _ := conn.Do("HGET", "topic", idValueStr+"_Author")
			ReplyTime, _ := conn.Do("HGET", "topic", idValueStr+"_ReplyTime")
			ReplyCount, _ := conn.Do("HGET", "topic", idValueStr+"_ReplyCount")
			Created, _ := conn.Do("HGET", "topic", idValueStr+"_Created")
			Updated, _ := conn.Do("HGET", "topic", idValueStr+"_Updated")

			topic.Id, _ = strconv.ParseInt(string(Id.([]byte)), 10, 0)
			topic.Title = string(Title.([]byte))
			topic.Category = string(Category.([]byte))
			topic.Lables = string(Lables.([]byte))
			topic.Content = string(Content.([]byte))
			topic.Attachment = string(Attachment.([]byte))
			topic.Views, _ = strconv.ParseInt(string(Views.([]byte)), 10, 0)
			topic.Author = string(Author.([]byte))
			topic.ReplyTime, _ = time.Parse("2006-01-02 15:04:05", string(ReplyTime.([]byte)))
			topic.ReplyCount, _ = strconv.ParseInt(string(ReplyCount.([]byte)), 10, 0)
			topic.Created, _ = time.Parse("2006-01-02 15:04:05", string(Created.([]byte)))
			topic.Updated, _ = time.Parse("2006-01-02 15:04:05", string(Updated.([]byte)))

			if isDesc && len(category) > 0 {
				if topic.Category == category {
					topics = append(topics, topic)
				}
			} else if isDesc && len(lable) > 0 {
				lable = strings.Trim(lable, " ")

				if strings.Contains(topic.Lables, "$"+lable+"#") {
					topics = append(topics, topic)
				}
			} else {
				topics = append(topics, topic)
			}

		}
	}

	return

}

///endRegion
