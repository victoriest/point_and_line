package utils

import (
	log "github.com/alecthomas/log4go"
	"go_server/goconfig"
	"os"
)

func CheckError(err error, isQuit bool) {
	if err != nil {
		log.Error(err.Error())
		if isQuit {
			os.Exit(2)
		}
	}
}

type ConfigReader struct {
	fileName string
	config   *goconfig.ConfigFile
}

func NewConfigReader(filePath string) (*ConfigReader, error) {

	cf, err := goconfig.LoadConfigFile(filePath)
	if err != nil {
		return nil, err
	}
	reader := &ConfigReader{fileName: filePath, config: cf}
	return reader, nil
}

func (self ConfigReader) Get(param string) (string, error) {
	return self.config.GetValue(goconfig.DEFAULT_SECTION, param)
}
