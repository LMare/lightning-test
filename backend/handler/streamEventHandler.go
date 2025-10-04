package handler


import (
	"fmt"
	"net/http"


	streamService "github.com/Lmare/lightning-test/backend/service/streamService"
)

// check the message from gRPC stream
func handleStreamEvent(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "text/event-stream")
	response.Header().Set("Cache-Control", "no-cache")
	response.Header().Set("Connection", "keep-alive")

	id := "uniqueSession"
	channel := streamService.GetChannel(id)
	/*
	if err != nil {
		//fail(response, request, "Impossible to access to the events", err)
		fmt.Fprintf(response, "event: error\ndata: %s\n\n", encode("Une erreur est survenue"))
		fmt.Fprintf(response, "event: end\ndata: done\n\n")
		http.Error(response, "stream expired", http.StatusGone) // 410 Gone
		response.(http.Flusher).Flush()
		return
	}*/

	for {
		select {
		case msg := <-channel :
			// Push SSE
			fmt.Println("Data Flush")
			fmt.Fprintf(response, "data: %s\n\n", msg)
			response.(http.Flusher).Flush()
		}
	}
/*
		go func(stream Stream) {

			for {
				msg, err := stream.Recv()
		    	if err == io.EOF {
		        	break // stream terminé
		    	}
		    	if err != nil {
		        	fmt.Println("Erreur stream:", err)
					fmt.Fprintf(response, "event: error\ndata: %s\n\n", encode(err))
					response.(http.Flusher).Flush()
		        	break
		    	}
				fmt.Println("message : ", msg)

				// Push SSE
		    	fmt.Fprintf(response, "data: %s\n\n", encode(msg))
		    	response.(http.Flusher).Flush()

			}
			fmt.Fprintf(response, "event: end\ndata: done\n\n")
			response.(http.Flusher).Flush()

		}(stream)
	}*/
}
