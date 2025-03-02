package tests

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	// Adjust these imports according to your project structure.
	BaseController "messenger_engine/controllers/base_controller"
	messagecontroller "messenger_engine/controllers/message_controller"
	Messages "messenger_engine/models/message"
	"messenger_engine/modules/database/database"
)

// DummyDatabase is a simple implementation of a database wrapper that implements GetConnection().
type DummyDatabase struct {
	db *sql.DB
}

// GetConnection returns the underlying sql.DB connection.
func (d *DummyDatabase) GetConnection() *sql.DB {
	return d.db
}

func newTestMessageController(db *sql.DB) *messagecontroller.MessageController {
	dummyDB := &database.Database{}
	baseCtrl := &BaseController.BaseController{
		Database: dummyDB,
	}
	return &messagecontroller.MessageController{
		BaseController: baseCtrl,
	}
}

func TestSaveMessage(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	// Create the MessageController with our dummy database.
	mmc := newTestMessageController(db)

	// Create a sample message.
	testMessage := Messages.Message{
		Message:    "Hello, world!",
		Timestamp:  time.Now(),
		AuthorId:   1,
		ChatId:     10,
		ReceiverId: 2,
	}

	// Set up the expectation for the Exec call.
	query := `
		INSERT INTO base_chatmessage \(content, timestamp, author_id, chat_id, receiver_id, is_edited, parent\)
		VALUES \(\$1, \$2, \$3, \$4, \$5, false, null\)`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(testMessage.Message, testMessage.Timestamp, testMessage.AuthorId, testMessage.ChatId, testMessage.ReceiverId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	// Call SaveMessage and check the error.
	if err := mmc.SaveMessage(testMessage); err != nil {
		t.Errorf("SaveMessage() returned an unexpected error: %v", err)
	}

	// Ensure all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestSaveMessageReply(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	// Create the MessageController.
	mmc := newTestMessageController(db)

	// Create a sample message reply.
	testReply := Messages.MessageReply{
		Message:         "This is a reply",
		Timestamp:       time.Now(),
		AuthorId:        3,
		ChatId:          10,
		ReceiverId:      1,
		ParentMessageId: 5,
	}

	// Set up the expectation for the Exec call.
	query := `
		INSERT INTO base_chatmessage \(content, timestamp, author_id, chat_id, receiver_id, parent_id, is_edited\)
		VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, false\)`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(testReply.Message, testReply.Timestamp, testReply.AuthorId, testReply.ChatId, testReply.ReceiverId, testReply.ParentMessageId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call SaveMessageReply and check the error.
	if err := mmc.SaveMessageReply(testReply); err != nil {
		t.Errorf("SaveMessageReply() returned an unexpected error: %v", err)
	}

	// Ensure all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestLoadMessages(t *testing.T) {
	// Create a new sqlmock database connection.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	// Create the MessageController.
	mmc := newTestMessageController(db)

	chatId := 10

	// Define the columns as expected by LoadMessages.
	columns := []string{
		"message_id",
		"content",
		"is_edited",
		"timestamp",
		"author_id",
		"chat_id",
		"receiver_id",
		"parent_message_id",
	}

	// Create sample rows.
	timestamp := time.Now()
	rows := sqlmock.NewRows(columns).
		AddRow(1, "Hello", false, timestamp, 1, chatId, 2, nil).
		AddRow(2, "Hi there", false, timestamp, 2, chatId, 1, nil)

	// Expect the query to be executed.
	query := `SELECT \* FROM base_chatmessage WHERE chat_id = \$1`
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(chatId).
		WillReturnRows(rows)

	// Call LoadMessages.
	msgs, err := mmc.LoadMessages(chatId)
	if err != nil {
		t.Errorf("LoadMessages() returned an unexpected error: %v", err)
	}

	// Check that the correct number of messages were returned.
	if len(msgs) != 2 {
		t.Errorf("expected 2 messages, got %d", len(msgs))
	}

	// Verify contents of the first message.
	if msgs[0].Message != "Hello" {
		t.Errorf("expected first message to be 'Hello', got '%s'", msgs[0].Message)
	}

	// Ensure all expectations were met.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
