package messagecontroller

import (
	"fmt"

	BaseController "messenger_engine/controllers/base_controller"
	Messages "messenger_engine/models/message"
)

// MessageController handles message-related logic, including saving and loading messages.
type MessageController struct {
	*BaseController.BaseController // Embeds the base controller for shared functionality
}

// SaveMessage saves a new message into the database.
// This function inserts a message with the provided content, timestamp, author ID, chat ID, and receiver ID.
// The message is saved as not edited and without a parent (indicating it's not a reply).
func (mmc *MessageController) SaveMessage(msg Messages.Message) error {
	db := mmc.Database.GetConnection()
	_, err := db.Exec(`
		INSERT INTO base_chatmessage (content, timestamp, author_id, chat_id, receiver_id, is_edited, parent)
		VALUES ($1, $2, $3, $4, $5, false, null)`,
		msg.Message, msg.Timestamp, msg.AuthorId, msg.ChatId, msg.ReceiverId)
	return err
}

// SaveMessageReply saves a reply to a message in the database.
// This function inserts a reply message with the provided content, timestamp, author ID, chat ID, receiver ID, and parent message ID.
// The reply message is saved as not edited.
func (mmc *MessageController) SaveMessageReply(msg Messages.MessageReply) error {
	db := mmc.Database.GetConnection()
	_, err := db.Exec(`
		INSERT INTO base_chatmessage (content, timestamp, author_id, chat_id, receiver_id, parent_id, is_edited)
		VALUES ($1, $2, $3, $4, $5, $6, false)`,
		msg.Message, msg.Timestamp, msg.AuthorId, msg.ChatId, msg.ReceiverId, msg.ParentMessageId)
	return err
}

// LoadMessages loads messages for a given chat ID from the database.
// This function retrieves all messages for a specific chat and returns them as a slice of Message objects.
func (mmc *MessageController) LoadMessages(chatId int) ([]Messages.Message, error) {
	db := mmc.Database.GetConnection()

	// Query to fetch messages for the given chat ID
	query := `SELECT * FROM base_chatmessage WHERE chat_id = $1`
	rows, err := db.Query(query, chatId)
	if err != nil {
		// Return an error if the query execution fails
		return nil, fmt.Errorf("error loading messages: %w", err)
	}
	defer rows.Close() // Ensure rows are closed after processing

	var messages []Messages.Message
	// Iterate through the query results
	for rows.Next() {
		var msg Messages.Message
		// Scan the row into the Message struct
		if err := rows.Scan(&msg.MessageId, &msg.Message, &msg.IsEdited, &msg.Timestamp, &msg.AuthorId, &msg.ChatId, &msg.ReceiverId, &msg.ParentMessageId); err != nil {
			// Return an error if scanning the row fails
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		// Append the message to the result list
		messages = append(messages, msg)
	}

	// Return the list of messages
	return messages, nil
}
