package streamService


import (
	"github.com/google/uuid"
	"sync"

	exception "github.com/Lmare/lightning-test/backend/exception"
)

type Stream interface {
    Recv() (any, error)
    Close() error
}

// gereric structure for the stream
type StreamWrapper[T any] struct {
    RecvCallback func() (*T, error)
    CloseCallback func() error
}
func (s StreamWrapper[T]) Recv() (any, error) {
    return s.RecvCallback()
}

func (s StreamWrapper[T]) Close() error {
    return s.CloseCallback()
}


// save server stream
// map[string]StreamWrapper
// TODO have une struct Envelop to have a batch garbage Collector in case
var streamMap = sync.Map{}

// save the steam in context of the server
func KeepStream[T any](stream StreamWrapper[T]) string {
	id := uuid.New().String()
	streamMap.Store(id, stream)
	return id
}

// Restore le stream in the actual thread context
func RestoreSteamByUuid(id string) (Stream, error) {
	val, ok := streamMap.Load(id)
	if !ok {
		return StreamWrapper[any]{}, exception.NewError("Stream don't exist", nil, exception.NewExampleError)
	}
	streamMap.Delete(id)

	stream, ok := val.(Stream)
	if !ok {
	    return StreamWrapper[any]{}, exception.NewError("Impossible to cast", nil, exception.NewExampleError)
	}
	return stream, nil
}
