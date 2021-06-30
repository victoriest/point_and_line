package main

import (
	"go_server/dao"
	"go_server/goconfig"
	"go_server/logic"
	"go_server/server"
	"go_server/utils"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"

	log "github.com/alecthomas/log4go"
)

func main() {
	log.LoadConfiguration("./log4go.config")
	port, ip, dbPort, account, pwd, scheme, protocolType, appId, secret := readServerPort()
	dbCon := new(dao.MysqlConnector)
	iDbPort, _ := strconv.Atoi(dbPort)
	isConnect := dbCon.Connect(&ip, iDbPort, &account, &pwd, &scheme)
	if !isConnect {
		err := log.Warn("mysql connect failed")
		if err != nil {
			return
		}
		return
	}
	var pt server.ProtocolType
	if protocolType == "tcp" {
		pt = server.ProtocolTypeTCP
	} else {
		pt = server.ProtocolTypeWebSocket
	}
	nexus := server.NewNexus(pt, port, logic.TCPHandler,
		logic.ConnectedHandler, logic.DisconnectingHander,
		dbCon, appId, secret)
	nexus.Startup()
}

// 读取配置文件
func readServerPort() (string, string, string, string, string, string, string, string, string) {
	exeFile, _ := exec.LookPath(os.Args[0])
	log.Info(filepath.Dir(exeFile))
	filePath := path.Join(filepath.Dir(exeFile), "./server.config")
	cf, err := goconfig.LoadConfigFile(filePath)
	utils.CheckError(err, true)
	port, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.port")
	utils.CheckError(err, true)
	protocolType, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.type")
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
	appid, err := cf.GetValue(goconfig.DEFAULT_SECTION, "wx.appid")
	utils.CheckError(err, false)
	secret, err := cf.GetValue(goconfig.DEFAULT_SECTION, "wx.secret")
	utils.CheckError(err, false)

	return port, ip, dbPort, account, pwd, schame, protocolType, appid, secret
}
