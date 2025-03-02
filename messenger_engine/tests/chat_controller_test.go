package tests

import (
	"encoding/json"
	"errors"
	"messenger_engine/controllers/chat_controller"
	"messenger_engine/controllers/url_controller"
	"messenger_engine/models/presigned_url"
	"messenger_engine/controllers/base_controller"
	"messenger_engine/modules/database/database"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// MockUrlFetcher is a struct for mocking the presigned URL fetch request.
type MockUrlFetcher struct{}

// Fetch simulates fetching a presigned URL and always returns a fixed URL.
func (m *MockUrlFetcher) Fetch(url string) (*presignedurl.PresignedUrlResponse, error) {
	return &presignedurl.PresignedUrlResponse{PresignedURL: "http://example.com/avatar.jpg"}, nil
}

// TestGetUserChats verifies the successful execution of the GetUserChats method.
//
// Steps:
//  1. Creates a mock database.
//  2. Adds test data to the messages table.
//  3. Calls `GetUserChats` and checks if the result is correct.
//
// Expected result:
//  - The function should return a valid JSON containing the chat data.
func TestGetUserChats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Mock query response
	mockRows := sqlmock.NewRows([]string{"content", "timestamp", "chat_id", "receiver_id", "receiver_username"}).
		AddRow("Hello!", "2025-03-02T12:00:00Z", 1, 2, "Alice")
	mock.ExpectQuery("SELECT ").WillReturnRows(mockRows)

	// Initialize ChatController with a mock database and URL fetcher
	baseController := &chatcontroller.ChatController{
		BaseController: &basecontroller.BaseController{
			Database: &database.Database{}, 
		},
		UrlFetcher: urlcontroller.HttpPresignedUrlFetcher{},
	}

	baseController.Database.GetConnection()

	// Execute the function and check for errors
	result, err := baseController.GetUserChats(1)
	assert.NoError(t, err)

	// Unmarshal JSON response
	var chats []map[string]interface{}
	err = json.Unmarshal(result, &chats)
	assert.NoError(t, err)

	// Verify the content of the response
	assert.Len(t, chats, 1)
	assert.Equal(t, "Hello!", chats[0]["messge_content"])
	assert.Equal(t, "http://example.com/avatar.jpg", chats[0]["user_avatar_url"])
}

// TestGetUserChats_QueryError verifies that GetUserChats handles query errors properly.
//
// Steps:
//  1. Creates a mock database that returns an error on query execution.
//  2. Calls `GetUserChats` and ensures it returns an error.
//
// Expected result:
//  - The function should return an error due to a failed query execution.
func TestGetUserChats_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Simulate a database query failure
	mock.ExpectQuery("SELECT ").WillReturnError(errors.New("database error"))

	// Initialize ChatController with a mock database and URL fetcher
	baseController := &chatcontroller.ChatController{
		BaseController: &basecontroller.BaseController{
			Database: &database.Database{}, 
		},
		UrlFetcher: urlcontroller.HttpPresignedUrlFetcher{},
	}

	baseController.Database.GetConnection()

	// Execute the function and check for errors
	_, err = baseController.GetUserChats(1)
	assert.Error(t, err)
}
