package errs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// UnAuthorisedErr error response object
func UnAuthorisedErr(message string) error {
	return ErrRes(UnathorizedErrorTitle,
		http.StatusUnauthorized,
		message)
}

// ForbiddenErr error response object
func ForbiddenErr(message string) error {
	return ErrRes(UnathorizedErrorTitle,
		http.StatusForbidden,
		message)
}

// RequestNotProcessed error response object
func RequestNotProcessed(message interface{}) error {
	return ErrRes(UnprocessableEntityMessage,
		http.StatusUnprocessableEntity,
		message)
}

// RequestRatelimitExceeded error response object
func RequestRatelimitExceeded(message interface{}) error {
	return ErrRes(TooManyRequests,
		http.StatusTooManyRequests,
		message)
}

// Internal error with message response object
func InternalErrWithMsg(message string) error {
	return ErrRes(InternalServerErrorTitle,
		http.StatusInternalServerError,
		message)
}

// Internal error response object
func InternalErr() error {
	return ErrRes(InternalServerErrorTitle,
		http.StatusInternalServerError,
		InternalServerErrorMessage)
}

// BadRequest error response object
func BadRequest(message interface{}) error {
	return ErrRes(BadRequestTitle,
		http.StatusBadRequest,
		message)
}

// This function handle all the body payload error messages
func ErrorReqHandler(err error) error {

	var (
		syntaxError        *json.SyntaxError
		unmarshalTypeError *json.UnmarshalTypeError
	)

	errMsg := make([]interface{}, 0)

	switch {
	case errors.As(err, &syntaxError):
		errMsg = append(errMsg, ErrMessage{
			Key:    "SyntaxError",
			Detail: SyntaxErrorMessageDetatil})

		return ErrResponse(BadRequestTitle,
			http.StatusBadRequest, errMsg)

	case errors.Is(err, io.ErrUnexpectedEOF):
		errMsg = append(errMsg, ErrMessage{
			Key:    "SyntaxError",
			Detail: SyntaxErrorMessageDetatil})

		return ErrResponse(BadRequestTitle,
			http.StatusBadRequest, errMsg)

	case errors.As(err, &unmarshalTypeError):
		detail := fmt.Sprintf(InputErrorMessageDetatil,
			unmarshalTypeError.Field)

		errMsg = append(errMsg, ErrMessage{
			Key:    "SyntaxError",
			Detail: detail})

		return ErrResponse(BadRequestTitle,
			http.StatusBadRequest, errMsg)

	case strings.Contains(err.Error(), "Error:Field validation for "):
		errMsg = append(errMsg, ErrMessage{
			Key:    "MissingField",
			Detail: MissingFieldErrorMessageDetail})

		return ErrResponse(BadRequestTitle,
			http.StatusBadRequest, errMsg)

	case errors.Is(err, io.EOF):
		errMsg = append(errMsg, ErrMessage{
			Key:    "ExceedLimit",
			Detail: BodyPayloadLimitErrorMessageDetail})

		return ErrResponse(BadRequestTitle,
			http.StatusBadRequest, errMsg)

	default:
		errMsg = append(errMsg, ErrMessage{
			Key:    "InternalServer",
			Detail: BadRequestErrorMessageDetail})

		return ErrResponse(BadRequestTitle,
			http.StatusBadRequest, errMsg)
	}
}

func ErrResponse(title string, status int,
	message []interface{}) error {

	return &HTTPErr{Title: title,
		Status: status,
		Ms:     message}
}
