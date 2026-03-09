package streamService

import (
	//"github.com/google/uuid"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"

	//"time"
	lnrpc "github.com/Lmare/lightning-playground/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc"
	//exception "github.com/Lmare/lightning-playground/backend/exception"
)

func NewStreamService() *StreamService {
	return &StreamService{}
}

type StreamService struct{} // TODO

type istream interface {
	Recv() (any, error)
	Close() error
}

// gereric structure for the stream
type StreamWrapper[T any] struct {
	RecvCallback  func() (*T, error)
	CloseCallback func() error
}

func (s StreamWrapper[T]) Recv() (any, error) {
	return s.RecvCallback()
}

func (s StreamWrapper[T]) Close() error {
	return s.CloseCallback()
}

// --------------------------------------------

type isession interface {
	add(http.ResponseWriter)
	remove(http.ResponseWriter)
	notifyAll(string)
	start()
	stream(istream)
}

type session struct {
	channel   chan string
	muSseList sync.Mutex
	sseList   []http.ResponseWriter
}

// add a sse
func (s *session) add(sse http.ResponseWriter) {
	s.muSseList.Lock()
	s.sseList = append(s.sseList, sse)
	s.muSseList.Unlock()
}

// remove a sse
func (s *session) remove(sse http.ResponseWriter) {
	s.muSseList.Lock()
	for i, w := range s.sseList {
		if w == sse {
			s.sseList = append(s.sseList[:i], s.sseList[i+1:]...)
			break
		}
	}
	s.muSseList.Unlock()
}

// send an event in all SSE
func (s *session) notifyAll(msg string) { // TODO : prévoir un type ?
	s.muSseList.Lock()
	for _, w := range s.sseList {
		fmt.Fprintf(w, "data: %s\n\n", strings.ReplaceAll(msg, "\n", " "))
		w.(http.Flusher).Flush()
	}
	s.muSseList.Unlock()
}

// listen incoming message in the channel and notify all the sse clients
func (s *session) start() {
	go func() {
		for {
			select {
			case msg := <-s.channel:
				// Push SSE
				fmt.Printf("Message brut envoyé : %#v\n", strings.ReplaceAll(msg, "\n", " "))
				s.notifyAll(msg)
			}
		}
	}()
}

// stream a ressource into the channel of the session
func (se *session) stream(st istream) {
	go func() {
		for {
			msg, err := st.Recv()
			if err == io.EOF {
				fmt.Println("fin de la goRoutine")
				break // stream terminé
			} else if err != nil {
				fmt.Println("Erreur sur le stream", err)
				se.channel <- fmt.Sprintf("Erreur : %s", err)
				break
			} else {
				fmt.Println("Data", msg)
				se.channel <- encode(msg)
			}
		}
	}()
}

// --------------------------------------------------------------

// channel for the session
// map[string]*session
var sessions = sync.Map{}

func SubscribeSse(sse http.ResponseWriter) {
	id := "uniqueSession"
	session := getSession(id)
	session.add(sse)
}

func RevoqueSse(sse http.ResponseWriter) {
	id := "uniqueSession"
	session := getSession(id)
	session.remove(sse)
}

/** TODO:
🛠️ Points d’attention
- Utilise des canaux bufferisés (make(chan Event, N)) pour éviter de bloquer les producteurs si le consommateur est lent.
- Ajoute un ping/keep‑alive régulier pour maintenir la connexion ouverte (et éviter que des proxies la coupent).
- Surveille la taille des listes de clients pour éviter les fuites mémoire si un utilisateur ouvre/ferme beaucoup d’onglets.
- Garbage collector pour supprimer les sessions si aucun client existant



donc du coup avec cette configuration là j'ai les notifications sur tous les onglets.
je me suis dit que si je veux des nofitications qui s'affiche uniquement sur certain onglets je peux faire ça : (note j'utilise HTMX, mais on pourrait avoir plus ou moin la même logique en rest classique)
dans le fait une action qui produit un stream gRPC, je génère un uuid que je met dans mon StreamWrapper,
je retourne au navigateur du html qui defini une class css qui dépends de cett uuid qui fait un display block.
dans les event SSE je génère une envelopper HTML sur le message qui ajoute une classe pour mettre les notif mono-onglet en display none + la classe unique qui permet d'afficher seulement dans l'onglet qui contient la définition.
*/

func getSession(sessionId string) isession {
	s, ok := sessions.Load(sessionId)
	if !ok {
		fmt.Println("initialisation de la session")
		s2 := &session{channel: make(chan string), sseList: make([]http.ResponseWriter, 0)}
		sessions.Store(sessionId, s2)
		s2.start()
		return s2
	}
	return s.(isession)
}

// save the steam in context of the server
func StreamResult[T any](stream StreamWrapper[T]) {
	//id := uuid.New().String()
	id := "uniqueSession"
	session := getSession(id)
	session.stream(stream)
}

// ------

// encode transforme n'importe quelle valeur en string pour SSE
func encode(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case *lnrpc.Payment:
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
