package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
)

const DEBUG = "D"
const ERROR = "E"

var BKSLogger *log.Logger
var LogLevel string

func init() {
	BKSLogger = log.New(os.Stderr, "BKSLog:\t", log.Ldate|log.Ltime)
	LogLevel = os.Getenv("BKS_log_level")
}

func BKSLog(logLevel string, message string, err error) {
	//RMS is acronym of the project name
	// Extract information about the caller of the log function, if requested.
	var callerInfo, callingFuncName, moduleAndFileName string

	pc, fullFilePath, line, ok := runtime.Caller(1)
	if ok {
		callingFuncName = runtime.FuncForPC(pc).Name()
		// We only want to print or examine file and package name, so use the
		// last two elements of the full path. The path package deals with
		// different path formats on different systems, so we use that instead
		// of just string-split.
		dirPath, fileName := path.Split(fullFilePath)
		var moduleName string
		if dirPath != "" {
			dirPath = dirPath[:len(dirPath)-1]
			_, moduleName = path.Split(dirPath)
		}

		moduleAndFileName = moduleName + "/" + fileName

		callerInfo = fmt.Sprintf("[ %s:%d (%s) : ] ", moduleAndFileName, line, callingFuncName)
	}

	switch LogLevel {
	case DEBUG:
		BKSLogger.Println(callerInfo+"\n"+message+" ** \n", err, "\n****\n")
	case ERROR:
		if logLevel == ERROR {
			BKSLogger.Println(callerInfo+"\n"+message+" ** \n", err, "\n****\n")
		}
	default:
		BKSLogger.Println(callerInfo+"\n"+message+" ** \n", err, "\n****\n")
	}
}

func PrintStruct(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
