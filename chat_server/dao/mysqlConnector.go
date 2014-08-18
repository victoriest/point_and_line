package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type MysqlConnector struct {
	connection *sql.DB
}

func (self *MysqlConnector) Connect(ip *string, port int, account *string, pwd *string, schame *string) bool {
	if self.connection != nil {
		return false
	}
	uri := *account + ":" + *pwd + "@tcp(" + *ip + ":" + strconv.Itoa(port) + ")/" + *schame + "?charset=utf8"
	db, err := sql.Open("mysql", uri)
	if db == nil || err != nil {
		return false
	}
	self.connection = db
	return true
}

func (self *MysqlConnector) IsClose() bool {
	err := self.connection.Ping()
	if err != nil {
		return true
	}
	return false
}

func (self *MysqlConnector) Close() bool {
	err := self.connection.Close()
	if err != nil {
		return false
	}
	return true
}

func (self *MysqlConnector) Insert(user *User) (int, error) {
	err := self.connection.Ping()
	if err != nil {
		return -1, err
	}
	_, err = self.connection.Exec("INSERT INTO user(`name`,`round`,`win_count`,`rank`) VALUES ('" + user.Name + "'," + strconv.Itoa(user.Round) + "," + strconv.Itoa(user.WinCount) + "," + strconv.Itoa(user.Rank) + ")")
	return 0, err
}

func (self *MysqlConnector) Query(user *User) interface {
	err := self.connection.Ping()
	if err != nil {
		return -1, err
	}
}

func (self *MysqlConnector) Update(user *User) int {
	err := self.connection.Ping()
	if err != nil {
		return -1, err
	}
	_, err = self.connection.Exec("UPDATE user SET name='" + user.Name + "', round=" + strconv.Itoa(user.Round) + ", win_count=" + strconv.Itoa(user.WinCount) + ", rank= WHERE id=" + strconv.Itoa(user.Id))
	return 0, err
}

func (self *MysqlConnector) Delete(userId int) int {
	err := self.connection.Ping()
	if err != nil {
		return -1, err
	}
	_, err = self.connection.Exec("DELETE FROM user WHERE id=" + strconv.Itoa(userId))
	return 0, err

}
