package tool

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)


/*
	一个简单的基于服务器内存的session, 用于检验登录状态 的 map
	结构：username : {
		userPassKey  一个给定长度的随机字符串
		deadline     session的过期时间
	}
*/


type sessinfo struct {
	userPassKey string
	deadline time.Time
}

type UserPassKeySessions map[string]sessinfo

var Upks UserPassKeySessions


//通过用户名，取出session中对应的 userPassKey
func (upks UserPassKeySessions) Get(username string) string {
	res, ok := upks[username]
	if ok {
		return res.userPassKey
	}
	return ""
}


//有新用户接入，添加一个新用户session
func (upks UserPassKeySessions) Add(username string, duration time.Duration) {
	upks.Set(username, CreateRandPassKey(32), duration)
}

//设置或更新一个用户session信息
func (upks UserPassKeySessions) Set(username, userPassKey string, duration time.Duration) {
	upks[username] = sessinfo{
		userPassKey: userPassKey,
		deadline: time.Now().Add(duration),
	}
}


//删除一个用户session信息
func (upks UserPassKeySessions) Delete(username string) {
	delete(upks, username)
}


//清空此服务器上所有用户的session信息
func (upks UserPassKeySessions) Clear() {
	for username, _ := range upks {
		delete(upks, username)
	}
}


//用于自动更新服务器session状态，定时清除过期的session信息
func (upks UserPassKeySessions) CheckTimeout(duration time.Duration) {
	for {
		now := time.Now()
		for username, info := range upks {
			if info.deadline.Before(now) {
				delete(upks, username)
			}
		}
		time.Sleep(time.Minute)
	}
}


func init() {
	Upks = make(map[string]sessinfo)
	go Upks.CheckTimeout(time.Hour * 2)
}


//用于判断是否登录的中间件
func AuthMiddleWare(c *gin.Context) {
	sess := sessions.Default(c)
	username, _ := sess.Get("username").(string)
	userPassKey, _ := sess.Get("userPassKey").(string)
	//fmt.Println("username:", username, "userPassKey:", userPassKey, "upks:", Upks)
	if userPassKey == Upks.Get(username){
		c.Next()
	} else {
		c.Redirect(http.StatusFound, "/login")
		//c.Abort()
		return
	}
}
