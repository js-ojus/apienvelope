package apienvelope

import (
	"bytes"
	"testing"

	"github.com/pkg/errors"
)

func TestSendSuccess001(t *testing.T) {
	buf := new(bytes.Buffer)
	out := `{"status":"OK","message":"Hello!"}`

	SendSuccess(buf, "Hello!")
	if string(buf.Bytes()) != out {
		t.Errorf("expected: '%s', got: '%s'\n", out, buf)
	}
}

func TestSendError001(t *testing.T) {
	buf := new(bytes.Buffer)
	out := `{"status":"Error","message":"Test error 1001"}`

	SendError(buf, errors.New("Test error 1001"))
	if string(buf.Bytes()) != out {
		t.Errorf("expected: '%s', got: '%s'\n", out, buf)
	}
}

func TestSendError002(t *testing.T) {
	buf := new(bytes.Buffer)
	out := `{"status":"Error","message":"Test error 1001: Test error 1002"}`

	err := errors.New("Test error 1002")
	SendError(buf, errors.Wrap(err, "Test error 1001"))
	if string(buf.Bytes()) != out {
		t.Errorf("expected: '%s', got: '%s'\n", out, buf)
	}
}

func TestSendResult001(t *testing.T) {
	buf := new(bytes.Buffer)
	out := `{"status":"OK","body":{"mode":"test","participants":100}}`

	obj := struct {
		Mode         string `json:"mode"`
		Participants int    `json:"participants"`
	}{"test", 100}
	SendResult(buf, obj)
	if string(buf.Bytes()) != out {
		t.Errorf("expected: '%s', got: '%s'\n", out, buf)
	}
}
