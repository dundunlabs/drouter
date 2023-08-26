package prenn

import "net/http"

type Exception struct {
	statusCode int
	err        error
}

func (e Exception) Error() string {
	if e.err == nil {
		return http.StatusText(e.statusCode)
	}
	return e.err.Error()
}

func (e Exception) WithError(err error) Exception {
	exception := e
	exception.err = err
	return exception
}

var (
	ExceptionBadRequest Exception = Exception{statusCode: http.StatusBadRequest}
)
