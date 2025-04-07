package constant

import (
	"fmt"
)

type CommonError struct {
	code      int64
	msg       string
	errReason string
}

func NewErrorWithoutReason(code int, msg string) CommonError {
	return CommonError{code: int64(code), msg: msg}
}

func NewError(code int, msg, errReason string) CommonError {
	return CommonError{code: int64(code), msg: msg, errReason: errReason}
}

func (err CommonError) Error() string {
	return fmt.Sprintf("Error: [%d] %s, %s", err.code, err.msg, err.errReason)
}

func (err CommonError) Code() int64 {
	return err.code
}

func (err CommonError) Msg() string {
	return err.msg
}
