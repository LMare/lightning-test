package handler

import (
	"fmt"
	"net/http"

	streamService "github.com/Lmare/lightning-playground/backend/service/streamService"
)

// check the message from gRPC stream
func handleStreamEvent(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "text/event-stream")
	response.Header().Set("Cache-Control", "no-cache")
	response.Header().Set("Connection", "keep-alive")

	fmt.Fprintf(response, "event: init\ndata: initialisation de la SSE\n\n")
	response.(http.Flusher).Flush()

	notify := request.Context().Done()
	// inscription au flux de notification
	streamService.SubscribeSse(response)

	// wait until the connection is interrupted
	<-notify
	fmt.Println("Client SSE déconnecté")
	streamService.RevoqueSse(response)
}
