package chatmessagehandler

import (
	"net/http"
	"log"

	"github.com/gorilla/websocket"

	Messages "messenger_engine/models/message"
	Broadcast "messenger_engine/controllers/broadcast_controller"
	MessageController "messenger_engine/controllers/message_controller"
	ErrorHandler "messenger_engine/controllers/websocket_controller/handlers/error_handler"
	MessageParser "messenger_engine/controllers/websocket_controller/parsers"
)

// ChatMessageHandler handles WebSocket connections for sending and receiving chat messages.
type ChatMessageHandler struct {
	upgrader      websocket.Upgrader // WebSocket upgrader for upgrading HTTP connection
	msgCtrl       *MessageController.MessageController // Controller for managing messages
	ErrorHandler  *ErrorHandler.ErrorHandler // Error handler for WebSocket errors
	MessageParser *MessageParser.Parser // Message parser for parsing incoming messages
	Broadcast     *Broadcast.Broadcast // Broadcast controller for broadcasting messages to clients
}

// NewChatMessageHandler initializes a new ChatMessageHandler with the given dependencies.
func NewChatMessageHandler(
	upgrader websocket.Upgrader, // WebSocket upgrader
	ctrl *MessageController.MessageController, // Message controller for handling message operations
	broadcast *Broadcast.Broadcast, // Broadcast controller to manage client broadcasts
) *ChatMessageHandler {
	if broadcast == nil {
		log.Fatal("NewChatMessageHandler: broadcast instance cannot be nil")
	}

	return &ChatMessageHandler{
		upgrader:      upgrader,
		msgCtrl:       ctrl,
		Broadcast:     broadcast,
		MessageParser: MessageParser.New(),
		ErrorHandler:  ErrorHandler.NewErrorHandler(),
	}
}

// ServeHTTP upgrades the HTTP connection to a WebSocket connection and handles incoming chat messages.
// It processes the WebSocket connection and delegates message handling to respective methods based on message type.
func (h *ChatMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    ws, err := h.upgrader.Upgrade(w, r, nil)
    if err != nil {
        // Handle WebSocket upgrade error
        h.ErrorHandler.HandleWebSocketError(err, nil, "Error upgrading to WebSocket: %s", err)
        return
    }
    defer ws.Close()

    // Register the client for broadcasts.
    h.Broadcast.Clients[ws] = true

    for {
        var msg map[string]interface{}
        if err := ws.ReadJSON(&msg); err != nil {
            // Handle error reading the message from WebSocket
            h.ErrorHandler.HandleWebSocketError(err, ws, "Error reading message: %s", err)
            delete(h.Broadcast.Clients, ws)
            break
        }

        // Handle different message types
        switch msg["type"] {
        case "initial":
            h.handleInitialMessage(ws, msg)
        case "message":
            h.handleMessage(ws, msg)
        case "message_reply":
            h.handleMessageReply(ws, msg)
        }
    }
}

// handleInitialMessage processes the initial message sent by the client. 
// It retrieves previous messages from the database and sends them back to the client.
func (h *ChatMessageHandler) handleInitialMessage(ws *websocket.Conn, msg map[string]interface{}) {
	chatID, err := h.MessageParser.ParseChatID(msg)
	if err != nil {
		// Handle error in parsing chat ID
		h.ErrorHandler.HandleWebSocketError(err, ws, "Invalid chat_id format: %s", err)
		return
	}

	messages, err := h.msgCtrl.LoadMessages(chatID)
	if err != nil {
		// Handle error loading messages
		h.ErrorHandler.HandleWebSocketError(err, ws, "Error loading messages: %s", err)
		return
	}

	// Send the initial messages back to the client
	if err := ws.WriteJSON(map[string]interface{}{"type": "initial", "messages": messages}); err != nil {
		// Handle error sending the messages
		h.ErrorHandler.HandleWebSocketError(err, ws, "Error sending initial messages: %s", err)
	}
}

// handleMessage processes a new message sent by the client. 
// It parses, saves the message to the database, and broadcasts it to other clients.
func (h *ChatMessageHandler) handleMessage(ws *websocket.Conn, msg map[string]interface{}) {
	messageData, err := h.MessageParser.ParseMessageData(msg)
	if err != nil {
		// Handle error in parsing message data
		h.ErrorHandler.HandleWebSocketError(err, ws, "Invalid message format: %s", err)
		return
	}

	if err := h.msgCtrl.SaveMessage(messageData); err != nil {
		// Handle error saving message to database
		h.ErrorHandler.HandleWebSocketError(err, ws, "Error saving message to database: %s", err)
		return
	}

	// Create a final message structure and broadcast to other clients
	finalMsg := Messages.FinalMessage{Type: "message", Message: messageData}
	h.Broadcast.Broadcast <- finalMsg
}

// handleMessageReply processes a message reply sent by the client. 
// It parses, saves the message reply to the database, and broadcasts it to other clients.
func (h *ChatMessageHandler) handleMessageReply(ws *websocket.Conn, msg map[string]interface{}) {
	messageReplyData, err := h.MessageParser.ParseMessageReplyData(msg)
	if err != nil {
		// Handle error in parsing message reply data
		h.ErrorHandler.HandleWebSocketError(err, ws, "Invalid message format: %s", err)
		return
	}

	if err := h.msgCtrl.SaveMessageReply(messageReplyData); err != nil {
		// Handle error saving message reply to database
		h.ErrorHandler.HandleWebSocketError(err, ws, "Error saving message reply to database: %s", err)
		return
	}

	// Create a final message reply structure and broadcast to other clients
	finalMsgReply := Messages.FinalMessageReply{Type: "message_reply", Message: messageReplyData}
	h.Broadcast.RepliesBroadcast <- finalMsgReply
}
