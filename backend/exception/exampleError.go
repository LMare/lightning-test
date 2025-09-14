package exception

// type only for the example and the syntaxe to do errors custom by enriching AppErrorBase

type ExampleError struct {
    base BaseError
}

func (e *ExampleError) Error() string {
    return e.base.Error()
}

func (e *ExampleError) Unwrap() error {
	return e.base.Unwrap();
}

func (e *ExampleError) File() string {
	return e.base.File();
}

func (e *ExampleError) Line() int {
	return e.base.Line();
}

func (e *ExampleError) Message() string {
	return e.base.Message();
}

// Constructor
// Instanciate like :
//   err := NewError("validation échouée", nil, NewExampleError)
func NewExampleError(m string, f string, l int, c error, args ...interface{}) *ExampleError {
	return &ExampleError {
		base: NewBaseErrorImpl(m, f, l, c),
	}
}
