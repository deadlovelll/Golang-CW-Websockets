package main

import (
	"encoding/json"
	"fmt"
	"log"
	"messenger_engine/modules/database"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"

	basecontroller "messenger_engine/modules/base_controller"
	getdatabasecontroller "messenger_engine/get_messenger_controller"
	postmessengercontroller "messenger_engine/post_messenger_controller"
)

type ChatsHandler struct {
	upgrader websocket.Upgrader
	GMContrl *getdatabasecontroller.GetMessengerController
}

type ChatMessageHandler struct {
	upgrader websocket.Upgrader
	MMC      *postmessengercontroller.MakeMessagesController
}

type ChatsMessage struct {
	UserID int `json:"user_id"`
}

// Function to handle WebSocket errors
func handleWebSocketError(err error, ws *websocket.Conn, message string) {
	if err != nil {
		log.Printf(message, err)
		if ws != nil {
			ws.WriteJSON(map[string]string{"error": fmt.Sprintf(message, err)})
		}
	}
}

// Handle incoming messages for chats
func (wsh *ChatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsh.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := wsh.upgrader.Upgrade(w, r, nil)
	handleWebSocketError(err, nil, "error %s when upgrading connection to websocket")

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			handleWebSocketError(err, c, "Error reading message %s")
			break
		}

		var msg ChatsMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			handleWebSocketError(err, c, "Error parsing message JSON: %s")
			continue
		}

		// Fetch and send user chats
		if jsonData, err := wsh.GMContrl.GetUserChats(msg.UserID); err != nil {
			handleWebSocketError(err, c, "Error fetching users by ID: %s")
			return
		} else if err := c.WriteMessage(messageType, jsonData); err != nil {
			handleWebSocketError(err, c, "Error sending message: %v")
			break
		}

		log.Printf("Received message %s", message)
	}
}

// Handle incoming chat messages
func (cmh *ChatMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cmh.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := cmh.upgrader.Upgrade(w, r, nil)
	handleWebSocketError(err, nil, "Error upgrading to WebSocket: %s")

	defer ws.Close()
	postmessengercontroller.Clients[ws] = true

	for {
		var msg map[string]interface{}
		if err := ws.ReadJSON(&msg); err != nil {
			handleWebSocketError(err, ws, "Error reading message: %s")
			delete(postmessengercontroller.Clients, ws)
			break
		}

		switch msg["type"] {
		case "initial":
			cmh.handleInitialMessage(ws, msg)
		case "message":
			cmh.handleMessage(ws, msg)
		case "message_reply":
			cmh.handleMessageReply(ws, msg)
		}
	}
}

func (cmh *ChatMessageHandler) handleInitialMessage(ws *websocket.Conn, msg map[string]interface{}) {
	chatId, err := parseChatID(msg)
	if err != nil {
		handleWebSocketError(err, ws, "Invalid chat_id format")
		return
	}

	messages, err := cmh.MMC.LoadMessages(chatId)
	if err != nil {
		handleWebSocketError(err, ws, "Error loading messages: %s")
		return
	}

	if err := ws.WriteJSON(map[string]interface{}{"type": "initial", "messages": messages}); err != nil {
		handleWebSocketError(err, ws, "Error sending initial messages: %s")
	}
}

func (cmh *ChatMessageHandler) handleMessage(ws *websocket.Conn, msg map[string]interface{}) {
	messageData, err := parseMessageData(msg)
	if err != nil {
		handleWebSocketError(err, ws, "Invalid message format")
		return
	}

	if err := cmh.MMC.SaveMessage(messageData); err != nil {
		handleWebSocketError(err, ws, "Error saving message to database: %s")
		return
	}

	finalMsg := postmessengercontroller.FinalMessage{Type: "message", Message: messageData}
	postmessengercontroller.Broadcast <- finalMsg
}

func (cmh *ChatMessageHandler) handleMessageReply(ws *websocket.Conn, msg map[string]interface{}) {
	messageReplyData, err := parseMessageReplyData(msg)
	if err != nil {
		handleWebSocketError(err, ws, "Invalid message format")
		return
	}

	if err := cmh.MMC.SaveMessageReply(messageReplyData); err != nil {
		handleWebSocketError(err, ws, "Error saving message reply to database: %s")
		return
	}

	finalMsgReply := postmessengercontroller.FinalMessageReply{Type: "message_reply", Message: messageReplyData}
	postmessengercontroller.RepliesBroadcast <- finalMsgReply
}

// Helper function to parse chat ID from the incoming message
func parseChatID(msg map[string]interface{}) (int, error) {
	chatIdFloat, ok := msg["chat_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid chat_id")
	}
	return int(chatIdFloat), nil
}

// Helper function to parse message data
func parseMessageData(msg map[string]interface{}) (postmessengercontroller.Mesaage, error) {
	messageData, ok := msg["message"].(map[string]interface{})
	if !ok {
		return postmessengercontroller.Mesaage{}, fmt.Errorf("invalid message format")
	}

	return postmessengercontroller.Mesaage{
		MessageId:  int(messageData["MessageId"].(float64)),
		AuthorId:   int(messageData["AuthorId"].(float64)),
		Timestamp:  time.Unix(int64(messageData["Timestamp"].(float64)), 0),
		ReceiverId: int(messageData["ReceiverId"].(float64)),
		Message:    messageData["Message"].(string),
		ChatId:     int(messageData["ChatId"].(float64)),
		IsEdited:   messageData["IsEdited"].(bool),
	}, nil
}

// Helper function to parse message reply data
func parseMessageReplyData(msg map[string]interface{}) (postmessengercontroller.MessageReply, error) {
	messageReplyData, ok := msg["message"].(map[string]interface{})
	if !ok {
		return postmessengercontroller.MessageReply{}, fmt.Errorf("invalid message format")
	}

	return postmessengercontroller.MessageReply{
		MessageId:       int(messageReplyData["MessageId"].(float64)),
		AuthorId:        int(messageReplyData["AuthorId"].(float64)),
		Timestamp:       time.Unix(int64(messageReplyData["Timestamp"].(float64)), 0),
		ReceiverId:      int(messageReplyData["ReceiverId"].(float64)),
		Message:         messageReplyData["Message"].(string),
		ChatId:          int(messageReplyData["ChatId"].(float64)),
		IsEdited:        messageReplyData["IsEdited"].(bool),
		ParentMessageId: int(messageReplyData["ParentMessageId"].(float64)),
	}, nil
}

func main() {
	DBPool := &database.DatabasePoolController{}
	DBPool.StartupEvent()

	BContrl := basecontroller.BaseController{Database: DBPool.GetDb()}
	GMContrl := getdatabasecontroller.GetMessengerController{BaseController: &BContrl}
	MMC := postmessengercontroller.MakeMessagesController{BaseController: &BContrl}

	wsHandler := &ChatsHandler{upgrader: websocket.Upgrader{}, GMContrl: &GMContrl}
	chatMessageHandler := &ChatMessageHandler{upgrader: websocket.Upgrader{}, MMC: &MMC}

	http.Handle("/chats", wsHandler)
	http.Handle("/chat", chatMessageHandler)

	log.Println("Starting server on http://localhost:8440")

	go postmessengercontroller.HandleMessages(&MMC)

	server := &http.Server{Addr: "localhost:8440"}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	DBPool.ShutdownEvent()

	if err := server.Close(); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	log.Println("Server successfully shut down.")
}