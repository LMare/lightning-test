package exception

// type only for the example and the syntaxe to do errors custom by enriching AppErrorBase

type ExampleErrorDetail struct {
	base BaseError
	Detail string
}

func (e *ExampleErrorDetail) Error() string {
    return e.base.Error()
}

func (e *ExampleErrorDetail) Unwrap() error {
	return e.base.Unwrap();
}

func (e *ExampleErrorDetail) File() string {
	return e.base.File();
}

func (e *ExampleErrorDetail) Line() int {
	return e.base.Line();
}

func (e *ExampleErrorDetail) Message() string {
	return e.base.Message();
}

// Constructor for the exemple
// Instanciate like :
//   err := NewError("validation échouée", nil, ExampleErrorDetail, "banana is not a vegetable")
func NewExampleErrorDetail(m string, f string, l int, c error,  args ...interface{}) *ExampleErrorDetail {
    if len(args) > 1 {
        panic("trop d'arguments")
    }

    detail := ""
    if len(args) == 1 {
        if s, ok := args[0].(string); ok {
            detail = s
        } else {
            panic("argument de détail invalide")
        }
    }

    return &ExampleErrorDetail{
        base: NewBaseErrorImpl(m, f, l, c),
        Detail: detail,
    }
}
