package models

import (
	_ "beegoblog/tools"
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
	guid, _ := conn.Do("HINCRBY", "reply", "guid", 1)
	guidStr := string(guid.([]byte))
	timeNow := time.Now()
	conn.Do("HMSET", "reply",
		guidStr+"_Id", guidStr,
		guidStr+"_Tid", tid,
		guidStr+"_Name", nickname,
		guidStr+"_Content", content,
		guidStr+"_Created", timeNow.Format("2006-01-02 15:04:05"),
		guidStr+"_Updated", timeNow.Format("2006-01-02 15:04:05"))

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
	if err != nil {
		return nil, err
	}

	for _, replyKey := range replyKeys {
		replyKeyStr := string(replyKey.([]byte))
		if strings.Contains(replyKeyStr, "_Id") {
			idValue, _ := conn.Do("HGET", "reply", replyKeyStr)
			idValueStr := string(idValue.([]byte))

			reply := new(Reply)
			Id, _ := conn.Do("HGET", "reply", idValueStr+"_Id")
			Tid, _ := conn.Do("HGET", "reply", idValueStr+"_Tid")
			Name, _ := conn.Do("HGET", "reply", idValueStr+"_Name")
			Content, _ := conn.Do("HGET", "reply", idValueStr+"_Content")
			Created, _ := conn.Do("HGET", "reply", idValueStr+"_Created")
			Updated, _ := conn.Do("HGET", "reply", idValueStr+"_Updated")

			reply.Id = int64(Id.(int64))
			reply.Tid = int64(Tid.(int64))
			reply.Name = string(Name.([]byte))
			reply.Content = string(Content.([]byte))
			reply.Created, _ = time.Parse("2006-01-02 15:04:05", string(Created.([]byte)))
			reply.Updated, _ = time.Parse("2006-01-02 15:04:05", string(Updated.([]byte)))
			replies = append(replies, reply)
		}
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
	if err != nil {
		return err
	}

	var replyTime string
	for _, replyKey := range replyKeys {
		replyKeyStr := string(replyKey.([]byte))
		if strings.Contains(replyKeyStr, "_Updated") {
			updatedValue, _ := conn.Do("HGET", "reply", replyKeyStr)
			updatedValueStr := string(updatedValue.([]byte))

			if updatedValueStr > replyTime {
				replyTime = updatedValueStr
			}

		}
	}

	//更新Topic
	conn.Do("HSET", "topic", tid+"_ReplyCount", replyCount)
	conn.Do("HSET", "topic", tid+"_ReplyTime", replyTime)
	return nil

}

///endRegion
