package dao

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlConnector struct {
	connection *sql.DB
}

func (connector *MysqlConnector) Connect(ip *string, port int,
	account *string, pwd *string, schame *string) bool {
	if connector.connection != nil {
		return false
	}
	uri := *account + ":" + *pwd + "@tcp(" +
		*ip + ":" + strconv.Itoa(port) + ")/" +
		*schame + "?charset=utf8"
	db, err := sql.Open("mysql", uri)
	if db == nil || err != nil {
		return false
	}
	connector.connection = db
	return true
}

func (connector *MysqlConnector) IsClose() bool {
	err := connector.connection.Ping()
	if err != nil {
		return true
	}
	return false
}

func (connector *MysqlConnector) Close() bool {
	err := connector.connection.Close()
	if err != nil {
		return false
	}
	return true
}

// Insert 添加用户
func (connector *MysqlConnector) Insert(user *User) (int64, error) {
	err := connector.connection.Ping()
	if err != nil {
		return -1, err
	}
	result, err := connector.connection.Exec(
		"INSERT INTO user(`uname`,`round`,`win_count`,`rank`,`pwd`, `open_id`) VALUES ('" +
			user.Name + "'," + strconv.Itoa(user.Round) + "," +
			strconv.Itoa(user.WinCount) + "," + strconv.Itoa(user.Rank) +
			",'" + user.Pwd + "','" + user.OpenId + "')")
	if err != nil {
		return -1, err
	}
	userID, _ := result.LastInsertId()
	return userID, err
}

// QueryByUserId 根据ID查询
func (connector *MysqlConnector) QueryByUserId(userID int64) ([]User, error) {
	err := connector.connection.Ping()
	if err != nil {
		return nil, err
	}
	result, err := connector.connection.Query("SELECT * FROM user WHERE id=" +
		strconv.Itoa(int(userID)))
	defer result.Close()

	users := []User{}
	for result.Next() {
		u := new(User)
		result.Scan(&u.Id, &u.Name, &u.Round, &u.WinCount, &u.Rank, &u.Pwd, &u.OpenId)
		users = append(users, *u)
	}
	return users, nil
}

// QueryByUserName 根据用户名查询
func (connector *MysqlConnector) QueryByUserName(userName string,
	pwd string) ([]User, error) {
	err := connector.connection.Ping()
	if err != nil {
		return nil, err
	}
	result, err := connector.connection.Query("SELECT * FROM user WHERE uname='" +
		userName + "' AND pwd='" + pwd + "'")
	defer result.Close()
	if result == nil || err != nil {
		return nil, err
	}
	users := []User{}
	for result.Next() {
		u := new(User)
		var p string
		result.Scan(&(u.Id), &(u.Name), &(u.Round), &(u.WinCount), &(u.Rank), &p, &(u.OpenId))
		users = append(users, *u)
	}
	if len(users) < 1 {
		return nil, nil
	}
	return users, nil
}

// QueryByOpenId 根据用户名查询
func (connector *MysqlConnector) QueryByOpenId(openId string) ([]User, error) {
	err := connector.connection.Ping()
	if err != nil {
		return nil, err
	}
	result, err := connector.connection.Query("SELECT * FROM user WHERE open_id='" +
		openId + "'")
	defer result.Close()
	if result == nil || err != nil {
		return nil, err
	}
	users := []User{}
	for result.Next() {
		u := new(User)
		var p string
		result.Scan(&(u.Id), &(u.Name), &(u.Round), &(u.WinCount), &(u.Rank), &p, &(u.OpenId))
		users = append(users, *u)
	}
	if len(users) < 1 {
		return nil, nil
	}
	return users, nil
}

// Update 更新用户数据
func (connector *MysqlConnector) Update(user *User) (int, error) {
	err := connector.connection.Ping()
	if err != nil {
		return -1, err
	}
	var base int
	userID := strconv.FormatInt(user.Id, base)
	_, err = connector.connection.Exec("UPDATE user SET name='" +
		user.Name + "', round=" + strconv.Itoa(user.Round) +
		", win_count=" + strconv.Itoa(user.WinCount) +
		", pwd=" + user.Pwd + ", open_id=" + user.OpenId + ", rank= WHERE id=" + userID)
	return 0, err
}

// Delete 删除用户
func (connector *MysqlConnector) Delete(userID int) (int, error) {
	err := connector.connection.Ping()
	if err != nil {
		return -1, err
	}
	_, err = connector.connection.Exec("DELETE FROM user WHERE id=" +
		strconv.Itoa(userID))
	return 0, err

}
