package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"

	"./dao"
	"./goconfig"
	"./logic"
	sev "./server"
	"./utils"
	log "github.com/alecthomas/log4go"
)

func main() {
	log.LoadConfiguration("./log4go.config")
	port, ip, dbPort, account, pwd, schame := readServerPort()
	dbCon := new(dao.MysqlConnector)
	iDbPort, _ := strconv.Atoi(dbPort)
	isConnect := dbCon.Connect(&ip, iDbPort, &account, &pwd, &schame)
	if !isConnect {
		log.Warn("mysql connect faild")
		return
	}
	server := sev.NewNexus(port, logic.TcpHandler,
		logic.ConnectedHandler, logic.DisconnectingHander,
		dbCon)
	server.Startup()
}

// 读取配置文件
func readServerPort() (string, string, string, string, string, string) {
	exefile, _ := exec.LookPath(os.Args[0])
	log.Info(filepath.Dir(exefile))
	filepath := path.Join(filepath.Dir(exefile), "./server.config")
	cf, err := goconfig.LoadConfigFile(filepath)
	utils.CheckError(err, true)
	port, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.port")
	utils.CheckError(err, true)
	ip, err := cf.GetValue(goconfig.DEFAULT_SECTION, "db.ip")
	utils.CheckError(err, true)
	dbPort, err := cf.GetValue(goconfig.DEFAULT_SECTION, "db.port")
	utils.CheckError(err, true)
	account, err := cf.GetValue(goconfig.DEFAULT_SECTION, "db.user")
	utils.CheckError(err, true)
	pwd, err := cf.GetValue(goconfig.DEFAULT_SECTION, "db.pwd")
	utils.CheckError(err, true)
	schame, err := cf.GetValue(goconfig.DEFAULT_SECTION, "db.schame")
	utils.CheckError(err, true)
	return port, ip, dbPort, account, pwd, schame
}
