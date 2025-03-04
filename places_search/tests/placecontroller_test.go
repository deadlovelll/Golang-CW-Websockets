package tests

import (
	"testing"
	"errors"
	"database/sql"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	"places_search/controllers/place_controller"
	"places_search/modules/database/database"
)

// MockDatabase is a mock implementation of DatabaseInterface.
type MockDatabase struct {
	mock.Mock
}

// Query mocks the database query method.
func (m *MockDatabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	argsCalled := m.Mock.Called(query, args...)
	rows, _ := argsCalled.Get(0).(*sql.Rows) // ✅ Ensure type assertion works
	return rows, argsCalled.Error(1)
}

// ✅ Ensure MockDatabase implements DatabaseInterface
var _ database.DatabaseInterface = (*MockDatabase)(nil)

// Mock Rows for simulating the rows returned by database query.
type MockRows struct {
	Mock mock.Mock
}

// Next mocks the Next method of sql.Rows.
func (r *MockRows) Next() bool {
	args := r.Mock.Called()
	return args.Bool(0)
}

// Scan mocks the Scan method of sql.Rows.
func (r *MockRows) Scan(dest ...interface{}) error {
	args := r.Mock.Called(dest)
	return args.Error(0)
}

// Close mocks the Close method of sql.Rows.
func (r *MockRows) Close() error {
	args := r.Mock.Called()
	return args.Error(0)
}

// TestGetPlaceByName tests the GetPlaceByName method of PlaceController.
func TestGetPlaceByName(t *testing.T) {
	// Create mock database and mock rows
	mockDB := new(MockDatabase)
	mockRows := new(MockRows)

	// Set up the mock behavior for database query.
	mockDB.Mock.On("Query", mock.Anything, mock.Anything).Return(mockRows, nil)

	// Set up mock behavior for rows.Next and rows.Scan.
	mockRows.Mock.On("Next").Return(true).Once()
	mockRows.Mock.On("Next").Return(false).Once()
	mockRows.Mock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create PlaceController with mock database.
	pc := placecontroller.NewPlaceController(mockDB)

	// Call the method to test.
	placeName := "test"
	jsonResponse, err := pc.GetPlaceByName(placeName)

	// Assert no error and that JSON response is not nil
	assert.Nil(t, err)
	assert.NotNil(t, jsonResponse)

	// Ensure all expectations were met
	mockDB.Mock.AssertExpectations(t)
	mockRows.Mock.AssertExpectations(t)
}

// TestGetPlaceByNameError tests the GetPlaceByName method when a query error occurs.
func TestGetPlaceByNameError(t *testing.T) {
	// Create mock database
	mockDB := new(MockDatabase)

	// Simulate a database query error
	mockDB.Mock.On("Query", mock.Anything, mock.Anything).Return(nil, errors.New("query error"))

	// Create PlaceController with mock database.
	pc := placecontroller.NewPlaceController(mockDB)

	// Call the method to test
	placeName := "test"
	_, err := pc.GetPlaceByName(placeName)

	// Assert error occurred
	assert.NotNil(t, err)
	assert.Equal(t, "query error", err.Error())

	// Ensure all expectations were met
	mockDB.Mock.AssertExpectations(t)
}

// TestGetPlaceWithHashtag tests the GetPlaceWithHashtag method of PlaceController.
func TestGetPlaceWithHashtag(t *testing.T) {
	// Create mock database and mock rows
	mockDB := new(MockDatabase)
	mockRows := new(MockRows)

	// Set up the mock behavior for database query.
	mockDB.Mock.On("Query", mock.Anything, mock.Anything).Return(mockRows, nil)

	// Set up mock behavior for rows.Next and rows.Scan.
	mockRows.Mock.On("Next").Return(true).Once()
	mockRows.Mock.On("Next").Return(false).Once()
	mockRows.Mock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create PlaceController with mock database.
	pc := placecontroller.NewPlaceController(mockDB)

	// Call the method to test.
	hashtag := "#test"
	jsonResponse, err := pc.GetPlaceWithHashtag(hashtag)

	// Assert no error and that JSON response is not nil
	assert.Nil(t, err)
	assert.NotNil(t, jsonResponse)

	// Ensure all expectations were met
	mockDB.Mock.AssertExpectations(t)
	mockRows.Mock.AssertExpectations(t)
}

// TestGetPlaceWithHashtagError tests the GetPlaceWithHashtag method when a query error occurs.
func TestGetPlaceWithHashtagError(t *testing.T) {
	// Create mock database
	mockDB := new(MockDatabase)

	// Simulate a database query error
	mockDB.Mock.On("Query", mock.Anything, mock.Anything).Return(nil, errors.New("query error"))

	// Create PlaceController with mock database.
	pc := placecontroller.NewPlaceController(mockDB)

	// Call the method to test
	hashtag := "#test"
	_, err := pc.GetPlaceWithHashtag(hashtag)

	// Assert error occurred
	assert.NotNil(t, err)
	assert.Equal(t, "query error", err.Error())

	// Ensure all expectations were met
	mockDB.Mock.AssertExpectations(t)
}
