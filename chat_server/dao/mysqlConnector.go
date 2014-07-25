package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlConnector struct {
}

func (self *MysqlConnector) Connect(ip *string, port int, pwd *string, schame *string) {

}

func (self *MysqlConnector) Check() int {

}

func (self *MysqlConnector) Query(sql *string) interface{}
