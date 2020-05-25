package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const addr = ":8764"

// Message represents single websocket message
type Message struct {
	HandleName string `json:"handleName"`
	MsgContent string `json:"msgContent"`
	Timestamp  int64  `json:"timestamp"`
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // channel for broadcasting messages to all connected clients

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {

	// handle main page
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// handle websocket connections
	http.HandleFunc("/websocket", handleWebsocketConn)

	// prepare service for broadcasting messages
	go broadcastMessages()

	// start server
	log.Printf("server is listening on %v", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("unable to run server due: %v", err)
	}
}

func handleWebsocketConn(w http.ResponseWriter, r *http.Request) {

	// attempt to upgrade http connection to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if checkError(err) {
		return
	}
	defer ws.Close()

	// register client as active client
	clients[ws] = true

	for {
		var msg Message

		// read in new message as JSON and map it to Message object
		err := ws.ReadJSON(&msg)
		if checkError(err) {
			delete(clients, ws)
			break
		}

		// send newly receive message to broadcast channel
		broadcast <- msg
	}
}

func broadcastMessages() {
	for {
		// get message from channel
		msg := <-broadcast

		// send it to every active client
		for client := range clients {
			err := client.WriteJSON(msg)
			if checkError(err) {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// define short function for checking error
func checkError(err error) bool {
	var isError bool = (err != nil)
	if isError {
		log.Println(err)
	}
	return isError
}
