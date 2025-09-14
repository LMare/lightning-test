package exception

import (
	"fmt"
	"runtime"
	"strings"
)
// brique de base pour faire des erreurs facile à identifier
type BaseError interface {
    error
    Unwrap() error
    File() string
    Line() int
    Message() string
}

type BaseErrorImpl struct {
    message string
    file    string
    line    int
    cause   error
}

// => implementation de l'interface error
func (e *BaseErrorImpl) Error() string {
	if e.cause != nil {
        return fmt.Sprintf("(%s:%d) %s → %v", e.file, e.line, e.message, e.cause.Error())
    }
    return fmt.Sprintf("(%s:%d) %s", e.file, e.line, e.message)
}

func (e *BaseErrorImpl) Unwrap() error {
	return e.cause;
}

func (e *BaseErrorImpl) File() string {
	return e.file;
}

func (e *BaseErrorImpl) Line() int {
	return e.line;
}

func (e *BaseErrorImpl) Message() string {
	return e.message;
}


// Base path of th project to remove it in the path of the file
var path string
func ConfigureProjectBasePath(p string) {
	path = p
}


// generic constructor to manage args at the runtime
func NewError[T BaseError](msg string, cause error, ctor func(string, string, int, error, ...interface{}) T, args ...interface{}) T {
    _, file, line, _ := runtime.Caller(1)
	file = strings.TrimPrefix(file, path)
    return ctor(msg, file, line, cause, args...)
}

// constructor used by custom exeption
func NewBaseErrorImpl(m string, f string, l int, c error) BaseError {
	return &BaseErrorImpl{
		message: m,
	    file: f,
	    line: l,
	    cause: c,
	}
}
