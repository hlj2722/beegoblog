package models

import (
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"strings"
)

///region UserRedis
func AddUserRedis(uname, pwd string) error {
	beego.Alert("================  AddUserRedis(uname, pwd string) error==============")
	//连接 Redis 服务器
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return err
	}

	//调用 AUTH 命令获得授权
	conn.Do("AUTH", beego.AppConfig.String("requirepass")) //配置文件中的Redis密码
	//延迟自动关闭连接
	defer conn.Close()
	beego.Alert("================  AddUserRedis(uname, pwd string) error==============")
	//Redis命令调用
	conn.Do("SET", "uname="+uname, pwd)
	return nil
}

func DeleteUserRedis(uname string) error {
	beego.Alert("================  DeleteUserRedis(uname string) error==============")
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()
	beego.Alert("================  DeleteUserRedis(uname string) error==============")
	conn.Do("DEL", "uname="+uname)
	return nil
}

func ModifyUserRedis(uname, pwd string) error {
	return AddUserRedis(uname, pwd)

}

func GetUserRedis(uname string) (*User, error) {
	beego.Alert("================ GetUserRedis(uname string) (*User, error)==============")
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()
	pwd, err := redis.String(conn.Do("GET", "uname="+uname))
	beego.Alert("================ GetUserRedis(uname string) (*User, error)==============")
	if err != nil {
		return nil, err
	}
	user := &User{
		Name:     uname,
		Password: pwd,
	}

	return user, nil

}

func GetAllUsersRedis(isDesc bool) (users []*User, err error) {
	beego.Alert("================ GetAllUsersRedis(isDesc bool) (users []*User, err error)==============")
	_ = isDesc
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}
	conn.Do("AUTH", beego.AppConfig.String("requirepass"))
	defer conn.Close()
	keys, err := redis.Values(conn.Do("KEYS", "uname=*"))
	if err != nil {
		return nil, err
	}
	beego.Alert("================ GetAllUsersRedis(isDesc bool) (users []*User, err error)==============")
	for _, key := range keys {
		value, err := redis.String(conn.Do("GET", key))
		if err != nil {
			return nil, err
		}
		beego.Alert(strings.TrimLeft(string(key.([]byte)), "uname="))
		user := &User{
			Name:     strings.TrimLeft(string(key.([]byte)), "uname="),
			Password: value,
		}
		users = append(users, user)
	}
	beego.Alert("================ GetAllUsersRedis(isDesc bool) (users []*User, err error)==============")
	return users, nil

}

///endregion
