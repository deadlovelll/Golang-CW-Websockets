package postmessengercontroller

import (
	"fmt"

	"github.com/gorilla/websocket"

	Messages "messenger_engine/models/message"
	MessageController "messenger_engine/controllers/message_controller"
)

var (
	Clients          = make(map[*websocket.Conn]bool)
	Broadcast        = make(chan Messages.FinalMessage)
	RepliesBroadcast = make(chan Messages.FinalMessageReply)
)

// HandleMessages handles broadcasting messages to all connected clients.
func HandleMessages(mmc *MessageController.MessageController) {
	for {
		msg := <-Broadcast
		for client := range Clients {
			if err := client.WriteJSON(msg); err != nil {
				handleClientError(client, err)
			}
		}
	}
}

// handleClientError handles errors by closing the client connection and cleaning it up.
func handleClientError(client *websocket.Conn, err error) {
	fmt.Printf("Error sending message: %v\n", err)
	client.Close()
	delete(Clients, client)
}
