package core

import "errors"

func Err[T error]() error {
	var err T
	return errors.Join(err, NewStacktrace())
}

type ErrCanceled struct{}

func (e ErrCanceled) Error() string {
	return "canceled"
}
