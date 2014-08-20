package main

import (
	"./goconfig"
	"./logic"
	sev "./server"
	"./utils"
	log "code.google.com/p/log4go"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func main() {
	log.LoadConfiguration("./log4go.config")
	server := sev.NewNexus(readServerPort(), logic.TcpHandler, logic.ConnectedHandler, logic.DisconnectingHander)
	server.Startup()
}

// 读取配置文件
func readServerPort() string {
	exefile, _ := exec.LookPath(os.Args[0])
	log.Info(filepath.Dir(exefile))
	filepath := path.Join(filepath.Dir(exefile), "./server.config")
	cf, err := goconfig.LoadConfigFile(filepath)
	utils.CheckError(err, true)
	port, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.port")
	utils.CheckError(err, true)
	return port
}
