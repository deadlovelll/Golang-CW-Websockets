package postmessengercontroller

import (
	"database/sql"
	"fmt"
	"time"

	basecontroller "messenger_engine/modules/base_controller"
	"github.com/gorilla/websocket"
)

var (
	Clients          = make(map[*websocket.Conn]bool)
	Broadcast        = make(chan FinalMessage)
	RepliesBroadcast = make(chan FinalMessageReply)
)

// MakeMessagesController handles message-related logic.
type MakeMessagesController struct {
	*basecontroller.BaseController
}

// Message represents a message structure.
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

// MessageReply represents a reply to a message.
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

// FinalMessage represents the final message format to be sent to the client.
type FinalMessage struct {
	Message Mesaage `json:"message"`
	Type    string  `json:"type"`
}

// FinalMessageReply represents the final reply format to be sent to the client.
type FinalMessageReply struct {
	Message MessageReply `json:"message"`
	Type    string       `json:"type"`
}

// SaveMessage saves a new message into the database.
func (mmc *MakeMessagesController) SaveMessage(msg Mesaage) error {
	db := mmc.Database.GetConnection()
	_, err := db.Exec(`
		INSERT INTO base_chatmessage (content, timestamp, author_id, chat_id, receiver_id, is_edited, parent)
		VALUES ($1, $2, $3, $4, $5, false, null)`,
		msg.Message, msg.Timestamp, msg.AuthorId, msg.ChatId, msg.ReceiverId)
	return err
}

// SaveMessageReply saves a reply to a message in the database.
func (mmc *MakeMessagesController) SaveMessageReply(msg MessageReply) error {
	db := mmc.Database.GetConnection()
	_, err := db.Exec(`
		INSERT INTO base_chatmessage (content, timestamp, author_id, chat_id, receiver_id, parent_id, is_edited)
		VALUES ($1, $2, $3, $4, $5, $6, false)`,
		msg.Message, msg.Timestamp, msg.AuthorId, msg.ChatId, msg.ReceiverId, msg.ParentMessageId)
	return err
}

// HandleMessages handles broadcasting messages to all connected clients.
func HandleMessages(mmc *MakeMessagesController) {
	for {
		msg := <-Broadcast
		for client := range Clients {
			if err := client.WriteJSON(msg); err != nil {
				handleClientError(client, err)
			}
		}
	}
}

// LoadMessages loads messages for a given chatId from the database.
func (mmc *MakeMessagesController) LoadMessages(chatId int) ([]Mesaage, error) {
	db := mmc.Database.GetConnection()

	query := `SELECT * FROM base_chatmessage WHERE chat_id = $1`
	rows, err := db.Query(query, chatId)
	if err != nil {
		return nil, fmt.Errorf("Error loading messages: %w", err)
	}
	defer rows.Close()

	var messages []Mesaage
	for rows.Next() {
		var msg Mesaage
		if err := rows.Scan(&msg.MessageId, &msg.Message, &msg.IsEdited, &msg.Timestamp, &msg.AuthorId, &msg.ChatId, &msg.ReceiverId, &msg.ParentMessageId); err != nil {
			return nil, fmt.Errorf("Error scanning row: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// handleClientError is used to handle client errors by closing the connection and cleaning up.
func handleClientError(client *websocket.Conn, err error) {
	fmt.Printf("Error sending message: %v\n", err)
	client.Close()
	delete(Clients, client)
}