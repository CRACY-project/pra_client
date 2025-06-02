// Log debug, info, warning messages to the console
// emum loglevel
package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/jimbersoftware/pra_client/logging/color"
	"github.com/sanity-io/litter"
)

type LogLevel = int64

var _logFile *os.File
var _debugLogFile string
var _logCounter int
var _logFileName string
var _logFileLimit int

const (
	DEBUG LogLevel = iota + 1
	INFO
	WARNING
	ERROR
	DEV
)

func InitLog(logfile string, debugEnableLogFile string, logfileLimit int) {
	_logFileName = logfile
	_logFileLimit = logfileLimit
	openFileForWriting()
	_debugLogFile = debugEnableLogFile
}

func openFileForWriting() {
	if _logFile != nil {
		_logFile.Close()
	}

	dirToMake, _ := filepath.Split(_logFileName)
	if _, err := os.Stat(dirToMake); os.IsNotExist(err) {
		err = os.Mkdir(dirToMake, 0755)
		if err != nil {
			log.Fatal("Could not create directory", _logFileName)
		}
	}

	ShrinkFileLines(_logFileName, _logFileLimit, 1000)
	// open file for writing
	var err = error(nil)
	_logFile, err = os.OpenFile(_logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Print("error opening log file: %v", err.Error())
	}
}

func LogD(keyword string, level LogLevel, msgs ...interface{}) { //log for extra debug statements, only if file exists. eg debug.log.netfilter
	if ShouldLogWithDebugKeyword(keyword) {
		Log(level, msgs...)
	}
}
func Log(level LogLevel, msgs ...interface{}) {
	if strings.HasSuffix(os.Args[0], ".test") { // print in unit tests
		log.Print(level, msgs)
		return
	}
	_logCounter++
	if _logCounter > 500 {
		openFileForWriting()
		_logCounter = 0
	}
	logLine := FormatMessagesToString(level, msgs...)
	if logLine != "" {
		_logFile.WriteString(logLine)
	}

	fmt.Print(logLine)

}
func FormatMessagesToString(level LogLevel, msgs ...interface{}) string {
	msg := []string{}
	for _, m := range msgs {
		msg = append(msg, fmt.Sprintf("%v", m))
	}
	msgConcatenated := strings.Join(msg, " ")
	currentTime := time.Now()
	ms := currentTime.Format(".000")
	timeForLog := currentTime.Format("2006-01-02 15:04:05") + ms
	logLine := ""
	switch level {
	case DEBUG:
		if _, err := os.Stat(_debugLogFile); err == nil {
			logLine += color.Yellow + timeForLog + " " + "DEBUG: " + msgConcatenated + color.Reset + "\n"
		}
	case INFO:
		logLine += color.Blue + timeForLog + " " + "INFO: " + msgConcatenated + color.Reset + "\n"
	case WARNING:
		logLine += color.Orange + timeForLog + " " + "WARNING: " + msgConcatenated + color.Reset + "\n"
	case ERROR:
		logLine += color.Red + timeForLog + " " + "ERROR: " + msgConcatenated + color.Reset + "\n"
	case DEV:
		logLine += color.Green + timeForLog + " " + "DEV: " + msgConcatenated + color.Reset + "\n"
	}
	return logLine
}

func ShouldLogWithDebugKeyword(keyword string) bool {
	if _, err := os.Stat(_debugLogFile + "." + keyword); err == nil {
		return true
	}
	return false
}

func Close() {
	_logFile.Close()
}

func LogReader(reader io.Reader, logLevel LogLevel, prefix string) {
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			Log(logLevel, string(buffer[:n]))
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			Log(ERROR, "Error reading from reader: %s", err)
			break
		}
	}
}

func GetLogFileIO() io.Writer {
	return _logFile
}

type Logger struct {
	prefix string
}

func CreateLogger(name string) Logger {
	return Logger{prefix: name}
}

func (l *Logger) Log(level LogLevel, msgs ...interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	// var fnName string
	var fileName string
	var lineNumber int
	if fn == nil {
		fileName = fmt.Sprintf("[%s]", l.prefix)
	} else {
		fileName, lineNumber = fn.FileLine(pc)
		fileName = filepath.Base(fileName)
		fileName = fmt.Sprintf("[%s:%d]", fileName, lineNumber)
		// functionNameParts := strings.Split(fn.Name(), ".")
		// end := functionNameParts[len(functionNameParts)-1]
		// subject := functionNameParts[len(functionNameParts)-2]
		// subject = strings.TrimPrefix(subject, "(*")
		// subject = strings.TrimSuffix(subject, ")")
		// fnName = fmt.Sprintf("[%s::%s]", subject, end)
	}
	Log(level, append([]interface{}{fileName}, msgs...)...)
}

func LogDev(msgs ...interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	var fileName string
	var lineNumber int
	if fn != nil {
		fileName, lineNumber = fn.FileLine(pc)
		fileName = filepath.Base(fileName)
		fileName = fmt.Sprintf("[%s:%d]", fileName, lineNumber)
	}

	processedMsgs := make([]interface{}, len(msgs))
	for index, msg := range msgs {
		switch {
		case msg == nil:
			processedMsgs[index] = "nil"
		case isStructOrPointer(msg):
			processedMsgs[index] = litter.Sdump(msg)
		default:
			processedMsgs[index] = msg
		}
	}

	Log(DEV, append([]interface{}{fileName}, processedMsgs...)...)
}

func isStructOrPointer(val interface{}) bool {
	kind := reflect.TypeOf(val).Kind()
	return kind == reflect.Struct || (kind == reflect.Ptr && reflect.TypeOf(val).Elem().Kind() == reflect.Struct)
}
