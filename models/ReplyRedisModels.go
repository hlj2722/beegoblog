package models

import (
	_ "github.com/hopehook/beegoblog/tools"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"strconv"
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
	guidStr := strconv.FormatInt(int64(guid.(int64)), 10)
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
		if strings.Contains(replyKeyStr, "_Tid") {
			tidValue, _ := conn.Do("HGET", "reply", replyKeyStr)
			tidValueStr := string(tidValue.([]byte))
			if tidValueStr == tid {
				idValueStr := strings.TrimRight(replyKeyStr, "_Tid")

				reply := new(Reply)
				Id, _ := conn.Do("HGET", "reply", idValueStr+"_Id")
				Tid, _ := conn.Do("HGET", "reply", idValueStr+"_Tid")
				Name, _ := conn.Do("HGET", "reply", idValueStr+"_Name")
				Content, _ := conn.Do("HGET", "reply", idValueStr+"_Content")
				Created, _ := conn.Do("HGET", "reply", idValueStr+"_Created")
				Updated, _ := conn.Do("HGET", "reply", idValueStr+"_Updated")

				reply.Id, _ = strconv.ParseInt(string(Id.([]byte)), 10, 0)
				reply.Tid, _ = strconv.ParseInt(string(Tid.([]byte)), 10, 0)
				reply.Name = string(Name.([]byte))
				reply.Content = string(Content.([]byte))
				reply.Created, _ = time.Parse("2006-01-02 15:04:05", string(Created.([]byte)))
				reply.Updated, _ = time.Parse("2006-01-02 15:04:05", string(Updated.([]byte)))
				replies = append(replies, reply)
			}

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
	tid, _ := conn.Do("HGET", "reply", rid+"_Tid")
	tidStr := string(tid.([]byte))

	//删除Reply
	conn.Do("HDEL", "reply", rid+"_Id")
	conn.Do("HDEL", "reply", rid+"_Tid")
	conn.Do("HDEL", "reply", rid+"_Name")
	conn.Do("HDEL", "reply", rid+"_Content")
	conn.Do("HDEL", "reply", rid+"_Created")
	conn.Do("HDEL", "reply", rid+"_Updated")

	//获取评论数量
	count, _ := conn.Do("HLEN", "reply")
	replyCount := int64(count.(int64)) / 6

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
	conn.Do("HSET", "topic", tidStr+"_ReplyCount", replyCount)
	conn.Do("HSET", "topic", tidStr+"_ReplyTime", replyTime)
	return nil

}

///endRegion
