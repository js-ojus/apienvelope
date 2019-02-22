package apienvelope

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
)

// RequestEnvelope holds the application layer method and an opaque
// request-specific JSON body.
//
// Examples:
//     {"method": "GET", "body": {"id": 1234}}
//     {"method": "POST", "body": {...}}
//     {"method": "DELETE", "body": {"id": 1234}}
//
// All client requests MUST follow this envelope structure when making
// API requests that do not involve blobs.
//
// Similarly, responses sent by the API server hold the application
// layer response status and one of the following types of JSON bodies:
//     - ("OK", informational message or acknowledgement),
//     - ("Error", error code, error message, a map of important parameters), and
//     - ("OK", an opaque response-specific JSON body).
//
// All clients MUST check the top-level status before parsing the
// response body.
type RequestEnvelope struct {
	Method string          `json:"method"`
	Body   json.RawMessage `json:"body"`
}

// OpenEnvelope opens the envelope by reading the full body of the
// request, and then reading it into an instance of `RequestEnvelope`.
func OpenEnvelope(r io.Reader) (*RequestEnvelope, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "request read failed")
	}
	var env RequestEnvelope
	err = json.Unmarshal(buf, &env)
	if err != nil {
		return nil, errors.Wrap(err, "request unmarshal failed")
	}

	return &env, nil
}

// SendSuccess prepares and writes an informational JSON response based
// on the given message.
func SendSuccess(w io.Writer, msg string) {
	r := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{"OK", msg}
	buf, err := json.Marshal(r)
	if err != nil {
		log.Printf("original error: %T %v\n", errors.Cause(err), errors.Cause(err))
		log.Printf("stack trace:\n%+v\n", err)
		io.WriteString(w, "internal system error")
		return
	}

	w.Write(buf)
}

// SendError prepares and writes an error JSON response based on the
// given values.
func SendError(w io.Writer, e error) {
	msg := fmt.Sprintf("%v", e)
	r := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{"Error", msg}
	buf, err := json.Marshal(r)
	if err != nil {
		log.Printf("original error: %T %v\n", errors.Cause(e), errors.Cause(e))
		log.Printf("stack trace:\n%+v\n", e)
		log.Printf("JSON error:\n%v\n", err)
		io.WriteString(w, "internal system error")
		return
	}

	w.Write(buf)
}

// SendResult prepares and writes a result-carrying response, using the
// given `json.RawMessage`.
func SendResult(w io.Writer, body interface{}) {
	r := struct {
		Status string      `json:"status"`
		Body   interface{} `json:"body"`
	}{"OK", body}
	buf, err := json.Marshal(r)
	if err != nil {
		log.Printf("original error: %T %v\n", errors.Cause(err), errors.Cause(err))
		log.Printf("stack trace:\n%+v\n", err)
		io.WriteString(w, "internal system error")
		return
	}

	w.Write(buf)
}
