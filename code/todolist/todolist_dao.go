package todolist

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
	"tweb/code/tool"
)

/*
	处理 todolist 相关服务器接口
*/
var todolist_Db_RWLock sync.RWMutex

func init() {
	if exist, _ := tool.FileExist("databases/todolist.db"); !exist {
		initDatabase()
	}
	go checkHistoryTasksTimeout(time.Hour * 6)
	//addActiveTask("今天很开心#^_^#")
	//addActiveTask("联盟连跪10盘，卸载卸载!!!#别了英雄联盟#")
	//addActiveTask("编程很有意思.#golang#我也想学#Vue#")
}

//初始化服务器 todolist.db
//todolist.db 只有一个表 works
func initDatabase() {
	todolist_Db_RWLock.Lock()
	defer todolist_Db_RWLock.Unlock()
	db, err := sql.Open("sqlite3", "databaseodolist.db")
	if err != nil {
		log.Panicln("[todolist.db]:  数据库初始化失败!")
		return
	}
	defer db.Close()

	sql_table := `CREATE TABLE IF NOT EXISTS "tasks" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"task" TEXT,
		"date" VARCHAR(10) NOT NULL,
		"time" VARCHAR(8) NOT NULL,
		"editstatus" VARCHAR(2) NOT NULL,
		"done" INT NOT NULL,
		"tag" VARCHAR(64) NOT NULL,
		"deleted" INT NOT NULL,
		"deletetime" VARCHAR(10) NOT NULL,
		"lefttime" INT NOT NULL
	)`
	db.Exec(sql_table)
}

//添加一个新的todo项目/任务
func addActiveTask(data string) (*TaskModel, bool) {
	todolist_Db_RWLock.Lock()
	defer todolist_Db_RWLock.Unlock()
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	if err != nil {
		log.Panicln("[todolist.db add new thing]:  数据库初始化失败!")
		return nil, false
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO tasks (task, date, time, editstatus, done, tag, deleted, deletetime, lefttime) values(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("[todolist.db add new thing]:  准备 插入数据失败!")
		return nil, false
	}
	//插入1条数据
	date, time := tool.TimeFormat(time.Now())
	task, tag := dataHandler(data)
	if _, err = stmt.Exec(task, date, time, "发布", 0, tag, 0, "00-00-00", 366); err != nil {
		log.Println("[todolist.db add new thing]:  数据库插入数据失败!")
		fmt.Println(err)
		return nil, false
	}
	tm := new(TaskModel)
	row := db.QueryRow("SELECT * FROM tasks WHERE id = (SELECT MAX(id) FROM tasks)")
	err = row.Scan(&tm.Id, &tm.Task, &tm.Date, &tm.Time, &tm.Editstatus, &tm.Done, &tm.Tag, &tm.Deleted, &tm.Deletetime, &tm.Lefttime)
	if err != nil {
		return nil, false
	}
	return tm, true
}

//更新活跃任务的内容
func updateActiveTaskContent(id int, task string, tag string) bool {
	todolist_Db_RWLock.Lock()
	defer todolist_Db_RWLock.Unlock()
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	if err != nil {
		log.Panicln("[todolist.db update]:  数据库初始化失败!")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE tasks SET task = ?, date = ?, time = ?, editstatus = ?, tag = ? WHERE id = ?")
	if err != nil {
		log.Println("[todolist.db update]:  准备更新数据失败!")
		return false
	}
	//更新1条数据
	date, time := tool.TimeFormat(time.Now())
	if _, err = stmt.Exec(task, date, time, "编辑", tag, id); err != nil {
		log.Println("[todolist.db update]:  数据库更新数据失败!")
		return false
	}
	return true
}

//更新活跃任务的完成状态
func updateActiveTaskStatus(id, done int) bool {
	todolist_Db_RWLock.Lock()
	defer todolist_Db_RWLock.Unlock()
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	if err != nil {
		log.Panicln("[todolist.db update]:  数据库初始化失败!")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE tasks SET done = ? WHERE id = ?")
	if err != nil {
		log.Println("[todolist.db update]:  准备更新数据失败!")
		return false
	}
	//更新1条数据
	if _, err = stmt.Exec(done, id); err != nil {
		log.Println("[todolist.db update]:  数据库更新数据失败!")
		return false
	}
	return true
}

//删除一条活跃任务 的记录，其实就是把 deleted 状态改为 true，并没有删除原记录
func deleteActiveTask(id int) bool {
	todolist_Db_RWLock.Lock()
	defer todolist_Db_RWLock.Unlock()
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	if err != nil {
		log.Panicln("[todolist.db delete current thing]:  数据库初始化失败!")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE tasks SET deleted = ?, deletetime = ?, lefttime = ? WHERE id = ?")
	if err != nil {
		log.Println("[todolist.db delete current thing]:  准备删除数据失败!")
		return false
	}
	//更新1条数据
	date, _ := tool.TimeFormat(time.Now())
	if _, err = stmt.Exec(1, date, 366, id); err != nil {
		log.Println("[todolist.db elete current thing]:  数据库删除数据失败!")
		return false
	}
	return true
}

//删除一条历史任务，也就是垃圾箱里的任务
//永久删除
func deleteHistoryTask(id int) bool {
	todolist_Db_RWLock.Lock()
	defer todolist_Db_RWLock.Unlock()
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	if err != nil {
		log.Panicln("[todolist.db delete history thing]:  数据库初始化失败!")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM tasks where id = ?")
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
func recoverHistoryTask(id int) bool {
	todolist_Db_RWLock.Lock()
	defer todolist_Db_RWLock.Unlock()
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	if err != nil {
		log.Panicln("[todolist.db add history thing]:  数据库初始化失败!")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE tasks SET deleted = ? WHERE id = ?")
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
	var task, tag string
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
				task += c
			}
		}

	}
	return task, tag
}

//请求数据库中所有的数据
func queryAllTasks() ([]TaskModel, error) {
	todolist_Db_RWLock.RLock()
	defer todolist_Db_RWLock.RUnlock()
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	var res []TaskModel
	if err != nil {
		log.Panicln("[todolist.db add history thing]:  数据库初始化失败!")
		return res, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM tasks ORDER BY id DESC")
	if err != nil {
		log.Println("[todolist.db add history thing]:  数据库删除数据失败!222")
		return res, err
	}
	for rows.Next() {
		var tm TaskModel
		err = rows.Scan(&tm.Id, &tm.Task, &tm.Date, &tm.Time, &tm.Editstatus, &tm.Done, &tm.Tag, &tm.Deleted, &tm.Deletetime, &tm.Lefttime)
		if err != nil {
			log.Println("[todolist.db add history thing]:  数据库删除数据失败!222")
			return make([]TaskModel, 0), err
		}
		res = append(res, tm)
	}
	return res, nil
}

//请求数据库中的活跃任务
func queryActiveTasks() ([]TaskModel, error) {
	todolist_Db_RWLock.RLock()
	defer todolist_Db_RWLock.RUnlock()
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	var res []TaskModel
	if err != nil {
		log.Panicln("[todolist.db add history thing]:  数据库初始化失败!")
		return res, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM tasks WHERE deleted = 0")
	if err != nil {
		log.Println("[todolist.db add history thing]:  数据库删除数据失败!222")
		return res, err
	}

	for rows.Next() {
		var tm TaskModel
		err = rows.Scan(&tm.Id, &tm.Task, &tm.Date, &tm.Time, &tm.Editstatus, &tm.Done, &tm.Tag, &tm.Deleted, &tm.Deletetime, &tm.Lefttime)
		if err != nil {
			log.Println("[todolist.db add history thing]:  数据库删除数据失败!222")
			return make([]TaskModel, 0), err
		}
		res = append(res, tm)
	}
	return res, nil
}

//请求数据库中的回收站任务
func queryHistoryTasks() ([]TaskModel, error) {
	todolist_Db_RWLock.RLock()
	defer todolist_Db_RWLock.RUnlock()
	db, err := sql.Open("sqlite3", "databases/todolist.db")
	var res []TaskModel
	if err != nil {
		log.Panicln("[todolist.db add history thing]:  数据库初始化失败!")
		return res, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM tasks WHERE deleted = 1")
	if err != nil {
		log.Println("[todolist.db add history thing]:  数据库删除数据失败!")
		return res, err
	}

	for rows.Next() {
		var tm TaskModel
		err = rows.Scan(&tm.Id, &tm.Task, &tm.Date, &tm.Time, &tm.Editstatus, &tm.Done, &tm.Tag, &tm.Deleted, &tm.Deletetime, &tm.Lefttime)
		if err != nil {
			log.Println("[todolist.db add history thing]:  数据库删除数据失败!222")
			return make([]TaskModel, 0), err
		}
		res = append(res, tm)
	}
	return res, nil
}

func checkHistoryTasksTimeout(duration time.Duration) {
	db, _ := sql.Open("sqlite3", "databases/todolist.db")
	defer db.Close()
	for {
		todolist_Db_RWLock.Lock()
		_, err := db.Exec("update tasks set lefttime = (deletetime + 365 - date('now')) WHERE deleted = 1")
		if err != nil {
			log.Println("[check history tasks timeout] update lefttime error!")
		}
		_, err = db.Exec("DELETE FROM tasks WHERE deleted = 1 AND lefttime < 0")
		if err != nil {
			log.Println("[check history tasks timeout] delete outtime history tasks error!")
		}
		todolist_Db_RWLock.Unlock()
		time.Sleep(duration)
	}
}
