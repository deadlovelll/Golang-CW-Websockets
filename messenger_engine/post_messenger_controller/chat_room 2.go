package postmessengercontroller

import (
	"database/sql"
	"fmt"
	"time"

	basecontroller "messenger_engine/modules/base_controller"

	"github.com/gorilla/websocket"
)

var Clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan FinalMessage)
var RepliesBroadcast = make(chan FinalMessageReply)

type MakeMessagesController struct {
	*basecontroller.BaseController
}

type Mesaage struct {
	MessageId       int           `json:"message_id"`
	AuthorId        int           `json:"author_id"`
	Timestamp       time.Time     `json:"timestamp"`
	ReceiverId      int           `json:"receiver_id"`
	Message         string        `json:"message"`
	ChatId          int           `json:"chat_id"`
	IsEdited        bool          `json:"is_edited"`
	ParentMessageId sql.NullInt64 `json:"parent_message_id"`
}

type MessageReply struct {
	MessageId       int       `json:"message_id"`
	AuthorId        int       `json:"author_id"`
	Timestamp       time.Time `json:"timestamp"`
	ReceiverId      int       `json:"receiver_id"`
	Message         string    `json:"message"`
	ChatId          int       `json:"chat_id"`
	IsEdited        bool      `json:"is_edited"`
	ParentMessageId int       `json:"parent_message_id"`
}

type FinalMessage struct {
	Message Mesaage `json:"message"`
	Type    string  `json:"type"`
}

type FinalMessageReply struct {
	Message MessageReply `json:"message"`
	Type    string       `json:"type"`
}

func (mmc *MakeMessagesController) SaveMessage(msg Mesaage) error {
	db := mmc.Database.GetConnection()
	_, err := db.Exec("INSERT INTO base_chatmessage (content, timestamp, author_id, chat_id, receiver_id, is_edited, parent) VALUES ($1, $2, $3, $4, $5, false, null)",
		msg.Message, msg.Timestamp, msg.AuthorId, msg.ChatId, msg.ReceiverId)
	return err
}

func (mmc *MakeMessagesController) SaveMessageReply(msg MessageReply) error {
	db := mmc.Database.GetConnection()
	_, err := db.Exec("INSERT INTO base_chatmessage (content, timestamp, author_id, chat_id, receiver_id, parent_id, is_edited) VALUES ($1, $2, $3, $4, $5, $6, false)",
		msg.Message, msg.Timestamp, msg.AuthorId, msg.ChatId, msg.ReceiverId, msg.ParentMessageId)
	return err
}

// Обрабатываем сообщения и рассылаем их клиентам
func HandleMessages(mmc *MakeMessagesController) {
	for {
		msg := <-Broadcast
		for client := range Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Printf("Ошибка отправки сообщения: %v", err)
				client.Close()
				delete(Clients, client)
			}
		}
	}
}

func (mmc *MakeMessagesController) LoadMessages(chatId int) ([]Mesaage, error) {

	db := mmc.Database.GetConnection()

	query := `
	SELECT * 

	FROM base_chatmessage

	WHERE chat_id = $1

	`

	rows, err := db.Query(query, chatId)

	if err != nil {
		fmt.Printf("Error occured: %s", err)
	}

	defer rows.Close()

	var messages []Mesaage

	for rows.Next() {

		var msg Mesaage
		err := rows.Scan(&msg.MessageId, &msg.Message, &msg.IsEdited, &msg.Timestamp, &msg.AuthorId, &msg.ChatId, &msg.ReceiverId, &msg.ParentMessageId)

		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	return messages, nil
}
