package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"main/features/frog/model/entity"
	v2entity "main/features/v2ws/model/entity"
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
	Project      string                 `json:"project"`
	Type         string                 `json:"type"`
	UserID       string                 `json:"userID"`
	Env          string                 `json:"env"`
	Created      string                 `json:"created"`
	Url          string                 `json:"url"`
	Method       string                 `json:"method"`
	Latency      int64                  `json:"latency"`
	HttpCode     int                    `json:"httpCode"`
	RequestID    string                 `json:"requestID"`
	RequestBody  map[string]interface{} `json:"requestBody,omitempty"`
	RequestPath  map[string]string      `json:"requestPath,omitempty"`
	RequestQuery map[string][]string    `json:"requestQuery,omitempty"`
	ErrorInfo    ErrorInfo              `json:"errorInfo,omitempty"`
	Message      string                 `json:"message,omitempty"`
}

type ErrorInfo struct {
	Stack     string `json:"stack,omitempty"`
	ErrorType string `json:"errorType,omitempty"`
	Msg       string `json:"msg,omitempty"`
	From      string `json:"from,omitempty"`
}

func InitLogging() error {
	infoFile, err := os.OpenFile("../logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open info log file: %v", err)
	}
	errorFile, err := os.OpenFile("../logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open error log file: %v", err)
	}
	warningFile, err := os.OpenFile("../logs/warning.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open error log file: %v", err)
	}
	infoLogger = log.New(infoFile, "", 0)
	warningLogger = log.New(warningFile, "", 0)
	errorLogger = log.New(errorFile, "", 0)
	return nil
}

func (l *Log) MakeLog(userID string, url string, method string, startTime time.Time, httpCode int, requestID string, requestBody map[string]interface{}, queryParams map[string][]string, pathValues map[string]string) error {
	l.Project = "frog"
	l.Type = "info"
	l.Env = Env.Env
	l.UserID = userID
	l.Created = startTime.Format("2006-01-02 15:04:05")
	l.Url = url
	l.Method = method
	l.Latency = time.Since(startTime).Milliseconds()
	l.HttpCode = httpCode
	l.RequestID = requestID
	if requestBody != nil {
		l.RequestBody = requestBody
	}
	if pathValues != nil {
		l.RequestPath = pathValues
	}
	if queryParams != nil {
		l.RequestQuery = queryParams
	}
	return nil
}

func (l *Log) MakeWSLog(msg entity.WSMessage) error {
	startTime := time.Now()
	l.Project = "frog"
	l.Type = "info"
	l.Env = Env.Env
	l.UserID = string(msg.UserID)
	l.Created = startTime.Format("2006-01-02 15:04:05")
	l.Url = msg.Event
	l.Method = "ws"
	l.Latency = time.Since(startTime).Milliseconds()
	l.HttpCode = 200
	l.RequestID = string(msg.RoomID)
	l.Message = msg.Message
	return nil
}
func (l *Log) V2MakeWSLog(msg v2entity.WSMessage) error {
	startTime := time.Now()
	l.Project = "frog"
	l.Type = "info"
	l.Env = Env.Env
	l.UserID = string(msg.UserID)
	l.Created = startTime.Format("2006-01-02 15:04:05")
	l.Url = msg.Event
	l.Method = "ws"
	l.Latency = time.Since(startTime).Milliseconds()
	l.HttpCode = 200
	l.RequestID = string(msg.RoomID)
	l.Message = msg.Message
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
		infoLogger.Printf("%s", getStringFromInterface(logContent))
	}
}

// LogWarning : warning level log
func LogWarning(logContent interface{}) {
	if Env.IsLocal {
		fmt.Printf("[WARNING] %s\n", getStringFromInterface(logContent))
	} else {
		warningLogger.Printf("%s", getStringFromInterface(logContent))
	}
}

// LogError : error level log
func LogError(logContent interface{}) {
	if Env.IsLocal {
		fmt.Printf("[ERROR] %s\n", getStringFromInterface(logContent))
	} else {
		errorLogger.Printf("%s", getStringFromInterface(logContent))
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
