package chatmessagehandler

import (
	"net/http"

	"github.com/gorilla/websocket"

	Messages "messenger_engine/models/message"
	Broadcast "messenger_engine/controllers/broadcast_controller"
	MessageController "messenger_engine/controllers/message_controller"
	ErrorHandler "messenger_engine/controllers/websocket_controller/handlers/error_handler"
	MessageParser "messenger_engine/controllers/websocket_controller/parsers"
)

// ChatMessageHandler handles WebSocket connections for sending and receiving chat messages.
type ChatMessageHandler struct {
	upgrader      websocket.Upgrader
	msgCtrl       *MessageController.MessageController
	ErrorHandler  *ErrorHandler.ErrorHandler
	MessageParser *MessageParser.Parser
	Broadcast	  *Broadcast.Broadcast
}

// NewChatMessageHandler returns a new instance of ChatMessageHandler.
func NewChatMessageHandler(
	upgrader websocket.Upgrader,
	ctrl *MessageController.MessageController,
) *ChatMessageHandler {
	return &ChatMessageHandler{
		upgrader:      upgrader,
		msgCtrl:       ctrl,
		MessageParser: MessageParser.New(),           // Auto-initialized
		ErrorHandler:  ErrorHandler.NewErrorHandler(), // Auto-initialized
	}
}

// ServeHTTP upgrades the connection and processes incoming chat messages.
func (h *ChatMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := h.upgrader.Upgrade(w, r, nil)
	h.ErrorHandler.HandleWebSocketError(err, nil, "Error upgrading to WebSocket: %s")
	defer ws.Close()

	// Register the client for broadcasts.
	h.Broadcast.Clients[ws] = true

	for {
		var msg map[string]interface{}
		if err := ws.ReadJSON(&msg); err != nil {
			h.ErrorHandler.HandleWebSocketError(err, ws, "Error reading message: %s")
			delete(h.Broadcast.Clients, ws)
			break
		}

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

func (h *ChatMessageHandler) handleInitialMessage(ws *websocket.Conn, msg map[string]interface{}) {
	chatID, err := h.MessageParser.ParseChatID(msg)
	if err != nil {
		h.ErrorHandler.HandleWebSocketError(err, ws, "Invalid chat_id format")
		return
	}

	messages, err := h.msgCtrl.LoadMessages(chatID)
	if err != nil {
		h.ErrorHandler.HandleWebSocketError(err, ws, "Error loading messages: %s")
		return
	}

	if err := ws.WriteJSON(map[string]interface{}{"type": "initial", "messages": messages}); err != nil {
		h.ErrorHandler.HandleWebSocketError(err, ws, "Error sending initial messages: %s")
	}
}

func (h *ChatMessageHandler) handleMessage(ws *websocket.Conn, msg map[string]interface{}) {
	messageData, err := h.MessageParser.ParseMessageData(msg)
	if err != nil {
		h.ErrorHandler.HandleWebSocketError(err, ws, "Invalid message format")
		return
	}

	if err := h.msgCtrl.SaveMessage(messageData); err != nil {
		h.ErrorHandler.HandleWebSocketError(err, ws, "Error saving message to database: %s")
		return
	}

	finalMsg := Messages.FinalMessage{Type: "message", Message: messageData}
	h.Broadcast.Broadcast <- finalMsg
}

func (h *ChatMessageHandler) handleMessageReply(ws *websocket.Conn, msg map[string]interface{}) {
	messageReplyData, err := h.MessageParser.ParseMessageReplyData(msg)
	if err != nil {
		h.ErrorHandler.HandleWebSocketError(err, ws, "Invalid message format")
		return
	}

	if err := h.msgCtrl.SaveMessageReply(messageReplyData); err != nil {
		h.ErrorHandler.HandleWebSocketError(err, ws, "Error saving message reply to database: %s")
		return
	}

	finalMsgReply := Messages.FinalMessageReply{Type: "message_reply", Message: messageReplyData}
	h.Broadcast.RepliesBroadcast <- finalMsgReply
}
