package main

import (
	"./goconfig"
	vcli "./victoriest.org/client"
	"./victoriest.org/utils"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func main() {
	client := vcli.NewVictoriestClient(readServerConfig())
	client.Startup()
}

func readServerConfig() (string, string) {
	exefile, _ := exec.LookPath(os.Args[0])

	filepath := path.Join(filepath.Dir(exefile), "./server.config")
	cf, err := goconfig.LoadConfigFile(filepath)
	utils.CheckError(err, true)

	host, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.host")
	utils.CheckError(err, true)

	port, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.port")
	utils.CheckError(err, true)
	return host, port
}
