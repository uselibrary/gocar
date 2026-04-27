package cli

import "errors"

type exitCodeError struct {
	code int
	err  error
}

func (e *exitCodeError) Error() string {
	return e.err.Error()
}

func (e *exitCodeError) Unwrap() error {
	return e.err
}

func (e *exitCodeError) ExitCode() int {
	if e.code <= 0 {
		return 1
	}
	return e.code
}

// WithExitCode wraps an error with a process exit code.
func WithExitCode(err error, code int) error {
	if err == nil {
		return nil
	}
	if code <= 0 {
		code = 1
	}
	return &exitCodeError{code: code, err: err}
}

// ExitCode extracts process exit code from error, defaulting to 1.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}

	var withCode interface{ ExitCode() int }
	if errors.As(err, &withCode) {
		code := withCode.ExitCode()
		if code > 0 {
			return code
		}
	}

	return 1
}
