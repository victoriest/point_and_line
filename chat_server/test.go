package main

import (
	"./goconfig"
	"fmt"
)

func main() {
	cf, _ := goconfig.LoadConfigFile("server_config.ini")
	fmt.Println(cf.GetValue(goconfig.DEFAULT_SECTION, "server.port"))
}
