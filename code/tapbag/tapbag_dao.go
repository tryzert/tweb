package tapbag

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
	"time"
	"tweb/code/tool"
)

var DB *sql.DB
var daoSrcPath string

func init() {
	exist, err := tool.FileExist("databases/tapbag.db")
	if err != nil {
		log.Panicln("数据库 tapbag.db 初始化失败！")
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
		"file_type" VARCHAR(10),
		"share_time" VARCHAR(10),
		"share_code" INT NOT NULL,
		"share_url" TEXT NOT NULL,
		"share_forever" INT NOT NULL DEFAULT 0,
		"shared_keep_days" INT NOT NULL DEFAULT 30
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


func queryAllRecycleBinRecords() error {
	rows, err := DB.Query(`SELECT * FROM recycleBin`)
	if err != nil {
		return err
	}
	var (
		id int
		orirelpath string
		filename string
		deletetime string
		keepdays string
	)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id, &orirelpath, &filename, &deletetime, &keepdays)
		fmt.Println(id, orirelpath, filename, deletetime, keepdays)
	}
	return nil
}

//delete from tablename;
//update sqlite_sequence SET seq = 0 where name ='tablename';
func clearRecycleBin() error {
	_, err := DB.Exec(`DELETE FROM recycleBin`)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`UPDATE sqlite_sequence SET seq = 0 WHERE name ='recycleBin';`)
	if err != nil {
		return err
	}
	err = os.RemoveAll(filepath.Join(daoSrcPath, ".tweb/recycleBin"))
	return err
}

//recover files from recycle bin to their origin dirs
func recoverFiles() {

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
				os.RemoveAll(filepath.Join(daoSrcPath, ".tweb/recycleBin", it.filename))
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
