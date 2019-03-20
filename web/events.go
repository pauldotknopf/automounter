package web

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type eventStruct struct {
	EventType string      `json:"eventType"`
	Data      interface{} `json:"data"`
}

var upgrader = websocket.Upgrader{}

func (server *Server) events(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	addedChannel, addedChannelCancel := server.mediaProvider.MediaAddded()
	removedChannel, removedChannelCancel := server.mediaProvider.MediaRemoved()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for media := range addedChannel {
			c.WriteJSON(eventStruct{"mediaAdded", convertMediaToJSON(media)})
		}
	}()
	go func() {
		defer wg.Done()
		for media := range removedChannel {
			c.WriteJSON(eventStruct{"mediaRemoved", media})
		}
	}()

	c.ReadMessage()

	addedChannelCancel()
	removedChannelCancel()

	wg.Wait()
}
