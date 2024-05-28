package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var messageQueue = make(chan interface{}, 1)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WsHandler(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	err = wsMessageWriter(ws)
	if err != nil {
		log.Println(err)
	}
}

func AddMessageToQueue(msg interface{}) {
	messageQueue <- msg
}

func wsMessageWriter(ws *websocket.Conn) error {

	for msg := range messageQueue {

		err := ws.WriteJSON(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
