package errno

import (
	"encoding/json"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestCopyWithMessage(t *testing.T) {
	oldMsg := DB.GetMessage()
	newMsg := "wrong params"
	e := DB.CopyWithMessage(newMsg)
	if e.GetMessage() != newMsg {
		t.Errorf("wanted: %v, got: %v", newMsg, e.GetMessage())
	}
	if DB.GetMessage() != oldMsg {
		t.Errorf("wanted: %v, got: %v", oldMsg, DB.GetMessage())
	}
}

func TestCopyWithPrompt(t *testing.T) {
	oldPrompt := DB.GetPrompt()
	newPrompt := "try again"
	e := DB.CopyWithPrompt(newPrompt)
	if e.GetPrompt() != newPrompt {
		t.Errorf("wanted: %v, got: %v", newPrompt, e.GetPrompt())
	}
	if DB.GetPrompt() != oldPrompt {
		t.Errorf("wanted: %v, got: %v", oldPrompt, DB.GetPrompt())
	}
}

func toStatusError(no ErrNo) error {
	str, _ := json.Marshal(no)
	return status.Errorf(codes.Code(no.GetCode()), string(str))
}

func TestClone(t *testing.T) {
	newErr := clone(ParamWrong).(*errNo)
	if *newErr != *ParamWrong.(*errNo) {
		t.Errorf("wanted: %#v, got: %#v", ParamWrong, newErr)
	}
}

func TestParse(t *testing.T) {
	if Parse(nil) != nil {
		t.Errorf("parse nil should return nil")
	}
	e := Parse(Unknown)
	if e != Unknown {
		t.Errorf("wanted :%v, got: %v", Unknown, e)
	}
	e = Parse(errors.New("an error"))
	wanted := Unknown.CopyWithMessage("an error")
	if *e.(*errNo) != *wanted.(*errNo) {
		t.Errorf("wanted: %v, got: %v", wanted, e)
	}

	statusErr := toStatusError(ParamWrong.CopyWithPrompt("this is new prompt"))
	newErr := Parse(statusErr)

	if ParamWrong.GetPrompt() != "参数错误，请确认后重试" {
		t.Errorf("wanted: %v, got: %v", "参数错误，请确认后重试", ParamWrong.GetPrompt())
	}

	if newErr.GetPrompt() != "this is new prompt" {
		t.Errorf("wanted: %v, got: %v", "this is new prompt", newErr.GetPrompt())
	}

	if newErr.GetHTTPStatus() != ParamWrong.GetHTTPStatus() {
		t.Errorf("wanted: %v, got: %v", ParamWrong.GetHTTPStatus(), newErr.GetHTTPStatus())
	}
}
