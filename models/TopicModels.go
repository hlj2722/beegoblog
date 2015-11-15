package models

import (
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

///region Topic
func AddTopic(title, category, lable, content, attachment, author string) error {
	// 处理标签
	lable = "$" + strings.Join(strings.Split(lable, " "), "#$") + "#"

	o := orm.NewOrm()

	topic := &Topic{
		Title:      title,
		Category:   category,
		Lables:     lable,
		Content:    content,
		Attachment: attachment,
		Author:     author,
		ReplyTime:  time.Now(),
		Created:    time.Now(),
		Updated:    time.Now(),
	}
	_, err := o.Insert(topic)
	if err != nil {
		return err
	}

	// 更新分类统计
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		if cate.TopicCount > 0 {
			cate.TopicCount++
		} else {
			cate.TopicCount = 1
		}

		_, err = o.Update(cate)

	} else {
		cate.Title = category
		cate.Created = time.Now()
		cate.Updated = time.Now()
		cate.TopicCount = 1
		_, err = o.Insert(cate)
	}

	return err
}

func GetTopic(tid string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	o := orm.NewOrm()

	topic := new(Topic)

	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return nil, err
	}

	topic.Views++
	_, err = o.Update(topic)

	topic.Lables =
		strings.Replace(strings.Replace(topic.Lables, "#", " ", -1), "$", "", -1)

	return topic, nil
}

func ModifyTopic(tid, title, category, lable, content, attachment string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	lable = "$" + strings.Join(strings.Split(lable, " "), "#$") + "#"

	var oldCate, oldAttach string
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		oldAttach = topic.Attachment
		topic.Title = title
		topic.Category = category
		topic.Lables = lable
		topic.Content = content
		topic.Attachment = attachment
		topic.Updated = time.Now()
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}

	// 更新分类统计
	if oldCate != category {

		cate1 := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", oldCate).One(cate1)
		if err == nil {
			cate1.TopicCount--
			_, err = o.Update(cate1)
		}

		cate2 := new(Category)
		err = qs.Filter("title", category).One(cate2)
		if err == nil {
			cate2.TopicCount++
			cate2.Updated = time.Now()
			_, err = o.Update(cate2)
		} else {
			cate2.Title = category
			cate2.Created = time.Now()
			cate2.Updated = time.Now()
			cate2.TopicCount = 1
			_, err = o.Insert(cate2)
		}

	}

	// 删除旧的附件
	if oldAttach != attachment {
		os.Remove(path.Join("attachment", oldAttach))
	}

	return nil
}

func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	var oldCate string
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		_, err = o.Delete(topic)
		if err != nil {
			return err
		}
	}

	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", oldCate).One(cate)
	if err == nil {
		cate.TopicCount--
		_, err = o.Update(cate)
	}

	return err
}

func GetAllTopics(category, lable string, isDesc bool) (topics []*Topic, err error) {
	o := orm.NewOrm()

	topics = make([]*Topic, 0)

	qs := o.QueryTable("topic")
	if isDesc {
		if len(category) > 0 {
			qs = qs.Filter("category", category)
		}
		if len(lable) > 0 {
			lable = strings.Trim(lable, " ")
			qs = qs.Filter("lables__contains", "$"+lable+"#")
		}
		_, err = qs.OrderBy("-created").All(&topics)

	} else {
		_, err = qs.All(&topics)
	}
	return topics, err
}

///endRegion
