package main

import (
	"main/internals/producer"
	ws "main/internals/websocket"
	"net/http"
	"time"
)

func main() {

	go producer.ProduceDeviationData(time.Microsecond, 3000, 5)

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", ws.WsHandler)
	http.ListenAndServe("localhost:5055", nil)
}
