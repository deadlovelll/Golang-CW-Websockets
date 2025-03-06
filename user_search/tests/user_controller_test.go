package tests

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"user_search/controllers/base_controller"
	"user_search/controllers/user_controller"
	"user_search/modules/database/database"

)

// MockDatabase is a mock implementation of the Database interface.
type MockDatabase struct {
	mock.Mock
}

// MockDatabasePoolController is a mock of the DatabasePoolController.
type MockDatabasePoolController struct {
	mock.Mock
}

func (m *MockDatabasePoolController) GetDb() *database.Database {
	args := m.Called()
	return args.Get(0).(*database.Database)
}

// Query mocks the Query method of the database instance.
func (m *MockDatabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	// Call the mock method with the query and args
	callArgs := m.Called(query, args)
	
	// Return the mocked results
	return callArgs.Get(0).(*sql.Rows), callArgs.Error(1)
}

type MockRows struct {
	mock.Mock
}

func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRows) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

type MockUrlFetcher struct {
	mock.Mock
}

func (m *MockUrlFetcher) Fetch(url string) {
	m.Called(url)
}

func TestGetUsers(t *testing.T) {
	// Setup mock database connection
	mockDb := new(MockDatabase)
	mockBaseCtrl := new(MockBaseController)

	// Mock DatabasePoolController's GetDb method to return the mock Database
	mockDbPoolCtrl := new(MockDatabasePoolController)
	mockDbPoolCtrl.On("GetDb").Return(&database.Database{Db: mockDb}) // Return the mocked database

	// Create the user controller instance
	userCtrl := usercontroller.NewUserController(mockBaseCtrl)

	// Mock database query result
	mockRows := new(MockRows)
	mockDb.On("Query", mock.Anything, mock.Anything).Return(mockRows, nil)

	// Simulate that rows.Next() is true and then false
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	// Mock scan behavior
	mockRows.On("Scan", mock.Anything).Return(nil)

	// Mock avatar URL fetcher
	mockUrlFetcher := new(MockUrlFetcher)
	mockUrlFetcher.On("Fetch", mock.Anything).Return()

	// Call the GetUsers method
	result, err := userCtrl.GetUsers(1)

	// Assert that there is no error
	assert.NoError(t, err)

	// Assert that result is a valid JSON array
	var response []map[string]interface{}
	err = json.Unmarshal(result, &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1) // One user in response

	// Verify mock interactions
	mockBaseCtrl.AssertExpectations(t)
	mockDb.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockUrlFetcher.AssertExpectations(t)
}

func TestGetUsersByUsername(t *testing.T) {
	// Setup mock database connection
	mockDb := new(MockDatabase)
	mockBaseCtrl := new(MockBaseController)
	mockBaseCtrl.On("Database").Return(&basecontroller.Database{Connection: mockDb})

	// Create the user controller instance
	userCtrl := usercontroller.NewUserController(mockBaseCtrl)

	// Mock database query result
	mockRows := new(MockRows)
	mockDb.On("Query", mock.Anything, mock.Anything).Return(mockRows, nil)

	// Simulate that rows.Next() is true and there is data
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()

	// Mock scan behavior
	mockRows.On("Scan", mock.Anything).Return(nil)

	// Mock avatar URL fetcher
	mockUrlFetcher := new(MockUrlFetcher)
	mockUrlFetcher.On("Fetch", mock.Anything).Return()

	// Call the GetUsersByUsername method
	result, err := userCtrl.GetUsersByUsername("testuser")

	// Assert that there is no error
	assert.NoError(t, err)

	// Assert that result is a valid JSON array
	var response []map[string]interface{}
	err = json.Unmarshal(result, &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1) // One user in response

	// Verify mock interactions
	mockBaseCtrl.AssertExpectations(t)
	mockDb.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockUrlFetcher.AssertExpectations(t)
}

func TestGetUsersQueryError(t *testing.T) {
	// Setup mock database connection
	mockDb := new(MockDatabase)
	mockBaseCtrl := new(MockBaseController)
	mockBaseCtrl.On("Database").Return(&basecontroller.Database{Connection: mockDb})

	// Create the user controller instance
	userCtrl := usercontroller.NewUserController(mockBaseCtrl)

	// Mock a query error
	mockDb.On("Query", mock.Anything, mock.Anything).Return(nil, errors.New("query error"))

	// Call the GetUsers method
	result, err := userCtrl.GetUsers(1)

	// Assert that an error occurred
	assert.Error(t, err)
	assert.Nil(t, result)

	// Verify mock interactions
	mockBaseCtrl.AssertExpectations(t)
	mockDb.AssertExpectations(t)
}

func TestScanUsersError(t *testing.T) {
	// Setup mock database connection
	mockDb := new(MockDatabase)
	mockBaseCtrl := new(MockBaseController)
	mockBaseCtrl.On("Database").Return(&basecontroller.Database{Connection: mockDb})

	// Create the user controller instance
	userCtrl := usercontroller.NewUserController(mockBaseCtrl)

	// Mock database query result with an error in scanning rows
	mockRows := new(MockRows)
	mockDb.On("Query", mock.Anything, mock.Anything).Return(mockRows, nil)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Scan", mock.Anything).Return(errors.New("scan error"))

	// Call the GetUsers method
	result, err := userCtrl.GetUsers(1)

	// Assert that an error occurred
	assert.Error(t, err)
	assert.Nil(t, result)

	// Verify mock interactions
	mockBaseCtrl.AssertExpectations(t)
	mockDb.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}
