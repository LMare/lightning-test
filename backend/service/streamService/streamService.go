package streamService


import (
	//"github.com/google/uuid"
	"sync"
	"io"
	"fmt"
	"reflect"
	"encoding/json"
	lnrpc "github.com/Lmare/lightning-test/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc"


	//exception "github.com/Lmare/lightning-test/backend/exception"
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
var sessionChannelMap = sync.Map{}

func GetChannel(sessionId string) chan string{
	channel, ok := sessionChannelMap.Load(sessionId)
	if !ok {
		channel = make(chan string)
		sessionChannelMap.Store(sessionId, channel)
	}
	return channel.(chan string)
}

// save the steam in context of the server
func StreamResult[T any](stream StreamWrapper[T]) {
	//id := uuid.New().String()
	id := "uniqueSession"
	channel := GetChannel(id)
	go func() {
        for {
            msg, err := stream.Recv()
			if err == io.EOF {
				break // stream terminé
			}
            if err != nil {
				fmt.Println("Erreur sur le stream", err)
                channel <-fmt.Sprintf(" Erreur : %s", err)
				break
            }
			fmt.Println("Stream Data", msg)
            channel <- encode(msg)
        }
    }()
}

// ------

// encode transforme n'importe quelle valeur en string pour SSE
func encode(v interface{}) string {
    switch val := v.(type) {
    case string:
        return val
	case *lnrpc.Payment :
		return fmt.Sprintf("💸 Paiement de %d sats — statut : %s", val.ValueSat, val.Status.String())
    case fmt.Stringer:
        return val.String()
    default:
        // Si c'est un type simple (int, float, bool, etc.)
        rv := reflect.ValueOf(v)
        switch rv.Kind() {
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
            reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
            reflect.Float32, reflect.Float64, reflect.Bool:
            return fmt.Sprintf("%v", v)
        default:
            // Pour les structs, slices, maps, etc. → JSON
            jsonData, err := json.Marshal(v)
            if err != nil {
                return fmt.Sprintf("error: %v", err)
            }
            return string(jsonData)
        }
    }
}
