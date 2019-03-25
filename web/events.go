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

	var lock sync.Mutex
	var doLock = func() {
		log.Println("doing lock")
		lock.Lock()
		log.Println("dot lock")
	}
	var doUnlock = func() {
		log.Println("doing unlock")
		lock.Unlock()
		log.Println("got unlock")
	}

	addedChannel, addedChannelCancel := server.mediaProvider.MediaAddded()
	removedChannel, removedChannelCancel := server.mediaProvider.MediaRemoved()
	mediaMountedChannel, mediaMountedChannelCancel := server.mediaProvider.MediaMounted()
	mediaUnmounteChannel, mediaUnmountedChannelCancel := server.mediaProvider.MediaUnmounted()

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		for media := range addedChannel {
			doLock()
			c.WriteJSON(eventStruct{"mediaAdded", convertMediaToJSON(media)})
			doUnlock()
		}
	}()
	go func() {
		defer wg.Done()
		for media := range removedChannel {
			doLock()
			c.WriteJSON(eventStruct{"mediaRemoved", media})
			doUnlock()
		}
	}()
	go func() {
		defer wg.Done()
		for media := range mediaMountedChannel {
			doLock()
			c.WriteJSON(eventStruct{"mediaMounted", media})
			doUnlock()
		}
	}()
	go func() {
		defer wg.Done()
		for media := range mediaUnmounteChannel {
			doLock()
			c.WriteJSON(eventStruct{"mediaUnmounted", media})
			doUnlock()
		}
	}()

	c.ReadMessage()

	addedChannelCancel()
	removedChannelCancel()
	mediaMountedChannelCancel()
	mediaUnmountedChannelCancel()

	wg.Wait()
}
