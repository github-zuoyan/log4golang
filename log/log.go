/**
* A log package for golang
* 
* not instead of golang's log but a replenish
*/

package log

import (
	"sync"
	"io"
	"fmt"
	"os"
	"strings"
)

type Level int   // loose enum type . maybe have some other define method
const (
	// content of emnu Level ,level of log
	DEBUG   = 1<< iota
	INFO
	TRACE
	WARNNING
	ERROR
	FATAL
)

type Logger struct {
	mtx        sync.Mutex
	logFd      io.Writer
	trcFd      io.Writer
	errFd      io.Writer
	level      Level
	buf        []byte
	path       string
	baseName   string
    logName    string
}

func NewLogger(path,baseName,logName string,level Level)( *Logger){
	var err error
	logger := &Logger{path:path,baseName:baseName,logName:logName,level:level}
	logger.buf = append(logger.buf,"["+logName+"]" ...)

	err = os.MkdirAll(path,os.ModePerm)
	if err != nil {
		panic(err)
	}

	path = strings.TrimSuffix(path,"/")
	var flag := os.O_WRONLY|os.O_APPEND|os.O_CREATE

	logger.logFd, err= os.OpenFile(path+"/"+baseName+".log",flag,0666)
	if err != nil {
		panic(err)
	}

	logger.errFd, err= os.OpenFile(path+"/"+baseName+".err",flag,0666)
	if err != nil {
		panic(err)
	}

	logger.trcFd, err= os.OpenFile(path+"/"+baseName+".trace",flag,0666)
	if err != nil {
		panic(err)
	}

	
	return logger
}

func (l *Logger) Debug(v ... interface{}){
	fmt.Printf("%s\n",l.buf)
}




















