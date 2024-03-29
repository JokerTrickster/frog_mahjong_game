package utils

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// 프론트엔드 받을 에러 형식
type ResError struct {
	ErrType string `json:"errType,omitempty"`
	Msg     string `json:"msg,omitempty"`
}

// 에러 로깅을 위한 에러 형식
type Err struct {
	HttpCode int    `json:"httpCode,omitempty"`
	ErrType  string `json:"errType,omitempty"`
	Msg      string `json:"msg,omitempty"`
	Trace    string `json:"trace,omitempty"`
	From     string `json:"from,omitempty"`
}

// 에러 타입을 구분
type ErrType string

// 에러가 어디서 발생했는지 확인용
type IErrFrom string

const (
	ErrFromClient   = IErrFrom("client")
	ErrFromInternal = IErrFrom("internal")
	ErrFromMongoDB  = IErrFrom("mongoDB")
	ErrFromMysqlDB  = IErrFrom("mysqlDB")
	ErrFromAws      = IErrFrom("aws")
	ErrFromAwsS3    = IErrFrom("aws_s3")
	ErrFromAwsSsm   = IErrFrom("aws_ssm")
	ErrFromNaver    = IErrFrom("naver")
)

const (
	ErrBadParameter         = ErrType("PARAM_BAD")
	ErrAuthFailed           = ErrType("AUTH_FAILED")
	ErrUserNotExist         = ErrType("USER_NOT_EXIST")
	ErrRoomNotExisted       = ErrType("ROOM_NOT_EXISTED")
	ErrRoomImpossibleJoin   = ErrType("ROOM_IMPOSSIBLE_JOIN")
	ErrNotFound             = ErrType("NOT_FOUND")
	ErrAuthInActive         = ErrType("AUTH_INACTIVE")
	ErrUserAlreadyExisted   = ErrType("USER_ALREADY_EXISTED")
	ErrBadToken             = ErrType("TOKEN_BAD")
	ErrAuthPolicyViolation  = ErrType("POLICY_VIOLATION")
	ErrInternalServer       = ErrType("INTERNAL_SERVER")
	ErrInternalDB           = ErrType("INTERNAL_DB")
	ErrPartner              = ErrType("PARTNER")
	ErrNotMatchedLoginInfo  = ErrType("NOT_MATCHED_LOGIN_INFO")
	ErrNotMatchedSignupInfo = ErrType("NOT_MATCHED_SIGNUP_INFO")
	ErrInvalidAuthCode      = ErrType("INVALID_AUTH_CODE")
	ErrExpiredAuthCode      = ErrType("EXPIRED_AUTH_CODE")
)

// 에러 타입에 따라서 httpCode 맵핑
var ErrHttpCode = map[string]int{
	"ROOM_IMPOSSIBLE_JOIN":    http.StatusBadRequest,
	"PARAM_BAD":               http.StatusBadRequest,
	"NOT_FOUND":               http.StatusNotFound,
	"AUTH_FAILED":             http.StatusUnauthorized,
	"AUTH_INACTIVE":           http.StatusForbidden,
	"USER_ALREADY_EXISTED":    http.StatusBadRequest,
	"TOKEN_BAD":               http.StatusUnauthorized,
	"POLICY_VIOLATION":        http.StatusUnauthorized,
	"INTERNAL_SERVER":         http.StatusInternalServerError,
	"INTERNAL_DB":             http.StatusInternalServerError,
	"PARTNER":                 http.StatusInternalServerError,
	"NOT_MATCHED_LOGIN_INFO":  http.StatusBadRequest,
	"NOT_MATCHED_SIGNUP_INFO": http.StatusBadRequest,
	"INVALID_AUTH_CODE":       http.StatusBadRequest,
	"EXPIRED_AUTH_CODE":       http.StatusBadRequest,
	"USER_NOT_EXIST":          http.StatusBadRequest,
	"ROOM_NOT_EXISTED":        http.StatusBadRequest,
}

func ErrorParsing(data string) Err {
	slice := strings.Split(data, "|")
	result := Err{
		HttpCode: ErrHttpCode[slice[0]],
		ErrType:  slice[0],
		Trace:    slice[1],
		Msg:      slice[2],
		From:     slice[3],
	}
	return result
}

func ErrorMsg(ctx context.Context, errType ErrType, trace string, msg string, from IErrFrom) error {

	return fmt.Errorf("%s|%s|%s|%s", errType, trace, msg, from)
}

func (e ErrType) New(errType string, msg string) *ResError {
	return &ResError{ErrType: errType, Msg: msg}
}

func Trace() string {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	_, line := runtime.FuncForPC(pc).FileLine(pc)
	return fmt.Sprintf("%s.L%d", funcName, line)
}
