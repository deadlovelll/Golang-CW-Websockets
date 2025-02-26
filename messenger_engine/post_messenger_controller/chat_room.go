package postmessengercontroller

import (
	"database/sql"
	"fmt"
	"time"

	basecontroller "messenger_engine/modules/base_controller"
	Messages "messenger_engine/models/message"

	"github.com/gorilla/websocket"
)

// ----------------------------------------------------------------------
// Repository Layer
// ----------------------------------------------------------------------

// MessageRepository defines operations for message persistence.
type MessageRepository interface {
	SaveMessage(msg Messages.Message) error
	SaveMessageReply(msg Messages.MessageReply) error
	LoadMessages(chatID int) ([]Messages.Message, error)
}

// PostgresMessageRepository is a concrete implementation of MessageRepository using PostgreSQL.
type PostgresMessageRepository struct {
	DB *sql.DB
}

// SaveMessage persists a new message in the database.
func (r *PostgresMessageRepository) SaveMessage(msg Messages.Message) error {
	query := `
		INSERT INTO base_chatmessage 
		(content, timestamp, author_id, chat_id, receiver_id, is_edited, parent) 
		VALUES ($1, $2, $3, $4, $5, false, null)
	`
	_, err := r.DB.Exec(query, msg.Message, msg.Timestamp, msg.AuthorId, msg.ChatId, msg.ReceiverId)
	return err
}

// SaveMessageReply persists a reply message in the database.
func (r *PostgresMessageRepository) SaveMessageReply(msg Messages.MessageReply) error {
	query := `
		INSERT INTO base_chatmessage 
		(content, timestamp, author_id, chat_id, receiver_id, parent_id, is_edited) 
		VALUES ($1, $2, $3, $4, $5, $6, false)
	`
	_, err := r.DB.Exec(query, msg.Message, msg.Timestamp, msg.AuthorId, msg.ChatId, msg.ReceiverId, msg.ParentMessageId)
	return err
}

// LoadMessages retrieves messages for a given chat from the database.
func (r *PostgresMessageRepository) LoadMessages(chatID int) ([]Messages.Message, error) {
	query := `
		SELECT 
			message_id, 
			content, 
			is_edited, 
			timestamp, 
			author_id, 
			chat_id, 
			receiver_id, 
			parent
		FROM 
			base_chatmessage
		WHERE 
			chat_id = $1
	`
	rows, err := r.DB.Query(query, chatID)
	if err != nil {
		return nil, fmt.Errorf("querying messages: %w", err)
	}
	defer rows.Close()

	var messages []Messages.Message
	for rows.Next() {
		var msg Messages.Message
		if err := rows.Scan(
			&msg.MessageId,
			&msg.Message,
			&msg.IsEdited,
			&msg.Timestamp,
			&msg.AuthorId,
			&msg.ChatId,
			&msg.ReceiverId,
			&msg.ParentMessageId,
		); err != nil {
			return nil, fmt.Errorf("scanning message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// ----------------------------------------------------------------------
// Controller Layer
// ----------------------------------------------------------------------

// MessageController handles message operations using a repository.
type MessageController struct {
	BaseController *basecontroller.BaseController
	Repo           MessageRepository
}

// NewMessageController creates a new MessageController instance.
func NewMessageController(bc *basecontroller.BaseController, repo MessageRepository) *MessageController {
	return &MessageController{
		BaseController: bc,
		Repo:           repo,
	}
}

// SaveMessage saves a new message using the repository.
func (mc *MessageController) SaveMessage(msg Messages.Message) error {
	return mc.Repo.SaveMessage(msg)
}

// SaveMessageReply saves a reply message using the repository.
func (mc *MessageController) SaveMessageReply(msg Messages.MessageReply) error {
	return mc.Repo.SaveMessageReply(msg)
}

// LoadMessages retrieves all messages for a specific chat.
func (mc *MessageController) LoadMessages(chatID int) ([]Messages.Message, error) {
	return mc.Repo.LoadMessages(chatID)
}

// ----------------------------------------------------------------------
// WebSocket Broadcaster
// ----------------------------------------------------------------------

// MessageBroadcaster manages WebSocket clients and message broadcasting.
type MessageBroadcaster struct {
	Clients          map[*websocket.Conn]bool
	Broadcast        chan Messages.FinalMessage
	RepliesBroadcast chan Messages.FinalMessageReply
}

// NewMessageBroadcaster initializes and returns a new MessageBroadcaster.
func NewMessageBroadcaster() *MessageBroadcaster {
	return &MessageBroadcaster{
		Clients:          make(map[*websocket.Conn]bool),
		Broadcast:        make(chan Messages.FinalMessage),
		RepliesBroadcast: make(chan Messages.FinalMessageReply),
	}
}

// HandleMessages continuously reads messages from the broadcast channel and sends them to connected clients.
func (mb *MessageBroadcaster) HandleMessages() {
	for msg := range mb.Broadcast {
		for client := range mb.Clients {
			if err := client.WriteJSON(msg); err != nil {
				fmt.Printf("Error sending message: %v\n", err)
				client.Close()
				delete(mb.Clients, client)
			}
		}
	}
}

// ----------------------------------------------------------------------
// Example Message struct definitions for context (from messenger_engine/models/message)
// ----------------------------------------------------------------------

// These types are defined elsewhere, but shown here for reference.
/*
type Message struct {
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
	Message Message `json:"message"`
	Type    string  `json:"type"`
}

type FinalMessageReply struct {
	Message MessageReply `json:"message"`
	Type    string       `json:"type"`
}
*/
