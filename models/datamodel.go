package models

import (
	"time"
)

// 分类
type Category struct {
	Id         int64
	Title      string
	Views      int64 `orm:"index"`
	TopicCount int64
	Created    time.Time `orm:"index"`
	Updated    time.Time `orm:"index"`
}

// 文章
type Topic struct {
	Id         int64
	Title      string
	Category   string
	Lables     string
	Content    string `orm:"size(5000)"`
	Attachment string
	Views      int64 `orm:"index"`
	Author     string
	ReplyTime  time.Time `orm:"index"` //最新评论时间
	ReplyCount int64     //评论数量
	Created    time.Time `orm:"index"`
	Updated    time.Time `orm:"index"`
}

// 评论
type Reply struct {
	Id      int64
	Tid     int64
	Name    string
	Content string    `orm:"size(1000)"`
	Created time.Time `orm:"index"`
	Updated time.Time `orm:"index"`
}

// 用户
type User struct {
	Id       int64
	Name     string
	Password string
	Created  time.Time `orm:"index"`
	Updated  time.Time `orm:"index"`
}
