package exception

import "fmt"

// brique de base pour faire des erreurs facile à identifier
type BaseError interface {
    error
    Unwrap() error
	// à voir si c'est pas trop lourd
    File() string
    Line() int
    Message() string
}

type BaseErrorImpl struct {
    Message string
    File    string
    Line    int
    Cause   error
}

// => implementation de l'interface error
func (e *BaseErrorImpl) Error() string {
	if e.Cause != nil {
        return fmt.Sprintf("%s (%s:%d) → %v ", e.Message, e.File, e.Line, e.Cause.Error())
    }
    return fmt.Sprintf("%s (%s:%d)", e.Message, e.File, e.Line)
}

func (e *BaseErrorImpl) Unwrap() error {
	return e.Cause;
}

func (e *BaseErrorImpl) File() string {
	return e.File;
}

func (e *BaseErrorImpl) Line() int {
	return e.Line;
}

func (e *BaseErrorImpl) Message() string {
	return e.Message;
}

// generic constructor to manage args at the runtime
func NewError[T BaseError](msg string, cause error, ctor func(msg string, cause error, file string, line int, args ...interface{}) T, args ...interface{}) T {
    _, file, line, _ := runtime.Caller(1)
    return ctor(msg, cause, file, line, args...)
}

// constructor used by custom exeption
func NewBaseErrorImpl(m string, f string, l int, c error) *BaseErrorImpl {
	return &BaseErrorImpl{
		Message: m,
	    File: f,
	    Line: l,
	    Cause: c,
	}
}
