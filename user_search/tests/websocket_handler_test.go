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

// MockUserController is a mock implementation of a user controller interface.
// It uses testify's mock package to simulate the GetUsers and GetUsersByUsername methods.
type MockUserController struct {
	mock.Mock
}

// GetUsers mocks fetching user data by user ID.
// It returns a byte slice containing the user data and an error if one occurs.
func (m *MockUserController) GetUsers(userID int) ([]byte, error) {
	args := m.Mock.Called(userID)
	return args.Get(0).([]byte), args.Error(1)
}

// GetUsersByUsername mocks fetching user data by username.
// It returns a byte slice containing the user data and an error if one occurs.
func (m *MockUserController) GetUsersByUsername(username string) ([]byte, error) {
	args := m.Mock.Called(username)
	return args.Get(0).([]byte), args.Error(1)
}

// setupTestServer creates a test WebSocket server using the provided mock user controller.
// It returns the WebSocket URL to connect to, as well as a teardown function to close the server.
func setupTestServer(t *testing.T, mockCtrl *MockUserController) (string, func()) {
	handler := websockethandler.NewWebSocketHandler(mockCtrl)
	server := httptest.NewServer(handler)
	url := "ws" + strings.TrimPrefix(server.URL, "http")
	return url, server.Close
}

// TestWebSocketHandler_FetchByID tests retrieving user data by user ID over a WebSocket connection.
// It sets up a test server with a mock user controller, sends a message with a valid user ID,
// and verifies that the response from the server matches the expected user data.
func TestWebSocketHandler_FetchByID(t *testing.T) {
	mockCtrl := new(MockUserController)
	expectedResponse := []byte(`{"id": 1, "name": "John Doe"}`)
	mockCtrl.Mock.On("GetUsers", 1).Return(expectedResponse, nil)

	wsURL, closeServer := setupTestServer(t, mockCtrl)
	defer closeServer()

	// Connect to the WebSocket server.
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Send a valid user ID query to fetch user data.
	request := websockethandler.Message{Query: "1"}
	reqBytes, _ := json.Marshal(request)
	err = conn.WriteMessage(websocket.TextMessage, reqBytes)
	require.NoError(t, err)

	// Receive and verify the response from the WebSocket server.
	_, response, err := conn.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, expectedResponse, response)

	mockCtrl.Mock.AssertExpectations(t)
}

// TestWebSocketHandler_FetchByUsername tests retrieving user data by username over a WebSocket connection.
// It sets up a test server with a mock user controller, sends a message with a valid username,
// and verifies that the response from the server matches the expected user data.
func TestWebSocketHandler_FetchByUsername(t *testing.T) {
	mockCtrl := new(MockUserController)
	expectedResponse := []byte(`{"id": 2, "name": "Jane Doe"}`)
	mockCtrl.Mock.On("GetUsersByUsername", "JaneDoe").Return(expectedResponse, nil)

	wsURL, closeServer := setupTestServer(t, mockCtrl)
	defer closeServer()

	// Connect to the WebSocket server.
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Send a username query to fetch user data.
	request := websockethandler.Message{Query: "JaneDoe"}
	reqBytes, _ := json.Marshal(request)
	err = conn.WriteMessage(websocket.TextMessage, reqBytes)
	require.NoError(t, err)

	// Receive and verify the response from the WebSocket server.
	_, response, err := conn.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, expectedResponse, response)

	mockCtrl.Mock.AssertExpectations(t)
}

// TestWebSocketHandler_InvalidJSON tests the WebSocket server's behavior when it receives invalid JSON input.
// The test sends a malformed JSON message and expects the server to handle the error gracefully,
// typically by not sending any response.
func TestWebSocketHandler_InvalidJSON(t *testing.T) {
	mockCtrl := new(MockUserController)

	wsURL, closeServer := setupTestServer(t, mockCtrl)
	defer closeServer()

	// Connect to the WebSocket server.
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Send invalid JSON.
	err = conn.WriteMessage(websocket.TextMessage, []byte(`{invalid json}`))
	require.NoError(t, err)

	// Set a short deadline and expect a timeout or error due to no valid response.
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, _, err = conn.ReadMessage()
	require.Error(t, err)

	mockCtrl.Mock.AssertExpectations(t)
}

// TestWebSocketHandler_EmptyUsername tests the behavior of the WebSocket server when an empty query is sent.
// It sends an empty username query and expects the server to handle it gracefully, such as by not sending any response.
func TestWebSocketHandler_EmptyUsername(t *testing.T) {
	mockCtrl := new(MockUserController)

	wsURL, closeServer := setupTestServer(t, mockCtrl)
	defer closeServer()

	// Connect to the WebSocket server.
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Send an empty username query.
	request := websockethandler.Message{Query: ""}
	reqBytes, _ := json.Marshal(request)
	err = conn.WriteMessage(websocket.TextMessage, reqBytes)
	require.NoError(t, err)

	// Set a short deadline and expect a timeout or error due to no response.
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, _, err = conn.ReadMessage()
	require.Error(t, err)

	mockCtrl.Mock.AssertExpectations(t)
}