// 不明白的代码 就往这里面写吧 测试用
package main

import (
	"database/sql"
	"fmt"
	dao2 "go_server/internal/dao"
	goconfig2 "go_server/pkg/goconfig"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// testReadConfig()
	// testSql()
	testConnector()
}

func testConnector() {
	connection := new(dao2.MysqlConnector)
	ip := "127.0.0.1"
	account := "root"
	pwd := "estest"
	schame := "point_and_line"
	connection.Connect(&ip, 3306, &account, &pwd, &schame)
	// user := new(dao.User)
	// user.Name = "est"
	// user.Round = 0
	// user.WinCount = 0
	// user.Rank = 0
	users, _ := connection.QueryByUserId(2)
	for _, user := range users {
		println(user.Id)
		println(user.Name)
	}
}

func testReadConfig() {
	cf, _ := goconfig2.LoadConfigFile("server_config.ini")
	fmt.Println(cf.GetValue(goconfig2.DEFAULT_SECTION, "server.port"))
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
