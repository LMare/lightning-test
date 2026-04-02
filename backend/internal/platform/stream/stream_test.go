package stream


// -----------------------------------------------------------------------------
// |/////////////////////////// istream \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
// -----------------------------------------------------------------------------


import (
    "bytes"
	"encoding/json"
	"errors"
	"fmt"
    "io"
    "net/http"
    "strings"
    "testing"
    "time"
)

// --- Fakes for dependencies ---

// fakeWriter simulates an http.ResponseWriter + Flusher
type fakeWriter struct {
    buf     bytes.Buffer
    flushed bool
}

func (f *fakeWriter) Header() http.Header { return http.Header{} }
func (f *fakeWriter) Write(b []byte) (int, error) { return f.buf.Write(b) }
func (f *fakeWriter) WriteHeader(statusCode int)  {}
func (f *fakeWriter) Flush()                      { f.flushed = true }

// fakeStream simule un istream
type fakeStream struct {
    messages []any
    errors   []error
    index    int
    closed   bool
}

func (fs *fakeStream) Recv() (any, error) {
    if fs.index < len(fs.messages) {
        msg := fs.messages[fs.index]
        fs.index++
        return msg, nil
    }
    if fs.index-len(fs.messages) < len(fs.errors) {
        err := fs.errors[fs.index-len(fs.messages)]
        fs.index++
        return nil, err
    }
    return nil, io.EOF
}

func (fs *fakeStream) Close() error {
    fs.closed = true
    return nil
}

// --- Unit tests ---

func TestSessionAddRemove(t *testing.T) {
    s := &session{sseList: []http.ResponseWriter{}}
    fw := &fakeWriter{}

    s.add(fw)
    if len(s.sseList) != 1 {
        t.Errorf("expected 1, got %d", len(s.sseList))
    }

    s.remove(fw)
    if len(s.sseList) != 0 {
        t.Errorf("expected 0, got %d", len(s.sseList))
    }
}

func TestSessionNotifyAll(t *testing.T) {
    s := &session{sseList: []http.ResponseWriter{}}
    fw := &fakeWriter{}
    s.add(fw)

    msg := "hello\nworld"
    s.notifyAll(msg)

    out := fw.buf.String()
    if !strings.Contains(out, "data: hello world") {
        t.Errorf("expected formatted message, got %q", out)
    }
    if !fw.flushed {
        t.Errorf("expected flush to be called")
    }
}

func TestSessionStart(t *testing.T) {
    s := &session{channel: make(chan string, 1), sseList: []http.ResponseWriter{}}
    fw := &fakeWriter{}
    s.add(fw)

    s.start()
    s.channel <- "ping"

    // give the goroutine time to process the message
    time.Sleep(50 * time.Millisecond)

    out := fw.buf.String()
    if !strings.Contains(out, "data: ping") {
        t.Errorf("expected message 'ping', got %q", out)
    }
}

func TestSessionStreamMessages(t *testing.T) {
    s := &session{channel: make(chan string, 10)}
    fs := &fakeStream{messages: []any{"msg1", "msg2"}}

    s.stream(fs)

    // wait for goroutines to send messages
    time.Sleep(50 * time.Millisecond)

    select {
    case m1 := <-s.channel:
        if !strings.Contains(m1, "msg1") && !strings.Contains(m1, "msg2") {
            t.Errorf("unexpected message %q", m1)
        }
    default:
        t.Errorf("expected a message in channel")
    }
}

func TestSessionStreamError(t *testing.T) {
    s := &session{channel: make(chan string, 10)}
    fs := &fakeStream{errors: []error{errors.New("boom")}}

    s.stream(fs)

    time.Sleep(50 * time.Millisecond)

    select {
    case m := <-s.channel:
        if !strings.Contains(m, "Error: boom") {
            t.Errorf("expected error message, got %q", m)
        }
    default:
        t.Errorf("expected an error message in channel")
    }
}


// -----------------------------------------------------------------------------
// |/////////////////////////// end istream \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
// -----------------------------------------------------------------------------




func TestGetSession_New(t *testing.T) {
    id := "abc"
    s := getSession(id)
    if s == nil {
        t.Fatal("expected a session, got nil")
    }
}

func TestGetSession_Existing(t *testing.T) {
    id := "xyz"
    s1 := getSession(id)
    s2 := getSession(id)
    if s1 != s2 {
        t.Fatal("expected same session instance")
    }
}

func TestSubscribeAndRevoqueSse(t *testing.T) {
    fw := &fakeWriter{}

    SubscribeSse(fw)
    // verify the session holds the writer
    s, _ := sessions.Load("uniqueSession")
    sess := s.(isession)
    if len(sess.(*session).sseList) != 1 {
        t.Errorf("expected 1 SSE client, got %d", len(sess.(*session).sseList))
    }

    RevoqueSse(fw)
    if len(sess.(*session).sseList) != 0 {
        t.Errorf("expected 0 SSE client, got %d", len(sess.(*session).sseList))
    }
}

/*
func TestStreamResult(t *testing.T) {
    msgs := []string{"hello"}
    idx := 0
    sw := StreamWrapper[string]{
        RecvCallback: func() (*string, error) {
            if idx < len(msgs) {
                m := msgs[idx]
                idx++
                return &m, nil
            }
            return nil, io.EOF
        },
        CloseCallback: func() error { return nil },
    }

    StreamResult(sw)

    // wait for the goroutine to push the message
    time.Sleep(50 * time.Millisecond)

    s, ok := sessions.Load("uniqueSession")
    if !ok {
        t.Fatal("expected session to be created")
    }
    sess := s.(isession)

    select {
    case got := <-sess.(*session).channel:
        if got != "hello" {
            t.Errorf("expected 'hello', got %q", got)
        }
    default:
        t.Errorf("expected a message in channel")
    }
}
*/


func TestEncodeBasicTypes(t *testing.T) {
    if got := encode("abc"); got != "abc" {
        t.Errorf("expected 'abc', got %q", got)
    }
    if got := encode(42); got != "42" {
        t.Errorf("expected '42', got %q", got)
    }
    if got := encode(true); got != "true" {
        t.Errorf("expected 'true', got %q", got)
    }
}

type myType int

func (m myType) String() string {
    return fmt.Sprintf("val=%d", m)
}

func TestEncodeStringer(t *testing.T) {
    if got := encode(myType(7)); got != "val=7" {
        t.Errorf("expected 'val=7', got %q", got)
    }
}

func TestEncodeStructToJSON(t *testing.T) {
    type user struct {
        Name string
        Age  int
    }
    u := user{"Alice", 30}
    got := encode(u)

    var decoded user
    if err := json.Unmarshal([]byte(got), &decoded); err != nil {
        t.Fatalf("expected valid JSON, got error %v", err)
    }
    if decoded.Name != "Alice" || decoded.Age != 30 {
        t.Errorf("unexpected JSON content: %+v", decoded)
    }
}

func TestEncodeErrorOnMarshal(t *testing.T) {
    // reflect.ValueOf(chan int) is not JSON-serializable
    ch := make(chan int)
    got := encode(ch)
    if !strings.Contains(got, "error:") {
        t.Errorf("expected error message, got %q", got)
    }
}
