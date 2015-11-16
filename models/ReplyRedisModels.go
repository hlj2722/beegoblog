package models

import (
	"beegoblog/tools"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"strings"
	"time"
)

///region ReplyRedis
func AddReplyRedis(tid, nickname, content string) error {

	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()
	guid := tools.GetGuid() //得到GUID
	timeNow := time.Now()
	conn.Do("HMSET", "reply",
		guid+"_Id", guid,
		guid+"_Tid", tid,
		guid+"_Name", nickname,
		guid+"_Content", content,
		guid+"_Created", timeNow,
		guid+"_Updated", timeNow)

	conn.Do("HSET", "topic", tid+"_ReplyTime", timeNow)
	conn.Do("HINCRBY", "topic", tid+"_ReplyCount", 1)

	return nil
}

func GetAllRepliesRedis(tid string) (replies []*Reply, err error) {

	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()

	replyKeys, err := redis.Values(conn.Do("HKEYS", "reply"))
	idValues := make([]string, 0)
	for _, replyKey := range replyKeys {
		replyKeyStr := string(replyKey.([]byte))
		if strings.Contains(replyKeyStr, "_Id") {
			idValue, _ := conn.Do("HGET", "reply", replyKeyStr+"_Id")
			idValues = append(idValues, string(idValue.([]byte)))
		}
	}
	for _, idValue := range idValues {
		reply := new(Reply)
		Id, _ := conn.Do("HGET", "reply", idValue+"_Id")
		Tid, _ := conn.Do("HGET", "reply", idValue+"_Tid")
		Name, _ := conn.Do("HGET", "reply", idValue+"_Name")
		Content, _ := conn.Do("HGET", "reply", idValue+"_Content")
		//Created, _ := conn.Do("HGET", "reply", idValue+"_Created")
		//Updated, _ := conn.Do("HGET", "reply", idValue+"_Updated")

		reply.Id = int64(Id.(int64))
		reply.Tid = int64(Tid.(int64))
		reply.Name = string(Name.([]byte))
		reply.Content = string(Content.([]byte))
		//reply.Created =
		//reply.Updated =
		replies = append(replies, reply)
	}

	return
}

func DeleteReplyRedis(rid string) error {

	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()
	//暂存Tid
	topicId, _ := conn.Do("HGET", "reply", rid+"_Tid")
	tid := string(topicId.([]byte))

	//删除Reply
	conn.Do("HDEL", "reply",
		rid+"_Id",
		rid+"_Tid",
		rid+"_Name",
		rid+"_Content",
		rid+"_Created",
		rid+"_Updated")

	//获取评论数量
	count, _ := conn.Do("HLEN", "reply")
	replyCount := int(count.(int)) / 6

	//获取最近的评论时间
	replyKeys, err := redis.Values(conn.Do("HKEYS", "reply"))
	updatedValues := make([]string, 0)
	for _, replyKey := range replyKeys {
		replyKeyStr := string(replyKey.([]byte))
		if strings.Contains(replyKeyStr, "_Updated") {
			updatedValue, _ := conn.Do("HGET", "reply", replyKeyStr+"_Updated")
			updatedValues = append(updatedValues, string(updatedValue.([]byte)))
		}
	}
	var replyTime string
	for i := 0; i < len(updatedValues); i++ {
		if updatedValues[i] > replyTime {
			replyTime = updatedValues[i]
		}

	}

	//更新Topic
	conn.Do("HSET", "topic", tid+"_ReplyCount", replyCount)
	conn.Do("HSET", "topic", tid+"_ReplyTime", replyTime)
	return nil

}

///endRegion
