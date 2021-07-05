package main

import (
	"go_server/internal/dao"
	"go_server/internal/logic"
	sev "go_server/internal/server"
	"go_server/pkg/goconfig"
	"go_server/pkg/log"
	utils2 "go_server/pkg/utils"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
)

func main() {
	log.InitZapLogger(true, "go_server", "./", "info", 1)
	port, ip, dbPort, account, pwd, scheme, protocolType, appId, secret := readServerPort()
	dbCon := new(dao.MysqlConnector)
	iDbPort, _ := strconv.Atoi(dbPort)
	isConnect := dbCon.Connect(&ip, iDbPort, &account, &pwd, &scheme)
	if !isConnect {
		log.Warn("mysql connect failed")
		return
	}
	var pt sev.ProtocolType
	if protocolType == "tcp" {
		pt = sev.ProtocolTypeTCP
	} else {
		pt = sev.ProtocolTypeWebSocket
	}
	nexus := sev.NewNexus(pt, port, logic.TCPHandler,
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
	utils2.CheckError(err, true)
	port, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.port")
	utils2.CheckError(err, true)
	protocolType, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.type")
	utils2.CheckError(err, true)
	ip, err := cf.GetValue(goconfig.DEFAULT_SECTION, "db.ip")
	utils2.CheckError(err, true)
	dbPort, err := cf.GetValue(goconfig.DEFAULT_SECTION, "db.port")
	utils2.CheckError(err, true)
	account, err := cf.GetValue(goconfig.DEFAULT_SECTION, "db.user")
	utils2.CheckError(err, true)
	pwd, err := cf.GetValue(goconfig.DEFAULT_SECTION, "db.pwd")
	utils2.CheckError(err, true)
	scheme, err := cf.GetValue(goconfig.DEFAULT_SECTION, "db.scheme")
	utils2.CheckError(err, true)
	appId, err := cf.GetValue(goconfig.DEFAULT_SECTION, "wx.appid")
	utils2.CheckError(err, false)
	secret, err := cf.GetValue(goconfig.DEFAULT_SECTION, "wx.secret")
	utils2.CheckError(err, false)

	return port, ip, dbPort, account, pwd, scheme, protocolType, appId, secret
}
