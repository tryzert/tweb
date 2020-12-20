package tapbag

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
	"time"
	"tweb/code/tool"
)

var DB *sql.DB

func init() {
	exist, err := tool.FileExist("databases/tapbag.db")
	if err != nil {
		panic("数据库 tapbag.db 初始化失败！")
	}
	if !exist {
		initDatabase()
	} else {
		db, err := sql.Open("sqlite3", "databases/tapbag.db")
		if err != nil {
			log.Panicln("数据库 tapbag.db 连接失败！")
			return
		}
		if db.Ping() != nil {
			log.Panicln("数据库 tapbag.db 连接失败！")
			return
		}
		DB = db
	}

	go checkRecycleBinKeepDays(time.Hour * 6)
	go checkShareKeepDays(time.Hour * 6)
}

func initDatabase() {
	db, err := sql.Open("sqlite3", "databases/tapbag.db")
	if err != nil {
		log.Panicln("数据库 tapbag.db 创建失败！")
		return
	}
	if db.Ping() != nil {
		log.Panicln("数据库 tapbag.db 创建失败！")
		return
	}
	sql_table := `CREATE TABLE IF NOT EXISTS "recycleBin" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"origin_rel_path" TEXT,
		"file_name" TEXT,
		"delete_time" VARCHAR(10),
		"deleted_keep_days" INT NOT NULL DEFAULT 366
)`
	db.Exec(sql_table)

	sql_table = `CREATE TABLE IF NOT EXISTS "shares" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"rel_path" TEXT,
		"file_name" TEXT,
		"share_time" VARCHAR(10),
		"share_forever" INT NOT NULL DEFAULT 0,
		"shared_keep_days" INT NOT NULL DEFAULT 366
)`
	db.Exec(sql_table)
	db = DB
}

// return error nums
func addDeleteRecord(relpaths []string) (int, error) {
	stmt, err := DB.Prepare(`INSERT INTO recycleBin (origin_rel_path, file_name, delete_time) VALUES (?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var errCount = 0
	date, _ := tool.TimeFormat(time.Now())
	for _, relpath := range relpaths {
		filename := filepath.Base(relpath)
		_, err = stmt.Exec(relpath, filename, date)
		if err != nil {
			errCount++
		}
	}
	return errCount, err
}

//recover files from recycle bin to their origin dirs
func recoverFiles(ids []int) {

}

func checkRecycleBinKeepDays(duration time.Duration) {
	for {
		rows, err := DB.Query(`SELECT id, file_name, julianday('now') - julianday(delete_time) FROM recycleBin`)
		if err != nil {
			log.Println("database [tapbag.db] check keep days error!")
			time.Sleep(time.Hour)
			continue
		}

		type resItem struct{
			id            int
			pastDays float64
			filename string
		}
		queryRes := []*resItem{}
		for rows.Next() {
			resIt := &resItem{}
			err = rows.Scan(&resIt.id, &resIt.filename, &resIt.pastDays)
			if err != nil {
				log.Println("database [tapbag.db] scan error!")
				continue
			}
			queryRes = append(queryRes, resIt)
		}
		rows.Close()

		stmt, err := DB.Prepare(`UPDATE recycleBin SET deleted_keep_days = ? WHERE id = ?`)
		if err != nil {
			log.Println("database [tapbag.db] prepare update keep days error!")
			time.Sleep(time.Hour)
			continue
		}
		for _, it := range queryRes {
			if _, err = stmt.Exec(366 - int(it.pastDays), it.id); err != nil {
				log.Println("database [tapbag.db] exec update keep days error!")
			}
			if 366 - int(it.pastDays) < 0 {
				// todo: how to get srcPath: define a global variable ?
				os.RemoveAll(filepath.Join(".tweb/recycleBin", it.filename))
			}
		}
		stmt.Close()
		if _, err = DB.Exec(`DELETE FROM recycleBin WHERE deleted_keep_days < 0`); err != nil {
			log.Println("database [tapbag.db] exec delete deleted_keep_days <= 0 error!")
		}
		time.Sleep(duration)
	}
}

func checkShareKeepDays(duration time.Duration) {
	//Todo
}
