// 不明白的代码 就往这里面写吧 测试用
package main

import (
	"./goconfig"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func main() {
	// testReadConfig()
	testSql()
}

func testReadConfig() {
	cf, _ := goconfig.LoadConfigFile("server_config.ini")
	fmt.Println(cf.GetValue(goconfig.DEFAULT_SECTION, "server.port"))
}

func testSql() {
	db, err := sql.Open("mysql", "victoriest:victoriestFuckHacker@asssql@tcp(115.239.252.4:3306)/rock_frog?charset=utf8")
	if err != nil {
		print(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		print(err.Error())
	}

	rows, err := db.Query("select id, name from user")

	defer rows.Close()
	var id int        //定义一个id 变量
	var name string   //定义lvs 变量
	for rows.Next() { //开始循环
		rerr := rows.Scan(&id, &name) //数据指针，会把得到的数据，往刚才id 和 lvs引入
		if rerr == nil {
			fmt.Println(" " + strconv.Itoa(id) + "   " + name) //输出来而已，看看
		}
	}
	// insert_sql := "INSERT INTO xiaorui(lvs) VALUES(?)"
	// _, e4 := db.Exec(insert_sql, "nima")
	// fmt.Println(e4)
}
