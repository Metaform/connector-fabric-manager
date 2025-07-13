package model

import (
	"errors"
	"fmt"
)

type RecoverableError interface {
	error
	IsRecoverable() bool
}

type ClientError interface {
	error
	IsClientError() bool
}

type FatalError interface {
	error
	IsFatal() bool
}

type GeneralRecoverableError struct {
	Message    string
	Cause      error
	badRequest bool
}

func (e GeneralRecoverableError) Error() string       { return e.Message }
func (e GeneralRecoverableError) Unwrap() error       { return e.Cause }
func (e GeneralRecoverableError) IsRecoverable() bool { return true }

type BadRequestError struct {
	Message string
	Cause   error
}

func (e BadRequestError) Error() string       { return e.Message }
func (e BadRequestError) IsClientError() bool { return true }
func (e BadRequestError) Unwrap() error       { return e.Cause }

type SystemError struct {
	Message string
	Cause   error
}

func (e SystemError) Error() string { return e.Message }
func (e SystemError) IsFatal() bool { return true }
func (e SystemError) Unwrap() error { return e.Cause }

func NewRecoverableError(message string, args ...any) error {
	return GeneralRecoverableError{Message: fmt.Sprintf(message, args...)}
}

func NewClientError(message string, args ...any) error {
	return BadRequestError{Message: fmt.Sprintf(message, args...)}
}

func NewFatalError(message string, args ...any) error {
	return SystemError{Message: fmt.Sprintf(message, args...)}
}

func NewRecoverableWrappedError(cause error, message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	if cause != nil {
		formattedMessage = fmt.Sprintf("%s: %s", formattedMessage, cause.Error())
	}
	return GeneralRecoverableError{
		Message: formattedMessage,
		Cause:   cause,
	}
}

func NewClientWrappedError(cause error, message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	if cause != nil {
		formattedMessage = fmt.Sprintf("%s: %s", formattedMessage, cause.Error())
	}
	return BadRequestError{
		Message: formattedMessage,
		Cause:   cause,
	}
}

func NewFatalWrappedError(cause error, message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	if cause != nil {
		formattedMessage = fmt.Sprintf("%s: %s", formattedMessage, cause.Error())
	}
	return SystemError{
		Message: formattedMessage,
		Cause:   cause,
	}
}

func IsRecoverable(err error) bool {
	var recErr RecoverableError
	return errors.As(err, &recErr) && recErr.IsRecoverable()
}

func IsClientError(err error) bool {
	var clientErr ClientError
	return errors.As(err, &clientErr) && clientErr.IsClientError()
}

func IsFatal(err error) bool {
	var fatalErr FatalError
	return errors.As(err, &fatalErr) && fatalErr.IsFatal()
}
