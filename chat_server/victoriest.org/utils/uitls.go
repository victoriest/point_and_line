package utils

import (
	log "code.google.com/p/log4go"
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
