package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

const (
	colorInfo    = "\033[1;36m%s\033[0m"
	colorWarning = "\033[1;33m%s\033[0m"
	colorError   = "\033[1;31m%s\033[0m"
)

type Logging interface {
	MakeLog() error
	MakeErrorLog(requestParam map[string]interface{}) error
}

type Log struct {
	Project   string    `json:"project"`
	Type      string    `json:"type"`
	UserID    string    `json:"userID"`
	Url       string    `json:"url"`
	Method    string    `json:"method"`
	Latency   int64     `json:"latency"`
	HttpCode  int       `json:"httpCode"`
	RequestID string    `json:"requestID"`
	ErrorInfo ErrorInfo `json:"errorInfo,omitempty"`
}

type ErrorInfo struct {
	Stack     string `json:"stack,omitempty"`
	ErrorType string `json:"errorType,omitempty"`
	Msg       string `json:"msg,omitempty"`
	From      string `json:"from,omitempty"`
}

func InitLogging() error {
	infoFile, err := os.OpenFile("/logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open info log file: %v", err)
	}
	errorFile, err := os.OpenFile("/logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open error log file: %v", err)
	}
	warningFile, err := os.OpenFile("/logs/warning.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open error log file: %v", err)
	}
	infoLogger = log.New(infoFile, fmt.Sprintf(colorInfo, "[INFO] "), log.LstdFlags)
	warningLogger = log.New(warningFile, fmt.Sprintf(colorWarning, "[WARNING] "), log.LstdFlags)
	errorLogger = log.New(errorFile, fmt.Sprintf(colorError, "[ERROR] "), log.LstdFlags)
	return nil
}

func (l *Log) MakeLog(userID string, url string, method string, startTime time.Time, httpCode int, requestID string) error {
	l.Project = "frog"
	l.Type = "access"
	l.UserID = userID
	l.Url = url
	l.Method = method
	l.Latency = time.Since(startTime).Milliseconds()
	l.HttpCode = httpCode
	l.RequestID = requestID

	return nil
}
func (l *Log) MakeErrorLog(res Err) error {
	l.Type = "error"
	errInfo := ErrorInfo{
		Stack:     res.Trace,
		Msg:       res.Msg,
		ErrorType: res.ErrType,
		From:      res.From,
	}
	l.ErrorInfo = errInfo
	return nil
}

// LogInfo : info level log
func LogInfo(logContent interface{}) {
	if Env.IsLocal {
		fmt.Printf("[INFO] %s\n", getStringFromInterface(logContent))
	} else {
		infoLogger.Printf("%s\n", getStringFromInterface(logContent))
	}
}

// LogWarning : warning level log
func LogWarning(logContent interface{}) {
	if Env.IsLocal {
		fmt.Printf("[WARNING] %s\n", getStringFromInterface(logContent))
	} else {
		warningLogger.Printf("%s\n", getStringFromInterface(logContent))
	}
}

// LogError : error level log
func LogError(logContent interface{}) {
	if Env.IsLocal {
		fmt.Printf("[ERROR] %s\n", getStringFromInterface(logContent))
	} else {
		errorLogger.Printf("%s\n", getStringFromInterface(logContent))
	}
}

// get string from any type.
func getStringFromInterface(logContent interface{}) string {
	var result string
	// if its struct, convert to json string
	if reflect.Indirect(reflect.ValueOf(logContent)).Kind() == reflect.Struct {
		raw, err := json.Marshal(logContent)
		if err != nil {
			return fmt.Sprintf("%v", logContent)
		}
		result = string(raw)
	} else {
		result = fmt.Sprintf("%v", logContent)
	}
	return result
}
