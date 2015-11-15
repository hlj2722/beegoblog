package models

import (
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// 设置数据库路径
	_DB_NAME = "data/beeblog.db"
	// 设置数据库名称
	_SQLITE3_DRIVER = "sqlite3"
)

func RegisterDB() {
	// 检查数据库文件
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}

	// 注册模型
	orm.RegisterModel(new(Category), new(Topic), new(Reply), new(User))
	// 注册驱动（“sqlite3” 属于默认注册，此处代码可省略）
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DR_Sqlite)
	// 注册默认数据库
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)
}

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
	beego.Alert("test------------------")
	if isListAll {
		beego.Alert("test------------------")
		_, err = qs.All(&cates) //过滤得到文章数量TopicCount大于0的分类

	} else {
		_, err = qs.Filter("TopicCount__gt", 0).All(&cates) //过滤得到文章数量TopicCount大于0的分类
	}
	beego.Alert(cates)
	return cates, err
}

///endReigon

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

///region Reply
func AddReply(tid, nickname, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	reply := &Reply{
		Tid:     tidNum,
		Name:    nickname,
		Content: content,
		Created: time.Now(),
		Updated: time.Now(),
	}
	o := orm.NewOrm()
	_, err = o.Insert(reply)
	if err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyTime = time.Now()
		topic.ReplyCount++
		_, err = o.Update(topic)
	}
	return err
}

func GetAllReplies(tid string) (replies []*Reply, err error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	replies = make([]*Reply, 0)

	o := orm.NewOrm()
	qs := o.QueryTable("Reply")
	_, err = qs.Filter("tid", tidNum).All(&replies)
	return replies, err
}

func DeleteReply(rid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	var tidNum int64
	reply := &Reply{Id: ridNum}
	if o.Read(reply) == nil {
		tidNum = reply.Tid
		_, err = o.Delete(reply)
		if err != nil {
			return err
		}
	}

	replies := make([]*Reply, 0)
	qs := o.QueryTable("Reply")
	_, err = qs.Filter("tid", tidNum).OrderBy("-created").All(&replies)
	if err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyTime = replies[0].Created
		topic.ReplyCount = int64(len(replies))
		_, err = o.Update(topic)
	}
	return err
}

///endRegion

///region User
func AddUser(uname, pwd string) error {
	o := orm.NewOrm()
	user := &User{
		Name:     uname,
		Password: pwd,
		Created:  time.Now(),
		Updated:  time.Now(),
	}

	_, err := o.Insert(user)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(uname string) error {
	o := orm.NewOrm()

	user := &User{Name: uname}
	if o.Read(user) == nil {
		_, err := o.Delete(user)
		if err != nil {
			return err
		}
	}
	return nil
}

func ModifyUser(uname, pwd string) error {
	o := orm.NewOrm()
	user := &User{Name: uname}
	if o.Read(user) == nil {
		user.Password = pwd
		user.Updated = time.Now()
		_, err := o.Update(user)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetUser(uname string) (*User, error) {

	o := orm.NewOrm()

	user := new(User)

	qs := o.QueryTable("user")
	err := qs.Filter("name", uname).One(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetAllUsers(isDesc bool) (users []*User, err error) {

	o := orm.NewOrm()

	users = make([]*User, 0)
	qs := o.QueryTable("user")
	if isDesc {
		_, err = qs.OrderBy("-created").All(&users)

	} else {
		_, err = qs.All(&users)
	}
	return users, err
}

///endregion
