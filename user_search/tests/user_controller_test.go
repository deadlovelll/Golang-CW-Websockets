package tests

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"user_search/controllers/base_controller"
	"user_search/controllers/user_controller"
	"user_search/modules/database/database"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *usercontroller.UserController) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	// Correctly create a Database instance with a mock connection
	mockDatabase := &database.Database{
		Db: db, // Assign the mock SQL database
	}

	// Pass the mock database to BaseController
	baseCtrl := &basecontroller.BaseController{
		Database: mockDatabase, // Correct reference
	}

	userCtrl := usercontroller.NewUserController(baseCtrl)
	return db, mock, userCtrl
}

// TestGetUsers_Success tests GetUsers method with valid database response.
func TestGetUsers_Success(t *testing.T) {
	db, mock, userCtrl := setupMockDB(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"userid", "username", "verifiedAccount", "fullName", "hasActiveStories", "userAvatar", "commonFriends"}).
		AddRow(1, "john_doe", true, "John Doe", true, "", "[2,3]")

	mock.ExpectQuery("SELECT .* FROM users").WithArgs(1).WillReturnRows(rows)

	res, err := userCtrl.GetUsers(1)
	require.NoError(t, err)
	require.NotNil(t, res)

	var users []map[string]interface{}
	require.NoError(t, json.Unmarshal(res, &users))
	require.Len(t, users, 1)
	require.Equal(t, "john_doe", users[0]["username"])
}

// TestGetUsersByUsername_Success tests GetUsersByUsername method with valid database response.
func TestGetUsersByUsername_Success(t *testing.T) {
	db, mock, userCtrl := setupMockDB(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"userid", "username", "verifiedAccount", "fullName", "hasActiveStories", "userAvatar", "commonFriends"}).
		AddRow(2, "jane_doe", true, "Jane Doe", false, "", "[1,3]")

	mock.ExpectQuery("SELECT .* FROM users").WithArgs("jane_doe").WillReturnRows(rows)

	res, err := userCtrl.GetUsersByUsername("jane_doe")
	require.NoError(t, err)
	require.NotNil(t, res)

	var users []map[string]interface{}
	require.NoError(t, json.Unmarshal(res, &users))
	require.Len(t, users, 1)
	require.Equal(t, "jane_doe", users[0]["username"])
}

// TestGetUsers_DBError tests GetUsers method when the database returns an error.
func TestGetUsers_DBError(t *testing.T) {
	db, mock, userCtrl := setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT .* FROM users").WithArgs(1).WillReturnError(errors.New("DB error"))

	res, err := userCtrl.GetUsers(1)
	require.Error(t, err)
	require.Nil(t, res)
}

// TestGetUsersByUsername_DBError tests GetUsersByUsername method when the database returns an error.
func TestGetUsersByUsername_DBError(t *testing.T) {
	db, mock, userCtrl := setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT .* FROM users").WithArgs("jane_doe").WillReturnError(errors.New("DB error"))

	res, err := userCtrl.GetUsersByUsername("jane_doe")
	require.Error(t, err)
	require.Nil(t, res)
}
