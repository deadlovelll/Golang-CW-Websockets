package tests

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"places_search/controllers/place_controller"
	"places_search/modules/database/database"
)

///////////////////////////////////////////////////////////////////////////////
// MockDatabase
///////////////////////////////////////////////////////////////////////////////

// MockDatabase is a mock implementation of database.DatabaseInterface.
// It uses testify's mock.Mock to record and return values for the Query method.
type MockDatabase struct {
	Mock mock.Mock
}

// Query mocks the Query method of the database.
// It registers the call with testify/mock and returns predetermined values.
// Parameters:
//   - query: The SQL query string.
//   - args:  The arguments for the query.
// Returns:
//   - *sql.Rows: A pointer to sql.Rows (or nil if an error occurs).
//   - error:     An error if one was set in the mock.
func (m *MockDatabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	argsCalled := m.Mock.Called(query, args...)
	rows, _ := argsCalled.Get(0).(*sql.Rows) // Ensure type assertion works.
	return rows, argsCalled.Error(1)
}

// Ensure MockDatabase implements database.DatabaseInterface.
var _ database.DatabaseInterface = (*MockDatabase)(nil)

///////////////////////////////////////////////////////////////////////////////
// MockRows
///////////////////////////////////////////////////////////////////////////////

// MockRows is a mock for simulating *sql.Rows behavior.
// It uses testify's mock.Mock to simulate Next, Scan, and Close methods.
type MockRows struct {
	Mock mock.Mock
}

// Next mocks the Next method of sql.Rows.
// It returns a boolean indicating whether another row is available.
func (r *MockRows) Next() bool {
	args := r.Mock.Called()
	return args.Bool(0)
}

// Scan mocks the Scan method of sql.Rows.
// It populates the destination arguments with data and returns an error if any.
func (r *MockRows) Scan(dest ...interface{}) error {
	args := r.Mock.Called(dest)
	return args.Error(0)
}

// Close mocks the Close method of sql.Rows.
// It returns an error if one is set in the mock.
func (r *MockRows) Close() error {
	args := r.Mock.Called()
	return args.Error(0)
}

///////////////////////////////////////////////////////////////////////////////
// Test Functions for PlaceController
///////////////////////////////////////////////////////////////////////////////

// TestGetPlaceByName tests the GetPlaceByName method of PlaceController.
// It verifies that a valid place name query returns a non-nil JSON response.
func TestGetPlaceByName(t *testing.T) {
	// Create a new instance of MockDatabase and MockRows.
	mockDB := new(MockDatabase)
	mockRows := new(MockRows)

	// Set up the expected behavior for the Query call:
	// It should be called with any query and any arguments,
	// and it should return the mockRows and no error.
	mockDB.Mock.On("Query", mock.Anything, mock.Anything).Return(mockRows, nil)

	// Set up expected behavior for the mockRows:
	// First call to Next returns true, then false.
	mockRows.Mock.On("Next").Return(true).Once()
	mockRows.Mock.On("Next").Return(false).Once()
	// When Scan is called with any arguments, return nil.
	mockRows.Mock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create a new PlaceController using the mock database.
	pc := placecontroller.NewPlaceController(mockDB)

	// Define a test place name.
	placeName := "test"
	// Call GetPlaceByName with the test place name.
	jsonResponse, err := pc.GetPlaceByName(placeName)

	// Assert that no error occurred and that the JSON response is not nil.
	assert.Nil(t, err)
	assert.NotNil(t, jsonResponse)

	// Assert that all expectations on the mocks were met.
	mockDB.Mock.AssertExpectations(t)
	mockRows.Mock.AssertExpectations(t)
}

// TestGetPlaceByNameError tests the behavior of GetPlaceByName when a query error occurs.
// It ensures that the error from the database query is returned.
func TestGetPlaceByNameError(t *testing.T) {
	// Create a new instance of MockDatabase.
	mockDB := new(MockDatabase)

	// Simulate a database query error by setting up the Query method to return an error.
	mockDB.Mock.On("Query", mock.Anything, mock.Anything).Return(nil, errors.New("query error"))

	// Create a new PlaceController using the mock database.
	pc := placecontroller.NewPlaceController(mockDB)

	// Define a test place name.
	placeName := "test"
	// Call GetPlaceByName expecting an error.
	_, err := pc.GetPlaceByName(placeName)

	// Assert that an error occurred and that it matches the expected error message.
	assert.NotNil(t, err)
	assert.Equal(t, "query error", err.Error())

	// Assert that all expectations on the mock were met.
	mockDB.Mock.AssertExpectations(t)
}

// TestGetPlaceWithHashtag tests the GetPlaceWithHashtag method of PlaceController.
// It verifies that a valid hashtag query returns a non-nil JSON response.
func TestGetPlaceWithHashtag(t *testing.T) {
	// Create a new instance of MockDatabase and MockRows.
	mockDB := new(MockDatabase)
	mockRows := new(MockRows)

	// Set up the expected behavior for the Query call.
	mockDB.Mock.On("Query", mock.Anything, mock.Anything).Return(mockRows, nil)

	// Set up expected behavior for mockRows:
	// First call to Next returns true, then false.
	mockRows.Mock.On("Next").Return(true).Once()
	mockRows.Mock.On("Next").Return(false).Once()
	// When Scan is called, return nil.
	mockRows.Mock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create a new PlaceController using the mock database.
	pc := placecontroller.NewPlaceController(mockDB)

	// Define a test hashtag.
	hashtag := "#test"
	// Call GetPlaceWithHashtag with the test hashtag.
	jsonResponse, err := pc.GetPlaceWithHashtag(hashtag)

	// Assert that no error occurred and that the JSON response is not nil.
	assert.Nil(t, err)
	assert.NotNil(t, jsonResponse)

	// Assert that all expectations on the mocks were met.
	mockDB.Mock.AssertExpectations(t)
	mockRows.Mock.AssertExpectations(t)
}

// TestGetPlaceWithHashtagError tests the behavior of GetPlaceWithHashtag when a query error occurs.
// It ensures that the error from the database query is properly propagated.
func TestGetPlaceWithHashtagError(t *testing.T) {
	// Create a new instance of MockDatabase.
	mockDB := new(MockDatabase)

	// Simulate a database query error by setting up the Query method to return an error.
	mockDB.Mock.On("Query", mock.Anything, mock.Anything).Return(nil, errors.New("query error"))

	// Create a new PlaceController using the mock database.
	pc := placecontroller.NewPlaceController(mockDB)

	// Define a test hashtag.
	hashtag := "#test"
	// Call GetPlaceWithHashtag expecting an error.
	_, err := pc.GetPlaceWithHashtag(hashtag)

	// Assert that an error occurred and that it matches the expected error message.
	assert.NotNil(t, err)
	assert.Equal(t, "query error", err.Error())

	// Assert that all expectations on the mock were met.
	mockDB.Mock.AssertExpectations(t)
}
