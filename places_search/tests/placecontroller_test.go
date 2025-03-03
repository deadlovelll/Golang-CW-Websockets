package tests

import (
	"testing"
	"errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	"places_search/controllers/place_controller"
	"places_search/modules/database/database"
)

// MockDatabase is a mock of the Database interface to simulate database interaction.
type MockDatabase struct {
	mock.Mock
}

// Query is a mocked method for executing a database query.
func (m *MockDatabase) Query(query string, args ...interface{}) (database.Rows, error) {
	// This will call the `Called` method from `mock.Mock` to record the method call.
	args = m.Called(query, args)
	return args.Get(0).(database.Rows), args.Error(1)
}

// Mock Rows for simulating the rows returned by database query.
type MockRows struct {
	mock.Mock
}

// Next mocks the Next method of sql.Rows.
func (r *MockRows) Next() bool {
	args := r.Called()
	return args.Bool(0)
}

// Scan mocks the Scan method of sql.Rows.
func (r *MockRows) Scan(dest ...interface{}) error {
	args := r.Called(dest)
	return args.Error(0)
}

// Close mocks the Close method of sql.Rows.
func (r *MockRows) Close() error {
	args := r.Called()
	return args.Error(0)
}

// TestGetPlaceByName tests the GetPlaceByName method of PlaceController.
func TestGetPlaceByName(t *testing.T) {
	// Create mock database and mock rows
	mockDB := new(MockDatabase)
	mockRows := new(MockRows)

	// Set up the mock behavior for database query.
	mockDB.On("Query", mock.Anything, mock.Anything).Return(mockRows, nil)

	// Set up mock behavior for rows.Next and rows.Scan.
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create PlaceController with mock database.
	pc := place_controller.NewPlaceController(mockDB)

	// Call the method to test.
	placeName := "test"
	jsonResponse, err := pc.GetPlaceByName(placeName)

	// Assert no error and that JSON response is not nil
	assert.Nil(t, err)
	assert.NotNil(t, jsonResponse)

	// Ensure all expectations were met
	mockDB.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

// TestGetPlaceByNameError tests the GetPlaceByName method when a query error occurs.
func TestGetPlaceByNameError(t *testing.T) {
	// Create mock database
	mockDB := new(MockDatabase)

	// Simulate a database query error
	mockDB.On("Query", mock.Anything, mock.Anything).Return(nil, errors.New("query error"))

	// Create PlaceController with mock database.
	pc := place_controller.NewPlaceController(mockDB)

	// Call the method to test
	placeName := "test"
	_, err := pc.GetPlaceByName(placeName)

	// Assert error occurred
	assert.NotNil(t, err)
	assert.Equal(t, "query error", err.Error())

	// Ensure all expectations were met
	mockDB.AssertExpectations(t)
}

// TestGetPlaceWithHashtag tests the GetPlaceWithHashtag method of PlaceController.
func TestGetPlaceWithHashtag(t *testing.T) {
	// Create mock database and mock rows
	mockDB := new(MockDatabase)
	mockRows := new(MockRows)

	// Set up the mock behavior for database query.
	mockDB.On("Query", mock.Anything, mock.Anything).Return(mockRows, nil)

	// Set up mock behavior for rows.Next and rows.Scan.
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create PlaceController with mock database.
	pc := place_controller.NewPlaceController(mockDB)

	// Call the method to test.
	hashtag := "#test"
	jsonResponse, err := pc.GetPlaceWithHashtag(hashtag)

	// Assert no error and that JSON response is not nil
	assert.Nil(t, err)
	assert.NotNil(t, jsonResponse)

	// Ensure all expectations were met
	mockDB.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

// TestGetPlaceWithHashtagError tests the GetPlaceWithHashtag method when a query error occurs.
func TestGetPlaceWithHashtagError(t *testing.T) {
	// Create mock database
	mockDB := new(MockDatabase)

	// Simulate a database query error
	mockDB.On("Query", mock.Anything, mock.Anything).Return(nil, errors.New("query error"))

	// Create PlaceController with mock database.
	pc := place_controller.NewPlaceController(mockDB)

	// Call the method to test
	hashtag := "#test"
	_, err := pc.GetPlaceWithHashtag(hashtag)

	// Assert error occurred
	assert.NotNil(t, err)
	assert.Equal(t, "query error", err.Error())

	// Ensure all expectations were met
	mockDB.AssertExpectations(t)
}
