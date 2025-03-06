package tests

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	websockethandler "user_search/handlers/websocket_handler"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserController is a mock implementation of UserControllerInterface.
type MockUserController struct {
	mock.Mock
}

// GetUsers mocks fetching a user by ID.
func (m *MockUserController) GetUsers(userID int) ([]byte, error) {
	args := m.Mock.Called(userID)
	return args.Get(0).([]byte), args.Error(1)
}

// GetUsersByUsername mocks fetching a user by username.
func (m *MockUserController) GetUsersByUsername(username string) ([]byte, error) {
	args := m.Mock.Called(username)
	return args.Get(0).([]byte), args.Error(1)
}

// setupTestServer creates a test WebSocket server and returns the URL and teardown function.
func setupTestServer(t *testing.T, mockCtrl *MockUserController) (string, func()) {
	handler := websockethandler.NewWebSocketHandler(mockCtrl)
	server := httptest.NewServer(handler)
	url := "ws" + strings.TrimPrefix(server.URL, "http")
	return url, server.Close
}

// TestWebSocketHandler_FetchByID tests retrieving user data by ID over WebSocket.
func TestWebSocketHandler_FetchByID(t *testing.T) {
	mockCtrl := new(MockUserController)
	expectedResponse := []byte(`{"id": 1, "name": "John Doe"}`)
	mockCtrl.Mock.On("GetUsers", 1).Return(expectedResponse, nil)

	wsURL, closeServer := setupTestServer(t, mockCtrl)
	defer closeServer()

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Send a valid user ID query
	request := websockethandler.Message{Query: "1"}
	reqBytes, _ := json.Marshal(request)
	err = conn.WriteMessage(websocket.TextMessage, reqBytes)
	require.NoError(t, err)

	// Receive response
	_, response, err := conn.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, expectedResponse, response)

	mockCtrl.Mock.AssertExpectations(t)
}

// TestWebSocketHandler_FetchByUsername tests retrieving user data by username over WebSocket.
func TestWebSocketHandler_FetchByUsername(t *testing.T) {
	mockCtrl := new(MockUserController)
	expectedResponse := []byte(`{"id": 2, "name": "Jane Doe"}`)
	mockCtrl.Mock.On("GetUsersByUsername", "JaneDoe").Return(expectedResponse, nil)

	wsURL, closeServer := setupTestServer(t, mockCtrl)
	defer closeServer()

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Send a username query
	request := websockethandler.Message{Query: "JaneDoe"}
	reqBytes, _ := json.Marshal(request)
	err = conn.WriteMessage(websocket.TextMessage, reqBytes)
	require.NoError(t, err)

	// Receive response
	_, response, err := conn.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, expectedResponse, response)

	mockCtrl.Mock.AssertExpectations(t)
}

// TestWebSocketHandler_InvalidJSON tests how the WebSocket server handles invalid JSON input.
func TestWebSocketHandler_InvalidJSON(t *testing.T) {
	mockCtrl := new(MockUserController)

	wsURL, closeServer := setupTestServer(t, mockCtrl)
	defer closeServer()

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Send invalid JSON
	err = conn.WriteMessage(websocket.TextMessage, []byte(`{invalid json}`))
	require.NoError(t, err)

	// Expect no response as the server should just log the error
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, _, err = conn.ReadMessage()
	require.Error(t, err) // Should timeout or return an error

	mockCtrl.Mock.AssertExpectations(t)
}

// TestWebSocketHandler_EmptyUsername tests how the WebSocket server handles an empty username query.
func TestWebSocketHandler_EmptyUsername(t *testing.T) {
	mockCtrl := new(MockUserController)

	wsURL, closeServer := setupTestServer(t, mockCtrl)
	defer closeServer()

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Send empty username
	request := websockethandler.Message{Query: ""}
	reqBytes, _ := json.Marshal(request)
	err = conn.WriteMessage(websocket.TextMessage, reqBytes)
	require.NoError(t, err)

	// Expect no response (server logs an error)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second)) // Set a short timeout for waiting
	_, _, err = conn.ReadMessage()
	require.Error(t, err) // Should timeout or return an error

	mockCtrl.Mock.AssertExpectations(t)
}