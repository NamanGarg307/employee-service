package errs

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Err is a custom error object complying to Headere and StatusCoder interfaces
// so that it can be readily used by custom HTTPHandler implementation
type HTTPError struct {
	Type        string      `json:"type"`
	Title       string      `json:"title,omitempty"`
	Status      int         `json:"status"`
	Message     interface{} `json:"message"`
	Errors      interface{} `json:"errors,omitempty"`
	HTTPHeaders http.Header
}

func (e *HTTPError) Error() string {
	msg, _ := json.Marshal(e.Message)
	return string(msg)
}

func (e *HTTPError) Headers() http.Header {
	return e.HTTPHeaders
}

func (e *HTTPError) StatusCode() int {
	if e.Status > 0 {
		return e.Status
	}
	return http.StatusOK

}

func ErrRes(title string, status int,
	message interface{}, types ...string) error {

	return &HTTPError{Type: strings.Join(types, ""),
		Title:   title,
		Status:  status,
		Message: message}
}

func (e *HTTPError) MarshalJSON() ([]byte, error) {
	if e.Type == "" {
		e.Type = "error"
	}

	// Error response structure
	val := struct {
		Type    string      `json:"type"`
		Title   string      `json:"title,omitempty"`
		Message interface{} `json:"message"`
		Errors  interface{} `json:"errors,omitempty"`
	}{Type: e.Type, Title: e.Title, Message: e.Message,
		Errors: e.Errors}

	return json.Marshal(val)
}

type HTTPErr struct {
	Title       string `json:"title"`
	Status      int    `json:"status"`
	Message     ErrMessage
	HTTPHeaders http.Header
	Ms          []interface{}
}

func (e *HTTPErr) MarshalJSON() ([]byte, error) {
	val := struct {
		Type    string      `json:"type"`
		Title   string      `json:"title"`
		Status  int         `json:"status"`
		Message interface{} `json:"message"`
	}{Type: "error", Title: e.Title, Status: e.Status, Message: e.Ms}

	return json.Marshal(val)
}

type ErrMessage struct {
	Key    string `json:"key"`
	Detail string `json:"detail"`
}

func (e *HTTPErr) Error() string {
	return e.Message.Detail
}

func (e *HTTPErr) Headers() http.Header {
	return e.HTTPHeaders
}

func (e *HTTPErr) StatusCode() int {
	if e.Status > 0 {
		return e.Status
	}
	return http.StatusOK

}
