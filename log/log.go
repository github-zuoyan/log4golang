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
	"errors"
	"runtime"
	"time"
)

type Level int   // loose enum type . maybe have some other define method
const (
	// content of emnu Level ,level of log
	NULL   = 1<< iota
	DEBUG
	INFO
	TRACE
	WARNNING
	ERROR
	FATAL
)

type Outputer int 
const (
	STD   = iota
	FILE
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
	debugOutputer Outputer
	debugSwitch bool
	callDepth  int
}
var loggers map[string] *Logger
var _logger *Logger

func Init() error {
	loggers = make(map[string] *Logger)
	_logger = NewLogger("./","log","Log4Golang",DEBUG)
	_logger.SetCallDepth(3)
	return nil
}

func GetLogger(logName string) *Logger{
	logger,err := loggers[logName]
	if err != true{
		return nil
	}
	return logger
}

func NewLogger(path,baseName,logName string,level Level)( *Logger){
	var err error
	logger := &Logger{path:path,baseName:baseName,logName:logName,level:level}

	err = os.MkdirAll(path,os.ModePerm)
	if err != nil {
		panic(err)
	}

	path = strings.TrimSuffix(path,"/")
	flag := os.O_WRONLY|os.O_APPEND|os.O_CREATE

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

	logger.debugSwitch   = true
	logger.debugOutputer = STD
	logger.callDepth     = 2

	loggers[logName] = logger
	return logger
}

func (l *Logger) SetCallDepth(d int){
	l.callDepth = d
}

func (l *Logger) OpenDebug(){
	l.debugSwitch = true
}

func (l *Logger) getFileLine() string{
	_, file, line, ok := runtime.Caller(l.callDepth)
	if !ok {
		file = "???"
		line = 0
	}
	
	return file+":"+itoa(line,-1)
}

/**
* Change from Golang's log.go
* Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
* Knows the buffer has capacity.
*/
func itoa(i int, wid int) string {
	var u uint = uint(i)
    if u == 0 && wid <= 1 {
		return "0"
	}

    // Assemble decimal in reverse order.
    var b [32]byte
    bp := len(b)
    for ; u > 0 || wid > 0; u /= 10 {
		bp--
		wid--
		b[bp] = byte(u%10) + '0'
	}
	return string(b[bp:])
}

func (l *Logger) getTime() string{
	// Time is yyyy-mm-dd hh:mm:ss.microsec
	var buf  []byte
	t := time.Now()
	year, month, day := t.Date()
	buf = append(buf, itoa(int(year),4)+"-" ...)
	buf = append(buf,itoa(int(month),2)+ "-" ...)
	buf = append(buf, itoa(int(day),2)+" " ...)
	
	hour, min, sec := t.Clock()
	buf = append(buf,itoa(hour,2)+ ":" ...)
	buf = append(buf,itoa(min,2)+ ":" ...) 
	buf = append(buf,itoa(sec,2) ...)

	buf = append(buf, '.')
	buf = append(buf,itoa(t.Nanosecond()/1e3,6) ...)

	return string(buf[:])
}

func (l *Logger) Output(fd io.Writer,level,prefix string,format string,v... interface{}) (err error) {
	var msg string
	if format== ""  {
		msg = fmt.Sprintln(v...)
	} else {
		msg = fmt.Sprintf(format,v...)
	}

	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.buf = l.buf[:0]
	
	l.buf = append(l.buf,"["+l.logName+"]" ...)
	l.buf = append(l.buf,level ...)
	l.buf = append(l.buf,prefix ...)

	l.buf = append(l.buf,":"+msg ... )
	if len(msg)>0 && msg[len(msg)-1]!= '\n'{
		l.buf = append(l.buf,'\n')
	}

	_,err = fd.Write(l.buf)
	return 
}


func (l *Logger) CloseDebug(){
	l.debugSwitch = false
}

func (l *Logger) SetDebugOutput(o Outputer){
	l.debugOutputer = o
}

func (l *Logger) Debug(format string,v... interface{}) error {
	var fd io.Writer
	fd = nil

	if ! l.debugSwitch {
		return nil
	}

	if l.debugOutputer == STD {
		fd  = os.Stdin
	} else if l.debugOutputer == FILE {
		fd = l.logFd
	} else {
		return errors.New("Debug output is invalied!")
	}
	
	
	l.Output(fd,"[DEBUG]","["+l.getTime()+"]["+l.getFileLine()+"]",format,v...)
	return nil
}

func (l *Logger) Info(format string,v...interface{}) error{
	err := l.Output(l.logFd,"[INFO]","",format,v...)
	return err
}

func (l *Logger) Warning(format string,v...interface{}) error{
	err := l.Output(l.logFd,"[WARNING]","",format,v...)
	return err
}


func (l *Logger) Trace(format string,v...interface{}) error{
	err := l.Output(l.trcFd,"[TRACE]","["+l.getTime()+"]["+l.getFileLine()+"]",format,v...)
	return err
}

func (l *Logger) Error(format string,v...interface{}) error{
	err := l.Output(l.errFd,"[ERROR]","["+l.getTime()+"]["+l.getFileLine()+"]",format,v...)
	return err
}

func (l *Logger) Fatal(format string,v... interface{}) error{
	err := l.Output(l.errFd,"[FATAL]","["+l.getTime()+"]["+l.getFileLine()+"]",format,v...)
	return err
}



func Debug(format string,v... interface{}) error{
	return _logger.Debug(format,v...)
}

func Info(format string,v... interface{}) error{
	return _logger.Info(format,v...)
}

func Warning(format string,v... interface{}) error{
	return _logger.Warning(format,v...)
}

func Trace(format string,v... interface{}) error{
	return _logger.Trace(format ,v... )
}


func Error(format string,v... interface{}) error{
	return _logger.Error(format,v...)
}

func Fatal(format string,v... interface{}) error{
	return _logger.Fatal(format,v...)
}




















