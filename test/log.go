/**
* test for log.go
*/
package main

import (
	"github.com/cz-it/log4golang/log"
	"fmt"
	"os"
)

func main(){
	fmt.Println("Testing:")
	logger := log.NewLogger("./a/b","base","nimei",log.DEBUG)
	logger.Debug("nimei")
	var err error
	err = os.Mkdir("./nimei",os.ModePerm)
	if err != nil {
		fmt.Println("Mkdir;",err)
	}
	err = os.MkdirAll("./nimei2",os.ModePerm)
	if err != nil {
		fmt.Println("MkdirAll;",err)
	}
}




















