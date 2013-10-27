// 不明白的代码 就往这里面写吧 测试用
package main

import (
	"./goconfig"
	"fmt"
)

func main() {
	testReadConfig()
}

func testReadConfig() {
	cf, _ := goconfig.LoadConfigFile("server_config.ini")
	fmt.Println(cf.GetValue(goconfig.DEFAULT_SECTION, "server.port"))
}
