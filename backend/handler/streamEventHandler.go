package handler


import (
	"fmt"
	"net/http"
	"io"

	streamService "github.com/Lmare/lightning-test/backend/service/streamService"
)

// check the message from gRPC stream
func handleStreamEvent(response http.ResponseWriter, request *http.Request) {

	streamId := request.FormValue("streamId")
	stream, err := streamService.RestoreSteamByUuid(streamId)
	if err != nil {
		fail(response, request, "Impossible to access to the events", err)
		return
	}
	/*
	response.Header().Set("Content-Type", "text/event-stream")
	response.Header().Set("Cache-Control", "no-cache")
	response.Header().Set("Connection", "keep-alive")
	*/
	for {
		msg, err := stream.Recv()
    	if err == io.EOF {
        	break // stream terminé
    	}
    	if err != nil {
        	fmt.Println("Erreur stream:", err)
        	break
    	}
		fmt.Println("message : ", msg)
		/*
		// Push SSE
    	fmt.Fprintf(w, "data: %s\n\n", encode(msg))
    	w.(http.Flusher).Flush()
		*/
	}

}
