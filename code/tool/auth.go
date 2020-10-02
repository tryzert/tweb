package tool

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type sessinfo struct {
	userPassKey string
	deadline time.Time
}

type UserPassKeySessions map[string]sessinfo

var Upks UserPassKeySessions

func (upks UserPassKeySessions) Get(username string) string {
	res, ok := upks[username]
	if ok {
		return res.userPassKey
	}
	return ""
}


func (upks UserPassKeySessions) Add(username string, duration time.Duration) {
	upks.Set(username, CreateRandPassKey(32), duration)
}
func (upks UserPassKeySessions) Set(username, userPassKey string, duration time.Duration) {
	upks[username] = sessinfo{
		userPassKey: userPassKey,
		deadline: time.Now().Add(duration),
	}
}

func (upks UserPassKeySessions) Delete(username string) {
	delete(upks, username)
}


func (upks UserPassKeySessions) Clear() {
	for username, _ := range upks {
		delete(upks, username)
	}
}


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


func AuthMiddleWare(c *gin.Context) {
	sess := sessions.Default(c)
	username, _ := sess.Get("username").(string)
	userPassKey, _ := sess.Get("userPassKey").(string)
	//fmt.Println("username:", username, "userPassKey:", userPassKey, "upks:", Upks)
	if username == "maple" && userPassKey == Upks.Get(username){
		c.Next()
	} else {
		c.Redirect(http.StatusFound, "/login")
		//c.Abort()
		return
	}
}


func UserLoginValidate(username, password string) bool {
	if username == "maple" && password == "maple" {
		return true
	}
	return false
}