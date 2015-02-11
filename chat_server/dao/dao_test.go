package dao

import (
	"testing"
)

func TestDao(t *testing.T) {
	dbCon := new(MysqlConnector)

	var ip string = "115.159.40.89"
	var name string = "victoriest"
	var pwd string = "victoriestFuckHacker@asssql"
	var schame string = "point_and_line"

	dbCon.Connect(&ip, 3306, &name, &pwd, &schame)

	user := &User{}
	user.Name = "testDao"
	user.Round = 0
	user.WinCount = 0
	user.Rank = 0
	user.Pwd = "testDaoPwd"

	userId, err1 := dbCon.Insert(user)

	if userId > 0 {
		t.Log("insert test ok!")
	} else {
		t.Errorf("insert test error", err1)
	}

	//user.Id = userId
	//user.Round = 10
	//user.WinCount = 11
	//user.Rank = 12

	//dbCon.Update(user)
}
