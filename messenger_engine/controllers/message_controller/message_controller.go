package messagecontroller

import (
	"fmt"

	BaseController "messenger_engine/controllers/base_controller"
	Messages "messenger_engine/models/message"
)

// MakeMessagesController handles message-related logic.
type MessageController struct {
	*BaseController.BaseController
}

// SaveMessage saves a new message into the database.
func (mmc *MessageController) SaveMessage(msg Messages.Mesaage) error {
	db := mmc.Database.GetConnection()
	_, err := db.Exec(`
		INSERT INTO base_chatmessage (content, timestamp, author_id, chat_id, receiver_id, is_edited, parent)
		VALUES ($1, $2, $3, $4, $5, false, null)`,
		msg.Message, msg.Timestamp, msg.AuthorId, msg.ChatId, msg.ReceiverId)
	return err
}

// SaveMessageReply saves a reply to a message in the database.
func (mmc *MessageController) SaveMessageReply(msg Messages.MessageReply) error {
	db := mmc.Database.GetConnection()
	_, err := db.Exec(`
		INSERT INTO base_chatmessage (content, timestamp, author_id, chat_id, receiver_id, parent_id, is_edited)
		VALUES ($1, $2, $3, $4, $5, $6, false)`,
		msg.Message, msg.Timestamp, msg.AuthorId, msg.ChatId, msg.ReceiverId, msg.ParentMessageId)
	return err
}

// LoadMessages loads messages for a given chatId from the database.
func (mmc *MessageController) LoadMessages(chatId int) ([]Messages.Mesaage, error) {
	db := mmc.Database.GetConnection()

	query := `SELECT * FROM base_chatmessage WHERE chat_id = $1`
	rows, err := db.Query(query, chatId)
	if err != nil {
		return nil, fmt.Errorf("error loading messages: %w", err)
	}
	defer rows.Close()

	var messages []Messages.Mesaage
	for rows.Next() {
		var msg Messages.Mesaage
		if err := rows.Scan(&msg.MessageId, &msg.Message, &msg.IsEdited, &msg.Timestamp, &msg.AuthorId, &msg.ChatId, &msg.ReceiverId, &msg.ParentMessageId); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
