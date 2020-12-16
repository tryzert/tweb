package tool

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
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
	deadline    time.Time
}

//type UserPassKeySessions map[string]sessinfo
//var Upks UserPassKeySessions
type UserPassKeySessionsPool struct {
	pool map[string]*sessinfo
	sync.RWMutex
}

var Upks *UserPassKeySessionsPool

//通过用户名，取出session中对应的 userPassKey
func (upks *UserPassKeySessionsPool) Get(username string) string {
	upks.RWMutex.RLock()
	defer upks.RWMutex.RUnlock()
	res, ok := upks.pool[username]
	if ok {
		return res.userPassKey
	}
	return ""
}

//有新用户接入，添加一个新用户session
func (upks *UserPassKeySessionsPool) Add(username string, duration time.Duration) {
	upks.RWMutex.Lock()

	if ssinfo, ok := upks.pool[username]; ok {
		ssinfo.deadline = time.Now().Add(duration)
		upks.RWMutex.Unlock()
	} else {
		upks.RWMutex.Unlock()
		upks.Set(username, CreateRandPassKey(32), duration)
	}

}

//设置或更新一个用户session信息
func (upks *UserPassKeySessionsPool) Set(username, userPassKey string, duration time.Duration) {
	upks.RWMutex.Lock()
	defer upks.RWMutex.Unlock()
	upks.pool[username] = &sessinfo{
		userPassKey: userPassKey,
		deadline:    time.Now().Add(duration),
	}
}

//删除一个用户session信息
func (upks *UserPassKeySessionsPool) Delete(username string) {
	upks.RWMutex.Lock()
	defer upks.RWMutex.Unlock()
	delete(upks.pool, username)
}

//清空此服务器上所有用户的session信息
func (upks *UserPassKeySessionsPool) Clear() {
	upks.RWMutex.Lock()
	defer upks.RWMutex.Unlock()
	for username, _ := range upks.pool {
		//delete(upks.pool, username)
		upks.Delete(username)
	}
}

//用于自动更新服务器session状态，定时清除过期的session信息
func (upks *UserPassKeySessionsPool) CheckTimeout(duration time.Duration) {
	for {
		now := time.Now()
		for username, info := range upks.pool {
			if info.deadline.Before(now) {
				upks.Delete(username)
			}
		}
		time.Sleep(duration)
	}
}

func init() {
	Upks = &UserPassKeySessionsPool{
		pool: make(map[string]*sessinfo),
	}
	go Upks.CheckTimeout(time.Hour * 2)
}

//用于判断是否登录的中间件
func AuthMiddleWare(c *gin.Context) {
	sess := sessions.Default(c)
	username, _ := sess.Get("username").(string)
	userPassKey, _ := sess.Get("userPassKey").(string)
	if username != "" && userPassKey == Upks.Get(username) {
		c.Next()
		return
	} else {
		c.Redirect(http.StatusFound, "/login")
		return
	}
}
