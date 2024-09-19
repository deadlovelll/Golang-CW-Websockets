package main

import (
	"encoding/json"
	"fmt"
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

func (wsh *ChatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsh.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	c, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("error %s when upgrading connection to websocket", err)
		return
	}

	for {

		messageType, message, err := c.ReadMessage()

		if err != nil {
			fmt.Printf("Error reading message %s", err)
			break
		}

		var msg ChatsMessage
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Printf("Error parsing message JSON: %s", err)
			continue
		}

		// Extract the "query" field from the JSON
		UserID := msg.UserID

		// Если преобразование успешно, то это целое число
		jsonData, err := wsh.GMContrl.GetUserChats(UserID)
		if err != nil {
			fmt.Printf("Error fetching users by ID: %s", err)
			return
		}

		// Send the JSON response back to the client
		err = c.WriteMessage(messageType, jsonData)
		if err != nil {
			fmt.Printf("Error sending message: %v", err)
			break
		}

		fmt.Printf("Получено сообщение %s", message)

	}
}

func (cmh *ChatMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cmh.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := cmh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.WriteJSON(map[string]string{"error": fmt.Sprintf("Error upgrading to WebSocket: %s", err)})
		return
	}
	defer ws.Close()

	postmessengercontroller.Clients[ws] = true

	for {
		var msg map[string]interface{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("Error reading message: %s", err)

			delete(postmessengercontroller.Clients, ws)

			ws.WriteJSON(map[string]string{"error": fmt.Sprintf("Error reading message: %s", err)})
			break
		}

		switch msg["type"] {

		case "initial":
			chatIdFloat, ok := msg["chat_id"].(float64)
			if !ok {
				fmt.Println("Invalid chat_id")

				ws.WriteJSON(map[string]string{"error": "Invalid chat_id format"})
				break
			}

			chatId := int(chatIdFloat)

			messages, err := cmh.MMC.LoadMessages(chatId)
			if err != nil {
				fmt.Printf("Error loading messages: %s", err)

				ws.WriteJSON(map[string]string{"error": fmt.Sprintf("Error loading messages: %s", err)})
				break
			}

			err = ws.WriteJSON(map[string]interface{}{
				"type":     "initial",
				"messages": messages,
			})
			if err != nil {
				fmt.Printf("Error sending initial messages: %s", err)
				ws.WriteJSON(map[string]string{"error": fmt.Sprintf("Error sending initial messages: %s", err)})
				break
			}

		case "message":
			messageData, ok := msg["message"].(map[string]interface{})
			if !ok {
				fmt.Println("Invalid message format")

				ws.WriteJSON(map[string]string{"error": "Invalid message format"})
				continue
			}

			messageIdFloat, _ := messageData["MessageId"].(float64)
			authorIdFloat, _ := messageData["AuthorId"].(float64)
			timestampFloat, _ := messageData["Timestamp"].(float64)
			receiverIdFloat, _ := messageData["ReceiverId"].(float64)
			messageText, _ := messageData["Message"].(string)
			chatIdFloat, _ := messageData["ChatId"].(float64)
			isEdited, _ := messageData["IsEdited"].(bool)

			msgData := postmessengercontroller.Mesaage{
				MessageId:  int(messageIdFloat),
				AuthorId:   int(authorIdFloat),
				Timestamp:  time.Unix(int64(timestampFloat), 0),
				ReceiverId: int(receiverIdFloat),
				Message:    messageText,
				ChatId:     int(chatIdFloat),
				IsEdited:   isEdited,
			}

			err = cmh.MMC.SaveMessage(msgData)
			if err != nil {
				fmt.Printf("Error saving message to database: %s", err)

				ws.WriteJSON(map[string]string{"error": fmt.Sprintf("Error saving message to database: %s", err)})
				break
			}

			FinalDict := postmessengercontroller.FinalMessage{Type: "message", Message: msgData}

			if err != nil {
				fmt.Printf("Error sending message back to client: %s", err)
				break
			}

			postmessengercontroller.Broadcast <- FinalDict

		case "message_reply":

			messageReplyData, ok := msg["message"].(map[string]interface{})

			if !ok {
				fmt.Println("Invalid message format")

				ws.WriteJSON(map[string]string{"error": "Invalid message format"})
				continue
			}

			fmt.Println(messageReplyData["ParentMessageId"])

			messageIdFloat, _ := messageReplyData["MessageId"].(float64)
			authorIdFloat, _ := messageReplyData["AuthorId"].(float64)
			timestampFloat, _ := messageReplyData["Timestamp"].(float64)
			receiverIdFloat, _ := messageReplyData["ReceiverId"].(float64)
			messageText, _ := messageReplyData["Message"].(string)
			chatIdFloat, _ := messageReplyData["ChatId"].(float64)
			isEdited, _ := messageReplyData["IsEdited"].(bool)
			parentMessageIdFloat, _ := messageReplyData["ParentMessageId"].(float64)

			msgData := postmessengercontroller.MessageReply{
				MessageId:       int(messageIdFloat),
				AuthorId:        int(authorIdFloat),
				Timestamp:       time.Unix(int64(timestampFloat), 0),
				ReceiverId:      int(receiverIdFloat),
				Message:         messageText,
				ChatId:          int(chatIdFloat),
				IsEdited:        isEdited,
				ParentMessageId: int(parentMessageIdFloat),
			}

			fmt.Println(msgData)

			err = cmh.MMC.SaveMessageReply(msgData)
			if err != nil {
				fmt.Printf("Error saving message to database: %s", err)

				ws.WriteJSON(map[string]string{"error": fmt.Sprintf("Error saving message to database: %s", err)})
				break
			}

			FinalDict := postmessengercontroller.FinalMessageReply{Type: "message_reply", Message: msgData}

			if err != nil {
				fmt.Printf("Error sending message back to client: %s", err)
				break
			}

			postmessengercontroller.RepliesBroadcast <- FinalDict

		}

	}
}

func main() {
	DBPool := &database.DatabasePoolController{}
	DBPool.StartupEvent()

	BContrl := basecontroller.BaseController{Database: DBPool.GetDb()}
	GMContrl := getdatabasecontroller.GetMessengerController{BaseController: &BContrl}
	MMC := postmessengercontroller.MakeMessagesController{BaseController: &BContrl}

	wsHandler := &ChatsHandler{
		upgrader: websocket.Upgrader{},
		GMContrl: &GMContrl,
	}

	chatMessageHandler := &ChatMessageHandler{
		upgrader: websocket.Upgrader{},
		MMC:      &MMC,
	}

	http.Handle("/chats", wsHandler)
	http.Handle("/chat", chatMessageHandler)

	fmt.Println("Запуск сервера на http://localhost:8440")

	go postmessengercontroller.HandleMessages(&MMC)

	server := &http.Server{Addr: "localhost:8440"}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Ошибка запуска сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	DBPool.ShutdownEvent()

	if err := server.Close(); err != nil {
		fmt.Printf("Ошибка при завершении работы сервера: %v", err)
	}

	fmt.Println("Сервер был корректно завершен.")
}
