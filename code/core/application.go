package core

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
	"tweb/code/admin"
	"tweb/code/home"
	"tweb/code/login"
	"tweb/code/logout"
	"tweb/code/notFound404"
	"tweb/code/tapbag"
	"tweb/code/todolist"
)

var SERVER_IS_RUNNING bool
var SERVER_ADMIN_SIGNAL chan int

func init() {
	SERVER_IS_RUNNING = false
	SERVER_ADMIN_SIGNAL = make(chan int, 1)
}

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
	if st.Services.TodolistConf.Run {
		htmlFiles = append(htmlFiles, "./template/todolist/todolist_index.html")
	}

	if st.Services.TapbagConf.Run {
		htmlFiles = append(htmlFiles, "./template/tapbag/tapbag_index.html")
	}

	r.LoadHTMLFiles(htmlFiles...)

	/**
	 * 以下为默认服务，必须开启
	 */
	//注册登录页login服务
	login.RegisterService(r)
	//logout服务
	logout.RegisterService(r)
	//注册主页home服务
	home.RegisterService(r)
	//注册管理页admin服务
	admin.RegisterService(r, SERVER_ADMIN_SIGNAL)
	//注册404页面
	notFound404.RegisterService(r)

	/*
	** 以下为注册服务
	 */
	//注册todolist服务
	if st.Services.TodolistConf.Run {
		todolist.RegisterService(r)
	}
	//注册tim服务
	//if st.Services.Run_tim {
	//	tim.RegisterService(r)
	//}
	//注册TapBag服务
	if st.Services.TapbagConf.Run {
		tapbag.RegisterService(r, st.Services.TapbagConf.SrcPath)
	}
}

func run() *http.Server {
	gin.SetMode(gin.ReleaseMode)
	//把user这个接头体注册进来，后面跨路由才可以获取到user数据
	//gob.Register(tool.Upks)
	r := gin.Default()
	//r := gin.New()

	//服务器端设置允许前端跨域请求
	r.Use(Cors())

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("twebLoginValidate", store))

	//读取配置文件
	st := getSettings()
	registerDefaultService(r, st)

	r.StaticFile("/favicon.ico", "public/favicon.ico")
	r.StaticFS("/public", http.Dir("public"))

	//文件下载权限鉴定
	//r.GET("/public/file", func(c *gin.Context) {
	//	c.File("public/vue/vue.js")
	//})
	fmt.Printf("Server running on: http://localhost%v\n", st.Port)
	//_ = r.Run(st.Port)
	server := &http.Server{Addr: st.Port, Handler: r}
	SERVER_IS_RUNNING = true
	go server.ListenAndServe()
	return server
}

//程序运行入口
func Run() {
	ctx, _ := context.WithCancel(context.Background())
	SERVER_ADMIN_SIGNAL <- 100
	var server *http.Server
	//server.Shutdown(ctx)
	for {
		select {
		case signal := <-SERVER_ADMIN_SIGNAL:
			switch signal {
			case 100: // start
				if !SERVER_IS_RUNNING {
					SERVER_IS_RUNNING = true
					server = run()
				}
			case 101: // shutdown
				if SERVER_IS_RUNNING {
					server.Shutdown(ctx)
					SERVER_IS_RUNNING = false
					os.Exit(0)
				}
			case 102: //reboot
				if SERVER_IS_RUNNING {
					server.Shutdown(ctx)
					SERVER_IS_RUNNING = false
					SERVER_ADMIN_SIGNAL <- 100
				}
			}
		case <-time.After(time.Second * 10):
			//fmt.Println("no signal ...")
		}
	}
}
