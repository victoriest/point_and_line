package main

import (
	"./goconfig"
	vserv "./victoriest.org/server"
	"./victoriest.org/utils"
	log "code.google.com/p/log4go"
	"os"
	"os/exec"
	// "os/signal"
	"path"
	"path/filepath"
)

// 退出信号量
// var quitSp chan bool

func main() {
	log.LoadConfiguration("./log4go.config")

	// // 监测退出程序的信号量
	// sign := make(chan os.Signal, 1)

	server := vserv.NewVictoriestServer(readServerPort())
	server.Startup()

	// signal.Notify(sign, os.Interrupt, os.Kill)
	// <-sign
	// log.Info("quit")
	// server.Shutdown()
}

/**
 * 读取配置文件
 */
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
