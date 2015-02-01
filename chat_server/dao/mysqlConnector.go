package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type MysqlConnector struct {
	connection *sql.DB
}

func (self *MysqlConnector) Connect(ip *string, port int,
	account *string, pwd *string, schame *string) bool {
	if self.connection != nil {
		return false
	}
	uri := *account + ":" + *pwd + "@tcp(" +
		*ip + ":" + strconv.Itoa(port) + ")/" +
		*schame + "?charset=utf8"
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

// 添加用户
func (self *MysqlConnector) Insert(user *User) (int64, error) {
	err := self.connection.Ping()
	if err != nil {
		return -1, err
	}
	result, err := self.connection.Exec(
		"INSERT INTO user(`uname`,`round`,`win_count`,`rank`,`pwd`) VALUES ('" +
			user.Name + "'," + strconv.Itoa(user.Round) + "," +
			strconv.Itoa(user.WinCount) + "," + strconv.Itoa(user.Rank) +
			"," + user.Pwd + ")")
	if err != nil {
		return -1, err
	}
	userId, _ := result.LastInsertId()
	return userId, err
}

// 根据ID查询
func (self *MysqlConnector) QueryByUserId(userId int64) ([]User, error) {
	err := self.connection.Ping()
	if err != nil {
		return nil, err
	}
	result, err := self.connection.Query("SELECT * FROM user WHERE id=" +
		strconv.Itoa(int(userId)))
	defer result.Close()

	users := []User{}
	for result.Next() {
		u := new(User)
		result.Scan(&u.Id, &u.Name, &u.Round, &u.WinCount, &u.Rank, &u.Pwd)
		users = append(users, *u)
	}
	return users, nil
}

// 根据用户名查询
func (self *MysqlConnector) QueryByUserName(userName string,
	pwd string) ([]User, error) {
	err := self.connection.Ping()
	if err != nil {
		return nil, err
	}
	result, err := self.connection.Query("SELECT * FROM user WHERE uname='" +
		userName + "' AND pwd='" + pwd + "'")
	defer result.Close()

	users := []User{}
	for result.Next() {
		u := new(User)
		result.Scan(&u.Id, &u.Name, &u.Round, &u.WinCount, &u.Rank)
		users = append(users, *u)
	}
	return users, nil
}

// 更新用户数据
func (self *MysqlConnector) Update(user *User) (int, error) {
	err := self.connection.Ping()
	if err != nil {
		return -1, err
	}
	var base int
	userId := strconv.FormatInt(user.Id, base)
	_, err = self.connection.Exec("UPDATE user SET name='" +
		user.Name + "', round=" + strconv.Itoa(user.Round) +
		", win_count=" + strconv.Itoa(user.WinCount) +
		", pwd=" + user.Pwd + ", rank= WHERE id=" + userId)
	return 0, err
}

// 删除用户
func (self *MysqlConnector) Delete(userId int) (int, error) {
	err := self.connection.Ping()
	if err != nil {
		return -1, err
	}
	_, err = self.connection.Exec("DELETE FROM user WHERE id=" +
		strconv.Itoa(userId))
	return 0, err

}
