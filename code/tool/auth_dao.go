package tool

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

/*
	用于登录权限检验，连接数据库等服务
*/

func init() {
	if exist, _ := FileExist("databases/login.db"); !exist {
		initDatabase()
	}
}

//初始化 数据库 ： login.db
func initDatabase() {
	db, err := sql.Open("sqlite3", "databases/login.db")
	if err != nil {
		log.Panicln("[login.db]:  数据库初始化失败!")
		return
	}
	defer db.Close()

	sql_table := `CREATE TABLE IF NOT EXISTS "userinfo" (
		"uid" INTEGER PRIMARY KEY AUTOINCREMENT,
		"username" VARCHAR(64) NULL,
		"password" VARCHAR(64) NULL
	)`
	db.Exec(sql_table)

	stmt, err := db.Prepare("INSERT INTO userinfo(username, password) values(?, ?)")
	if err != nil {
		log.Println("[login.db]:  数据库插入数据失败!")
		return
	}

	//插入4条数据，密码以加密方式存入数据库
	if _, err = stmt.Exec("maple", Encryption("maple")); err != nil {
		log.Println("[login.db]:  数据库插入数据失败!")
		return
	}
	if _, err = stmt.Exec("syrup", Encryption("syrup")); err != nil {
		log.Println("[login.db]:  数据库插入数据失败!")
		return
	}
	if _, err = stmt.Exec("tang", Encryption("tang")); err != nil {
		log.Println("[login.db]:  数据库插入数据失败!")
		return
	}
	if _, err = stmt.Exec("zheng", Encryption("zheng")); err != nil {
		log.Println("[login.db]:  数据库插入数据失败!")
		return
	}
}

//检验用户名、密码是否正确
//用于检验用户名和密码是否存在于数据库
func UserLoginValidate(username, password string) bool {
	db, err := sql.Open("sqlite3", "databases/login.db")
	if err != nil {
		log.Panicln("[login.db]:  数据库初始化失败!")
		return false
	}
	defer db.Close()

	rows, err := db.Query("SELECT username, password FROM userinfo")
	if err != nil {
		log.Panicln("[login.db]:  数据库查询数据失败!")
		return false
	}

	var (
		//uid int
		uname string
		pword string
	)
	for rows.Next() {
		err = rows.Scan(&uname, &pword)
		if err != nil {
			fmt.Println(err)
			break
		}
		//fmt.Println(uname, pword)
		if username == uname && Encryption(password) == pword {
			return true
		}
	}
	return false
}
