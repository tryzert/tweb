package core

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"tweb/code/admin"
	"tweb/code/home"
	"tweb/code/login"
	"tweb/code/notFound404"
	"tweb/code/tim"
	"tweb/code/todolist"
)


//注册服务
func registerDefaultService(r *gin.Engine, st *Settings) {
	//载入html
	htmlFiles := []string{
		"./template/home/home_index.html",
		"./template/login/login_index.html",
		"./template/admin/admin_index.html",
		"./template/notFound404/notFound404_index.html",
	}

	//如果配置文件设置开启服务，则打开服务
	if st.Services.Run_todolist {
		htmlFiles = append(htmlFiles, "./template/todolist/todolist_index.html")
	}
	if st.Services.Run_todolist {
		htmlFiles = append(htmlFiles, "./template/tim/tim_index.html")
	}
	r.LoadHTMLFiles(htmlFiles...)


	/**
	 * 以下为默认服务，必须开启
	*/
	//注册登录页login服务
	login.RegisterSevice(r)
	//注册主页home服务
	home.RegisterSevice(r)
	//注册管理页admin服务
	admin.RegisterSevice(r)
	//注册404页面
	notFound404.RegisterSevice(r)


	/*
	** 以下为注册服务
	*/
	//注册todolist服务
	if st.Services.Run_todolist {
		todolist.RegisterSevice(r)
	}
	//注册tim服务
	if st.Services.Run_tim {
		tim.RegisterSevice(r)
	}
}


type user struct {
	Name string
	PassKey string
}


//程序运行入口
func Run() {
	gin.SetMode(gin.ReleaseMode)
	//把user这个接头体注册进来，后面跨路由才可以获取到user数据
	//gob.Register(tool.Upks)
	r := gin.Default()
	//r := gin.New()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("twebLoginValidate", store))

	//读取配置文件
	st := getSettings()
	registerDefaultService(r, st)

	r.StaticFile("/favicon.ico", "other/favicon.ico")

	fmt.Printf("Server running on: http://localhost%v\n", st.Port)
	_ = r.Run(st.Port)
}
