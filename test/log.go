/**
* test for log.go
*/
package main

import (
	"github.com/cz-it/log4golang/log"
	"fmt"
)

func main(){
	fmt.Println("Testing:")
	err := log.Init()
	if err != nil{
		panic("log init error")
	}
	logger := log.NewLogger("./a/b","test","nimei",log.DEBUG)
	logger.Debug("logger's debug")
	logger.Info("logger's info")
	logger.Warning("logger's Waring")
	logger.Trace("logger's Trace")
	logger.Error("logger's error")
	logger.Fatal("logger's fatal")

	log.Debug("Debug")
	log.Info("Info")
	log.Warning("Warning")
	log.Trace("Trace")
	log.Error("Error")
	log.Fatal("fatal end")

	log.Debug("","logger is",logger)
}




















