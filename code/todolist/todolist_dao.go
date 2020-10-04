package todolist

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"tweb/code/tool"
)

/*
	处理 todolist 相关服务器接口
*/

func init() {
	if exist, _ := tool.FileExist("databases/todolist.db"); !exist {
		initDatabase()
	}
	//addCurrentThing("吃饭#休息##急😂#")
	//updateCurrentThing(3, "吃饭睡觉打豆豆，明天还有事要早起", 1, "#好好休息#")
	deleteCurrentThing(2)
	//deleteHistoryThing(1)
	//recoverHistoryThing(2)
	fmt.Println(queryAll())
	fmt.Println(queryCurrentThings())
	fmt.Println(queryHistoryThings())
}

//初始化服务器 todolist.db
//todolist.db 只有一个表 works
func initDatabase() {
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	if err != nil {
		log.Panicln("[todolist.db]:  数据库初始化失败!")
		return
	}
	defer db.Close()

	sql_table := `CREATE TABLE IF NOT EXISTS "works" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"thing" TEXT,
		"date" VARCHAR(10) NOT NULL,
		"time" VARCHAR(8) NOT NULL,
		"editstatus" VARCHAR(2) NOT NULL,
		"done" INT NOT NULL,
		"tag" VARCHAR(64) NOT NULL,
		"isdeleted" INT NOT NULL,
		"deletetime" VARCHAR(10) NOT NULL,
		"lefttime" INT NOT NULL
	)`
	db.Exec(sql_table)
}

//添加一个新的todo项目/任务
func addCurrentThing(data string) bool {
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	if err != nil {
		log.Panicln("[todolist.db add new thing]:  数据库初始化失败!")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO works (thing, date, time, editstatus, done, tag, isdeleted, deletetime, lefttime) values(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("[todolist.db add new thing]:  准备 插入数据失败!")
		return false
	}
	//插入1条数据
	date, time := tool.TimeFormat(time.Now())
	thing, tag := dataHandler(data)
	if _, err = stmt.Exec(thing, date, time, "发布", 0, tag, 0, "", -1); err != nil {
		log.Println("[todolist.db add new thing]:  数据库插入数据失败!")
		fmt.Println(err)
		return false
	}
	return true
}

//更新活跃任务的状态
func updateCurrentThing(id int, thing string, done int, tag string) bool {
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	if err != nil {
		log.Panicln("[todolist.db update]:  数据库初始化失败!")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE works SET thing = ?, date = ?, time = ?, editstatus = ?, done = ?, tag = ? WHERE id = ?")
	if err != nil {
		log.Println("[todolist.db update]:  准备更新数据失败!")
		return false
	}
	//更新1条数据
	date, time := tool.TimeFormat(time.Now())
	if _, err = stmt.Exec(thing, date, time, "编辑", done, tag, id); err != nil {
		log.Println("[todolist.db update]:  数据库更新数据失败!")
		return false
	}
	return true
}

//删除一条活跃任务 的记录，其实就是把 deleted 状态改为 true，并没有删除原记录
func deleteCurrentThing(id int) bool {
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	if err != nil {
		log.Panicln("[todolist.db delete current thing]:  数据库初始化失败!")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE works SET isdeleted = ?, deletetime = ? WHERE id = ?")
	if err != nil {
		log.Println("[todolist.db delete current thing]:  准备删除数据失败!")
		return false
	}
	//更新1条数据
	date, _ := tool.TimeFormat(time.Now())
	if _, err = stmt.Exec(1, date, id); err != nil {
		log.Println("[todolist.db elete current thing]:  数据库删除数据失败!")
		return false
	}
	return false
}

//删除一条历史任务，也就是垃圾箱里的任务
//永久删除
func deleteHistoryThing(id int) bool {
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	if err != nil {
		log.Panicln("[todolist.db delete history thing]:  数据库初始化失败!")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM works where id = ?")
	if err != nil {
		log.Println("[todolist.db delete history thing]:  准备删除数据失败!")
		return false
	}
	//插入1条数据
	if _, err = stmt.Exec(id); err != nil {
		log.Println("[todolist.db add history thing]:  数据库删除数据失败!")
		return false
	}
	return true
}

//将回收站里的历史任务恢复到当前活跃任务中
func recoverHistoryThing(id int) bool {
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	if err != nil {
		log.Panicln("[todolist.db add history thing]:  数据库初始化失败!")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE works SET isdeleted = ? WHERE id = ?")
	if err != nil {
		log.Println("[todolist.db add history thing]:  准备恢复数据失败!")
		return false
	}
	//插入1条数据
	if _, err = stmt.Exec(0, id); err != nil {
		log.Println("[todolist.db add history thing]:  数据库恢复数据失败!")
		return false
	}
	return true
}

//用于格式化前端过来的数据， 把文本信息转化为content和tag
func dataHandler(data string) (string, string) {
	var thing, tag string
	ist := false
	for _, r := range data {
		c := string(r)
		if c == "#" {
			ist = !ist
			tag += c
		} else {
			if ist {
				tag += c
			} else {
				thing += c
			}
		}

	}
	return thing, tag
}

//请求数据库中所有的数据
func queryAll() ([]ThingModel, error) {
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	var res []ThingModel
	if err != nil {
		log.Panicln("[todolist.db add history thing]:  数据库初始化失败!")
		return res, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM works")
	if err != nil {
		log.Println("[todolist.db add history thing]:  数据库删除数据失败!222")
		return res, err
	}

	for rows.Next() {
		var tm ThingModel
		err = rows.Scan(&tm.Id, &tm.Thing, &tm.Date, &tm.Time, &tm.Editstatus, &tm.Done, &tm.Tag, &tm.Isdeleted, &tm.Deletetime, &tm.Lefttime)
		if err != nil {
			log.Println("[todolist.db add history thing]:  数据库删除数据失败!222")
			return make([]ThingModel, 0), err
		}
		res = append(res, tm)
	}
	return res, nil
}

//请求数据库中的活跃任务
func queryCurrentThings() ([]ThingModel, error) {
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	var res []ThingModel
	if err != nil {
		log.Panicln("[todolist.db add history thing]:  数据库初始化失败!")
		return res, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM works WHERE isdeleted = ?", 0)
	if err != nil {
		log.Println("[todolist.db add history thing]:  数据库删除数据失败!222")
		return res, err
	}

	for rows.Next() {
		var tm ThingModel
		err = rows.Scan(&tm.Id, &tm.Thing, &tm.Date, &tm.Time, &tm.Editstatus, &tm.Done, &tm.Tag, &tm.Isdeleted, &tm.Deletetime, &tm.Lefttime)
		if err != nil {
			log.Println("[todolist.db add history thing]:  数据库删除数据失败!222")
			return make([]ThingModel, 0), err
		}
		res = append(res, tm)
	}
	return res, nil
}

//请求数据库中的回收站任务
func queryHistoryThings() ([]ThingModel, error) {
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	var res []ThingModel
	if err != nil {
		log.Panicln("[todolist.db add history thing]:  数据库初始化失败!")
		return res, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM works WHERE isdeleted = ?", 1)
	if err != nil {
		log.Println("[todolist.db add history thing]:  数据库删除数据失败!222")
		return res, err
	}

	for rows.Next() {
		var tm ThingModel
		err = rows.Scan(&tm.Id, &tm.Thing, &tm.Date, &tm.Time, &tm.Editstatus, &tm.Done, &tm.Tag, &tm.Isdeleted, &tm.Deletetime, &tm.Lefttime)
		if err != nil {
			log.Println("[todolist.db add history thing]:  数据库删除数据失败!222")
			return make([]ThingModel, 0), err
		}
		res = append(res, tm)
	}
	return res, nil
}
