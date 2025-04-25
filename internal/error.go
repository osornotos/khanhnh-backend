package internal

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type ApplicationError struct {
	HTTPStatusCode int    `json:"statusCode,omitempty"`
	ErrorCode      string `json:"code,omitempty"`
	ErrorText      string `json:"message,omitempty"`
	// e.g. external error when validate error
	ErrorMetadata interface{} `json:"error,omitempty"`
	// for sentry, not client
	ErrorStack error `json:"-"`
}

type ValidatorError struct {
	Key     string `json:"key,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (err ApplicationError) Error() string {
	if err.ErrorMetadata != nil {
		b, _ := json.Marshal(err.ErrorMetadata)
		return fmt.Sprintf("status: %d, code: %s, message: %s, error: %s", err.HTTPStatusCode, err.ErrorCode, err.ErrorText, string(b))
	}
	return fmt.Sprintf("status: %d, code: %s, message: %s", err.HTTPStatusCode, err.ErrorCode, err.ErrorText)
}

func (err ApplicationError) WithMessage(message string) ApplicationError {
	err.ErrorText = message
	return err
}

func (err ApplicationError) WithMessagef(format string, v ...interface{}) ApplicationError {
	err.ErrorText = fmt.Sprintf(format, v...)
	return err
}

func (err ApplicationError) WithMetadata(errs interface{}) ApplicationError {
	err.ErrorMetadata = errs
	return err
}

// WithError separate client response and sentry reports
//
//	err = errors.New("to sentry")
//	return ErrorInternalServerError.WithMessage("to client").WithError(err)
func (err ApplicationError) WithError(e error) ApplicationError {
	if !hasStacktrace(e) {
		e = errors.WithStack(e)
	}
	err.ErrorStack = e
	return err
}

func (err ApplicationError) HasError() bool {
	return err.ErrorStack != nil
}

func (err ApplicationError) ErrorStackMessage() string {
	if err.ErrorStack == nil {
		return ""
	}
	return err.ErrorStack.Error()
}

func (err ApplicationError) IsServerError() bool {
	return err.HTTPStatusCode >= 500
}

func (err ApplicationError) GetDefaultErrorCode() string {
	if err.ErrorCode != "" {
		return err.ErrorCode
	}
	switch err.HTTPStatusCode {
	case 400:
		return "BAD_REQUEST"
	case 401:
		return "UNAUTHORIZED"
	case 403:
		return "FORBIDDEN"
	case 404:
		return "RESOURCE_NOT_FOUND"
	case 500:
		return "INTERNAL_SERVER_ERROR"
	default:
		return "UNKNOWN"
	}
}

func ApplicationErrorFromJson(body []byte, status int) ApplicationError {
	if status >= 500 {
		return ErrorInternalServerError.WithMessage(string(body))
	}
	appErr := ApplicationError{}
	err := json.Unmarshal(body, &appErr)
	if err != nil {
		return ErrorInternalServerError.WithMessage(fmt.Sprintf("body %s; err %s", string(body), err))
	}
	appErr.HTTPStatusCode = status
	return appErr
}

var (
	ErrorCodeValidatorMaxMin       = "INVALID_LENGTH"
	ErrorCodeValidatorRequire      = "FIELD_REQUIRED"
	ErrorCodeValidatorInvalidInput = "INVALID_INPUT"
)

var (
	ErrorInvalidInput        = ApplicationError{HTTPStatusCode: 400, ErrorCode: "INVALID_INPUT", ErrorText: "invalid input"}
	ErrorTimeout             = ApplicationError{HTTPStatusCode: 400, ErrorCode: "TIME_OUT", ErrorText: "time out"}
	ErrorBadRequest          = ApplicationError{HTTPStatusCode: 400, ErrorCode: "BAD_REQUEST", ErrorText: "bad request"}
	ErrorUnauthorized        = ApplicationError{HTTPStatusCode: 401, ErrorCode: "UNAUTHORIZED", ErrorText: "unauthorized"}
	ErrorForbidden           = ApplicationError{HTTPStatusCode: 403, ErrorCode: "FORBIDDEN", ErrorText: "forbidden"}
	ErrorNotFound            = ApplicationError{HTTPStatusCode: 404, ErrorCode: "RESOURCE_NOT_FOUND", ErrorText: "not found"}
	ErrorTooManyRequests     = ApplicationError{HTTPStatusCode: 429, ErrorCode: "TOO_MANY_REQUESTS", ErrorText: "too many requests"}
	ErrorInternalServerError = ApplicationError{HTTPStatusCode: 500, ErrorCode: "INTERNAL_SERVER_ERROR", ErrorText: "unknown error"}
)

func NewInvalidInputError(key string, code string, message string) ApplicationError {
	return ErrorInvalidInput.WithMessage("invalid form").WithMetadata([]ValidatorError{
		{Key: key, Code: code, Message: message},
	})
}

// type assertion "github.com/pkg/errors"
type stackTracer interface {
	StackTrace() errors.StackTrace
}

func hasStacktrace(err error) bool {
	_, ok := err.(stackTracer)
	return ok
}
