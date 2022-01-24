package errno

import (
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/status"
	"reflect"
)

// api, service 错误码尽量收敛到此
var (
	// 10000 ~ 10999 服务内部错误
	Unknown    = register(500, 10000, "internal error", "")
	DB         = register(500, 10001, "internal error", "")
	Redis      = register(500, 10002, "internal error", "")
	RPC        = register(500, 10003, "internal error", "")
	ThirdParty = register(500, 10004, "call third_party api failed", "调用第三方API异常")

	// 11000 ~ 11999 通用业务错误
	ParamWrong       = register(400, 11000, "invalid param", "参数错误，请确认后重试")
	SessionExpired   = register(401, 11001, "session expired", "账号授权过期，请重新登录")
	TokenInvalid     = register(401, 11002, "invalid token", "")
	TicketInvalid    = register(401, 11003, "invalid ticket", "")
	NoPermission     = register(403, 11004, "no permission", "无权限")
	RedirectIllegal  = register(403, 11005, "redirect url illegal", "")
	CodeInvalid      = register(401, 11006, "invalid code", "手机验证码错误")
	Existed          = register(409, 11007, "resource is already existed", "资源已存在")
	FrequencyExceeds = register(429, 11008, "request too frequent", "请稍后再试")
	InvalidID        = register(401, 11009, "invalid id", "无效的ID")
	NotFound         = register(404, 11010, "not found", "资源不存在")
	Deprecated       = register(410, 11011, "api deprecated", "该方法已废弃")

	// 12000 ~ 12999 内部账户相关错误
	UserNotFound = register(404, 12000, "user not found", "用户不存在")
)

var (
	errNos = make(map[int32]ErrNo)
)

// ErrNo represents an error condition which is generally used through the whole call chain
// from front-end to inner services.
type ErrNo interface {
	Error() string
	GetHTTPStatus() int
	GetCode() int32
	GetMessage() string
	GetPrompt() string

	CopyWithMessage(message string) ErrNo
	CopyWithPrompt(message string) ErrNo
}

// errNo implements ErrNo.
// We keep errNo unexported so it can only be declared in this package.
// Defining errNo in API's own package should never happen.
type errNo struct {
	// HTTPStatus is the http status code of a http response.
	HTTPStatus int `json:"-"`
	// Code is for documentation-specific notation of errors. Client can also use it to distinguish different errors.
	Code int32 `json:"code"`
	// Message is the human readable messages that summarize the context, cause and general solution for the error.
	Message string `json:"message"`
	// Prompt is the hint presented to client users when the error occurs.
	Prompt string `json:"prompt,omitempty"`
}

func newErrNo(httpStatus int, code int32, message, prompt string) *errNo {
	return &errNo{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
		Prompt:     prompt,
	}

}

func Register(httpStatus int, code int32, message, prompt string) ErrNo {
	return register(httpStatus, code, message, prompt)
}

func register(httpStatus int, code int32, message, prompt string) ErrNo {
	if _, exist := errNos[code]; exist {
		panic(fmt.Sprintf("errno %d already exists", code))
	}
	e := newErrNo(httpStatus, code, message, prompt)
	errNos[code] = e
	return e
}

// Error implements error interface.
func (en *errNo) Error() string {
	return fmt.Sprintf("code: %d, message: %s, prompt: %s", en.Code, en.Message, en.Prompt)
}

// GetHTTPStatus returns the HTTPStatus field of errNo.
func (en *errNo) GetHTTPStatus() int {
	return en.HTTPStatus
}

// GetCode returns the Code field of errNo.
func (en *errNo) GetCode() int32 {
	return en.Code
}

// GetMessage returns the Message field of errNo.
func (en *errNo) GetMessage() string {
	return en.Message
}

// GetPrompt returns the Prompt field of errNo.
func (en *errNo) GetPrompt() string {
	return en.Prompt
}

// CopyWithMessage copies the errNo instance and returns a new one with given message.
func (en *errNo) CopyWithMessage(message string) ErrNo {
	return newErrNo(en.HTTPStatus, en.Code, message, en.Prompt)
}

// CopyWithPrompt copies the errNo instance and returns a new one with given prompt.
func (en *errNo) CopyWithPrompt(prompt string) ErrNo {
	return newErrNo(en.HTTPStatus, en.Code, en.Message, prompt)
}

// Parse converts an interface to an ErrNo.
func Parse(err interface{}) ErrNo {
	if err == nil {
		return nil
	}
	switch e := err.(type) {
	case ErrNo:
		return e
	case error:
		if s, ok := status.FromError(e); ok {
			en := GetErrNo(int32(s.Code()))
			newEn := clone(en)
			_ = json.Unmarshal([]byte(s.Message()), newEn)
			return newEn.(ErrNo)
		}
		return Unknown.CopyWithMessage(e.Error())
	default:
		return Unknown
	}
}

func clone(expected interface{}) interface{} {
	v := reflect.ValueOf(expected)
	t := v.Type()
	if reflect.Ptr == t.Kind() {
		t = reflect.Indirect(v).Type()
	}
	rel := reflect.New(t)
	for i := 0; i < t.NumField(); i++ {
		rel.Elem().Field(i).Set(v.Elem().Field(i))
	}
	return rel.Interface()
}

// GetErrNo returns the ErrNo by code.
func GetErrNo(code int32) ErrNo {
	if en, exist := errNos[code]; exist {
		return en
	}
	return Unknown
}
