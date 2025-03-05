package tests

import (
	"database/sql"
	"encoding/json"
	"testing"

	"hashtags_search/controllers/hashtag_controller"
)

// FakeDatabase is a stub implementation used to satisfy the Database dependency
// in BaseController during testing.
type FakeDatabase struct{}

func (fd *FakeDatabase) GetConnection() *sql.DB {
	return &sql.DB{} // Returning a mock *sql.DB object
}

// Mock the Query method (if needed)
func (fd *FakeDatabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil // Mock returning empty result, adjust as needed
}

// TestHashtagController_GetHashtags_Success verifies that GetHashtags returns the expected JSON
// when the database query succeeds.
func TestHashtagController_GetHashtags_Success(t *testing.T) {
	// Save the original function and restore it after the test.
	originalGetHashtagsFromDBB := hashtagcontroller.GetHashtagsFromDBB
	defer func() {
		hashtagcontroller.GetHashtagsFromDBB = originalGetHashtagsFromDBB
	}()

	// Override GetHashtagsFromDBB with a mock function.
	hashtagcontroller.GetHashtagsFromDBB = func(conn *sql.DB, hashtag string) ([]map[string]interface{}, error) {
		return []map[string]interface{}{
			{"name": hashtag, "match_count": 42},
		}, nil
	}

	// Create a new HashtagController and assign the fake database connection.
	hc := hashtagcontroller.NewHashtagController()
	hc.Database = &FakeDatabase{} // Assigning the mocked database

	// Execute the method under test.
	jsonData, err := hc.GetHashtags("test")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Parse the JSON response.
	var results []map[string]interface{}
	if err := json.Unmarshal(jsonData, &results); err != nil {
		t.Fatalf("error unmarshaling JSON: %v", err)
	}

	// Validate the response.
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got: %d", len(results))
	}
	if results[0]["name"] != "test" {
		t.Errorf("expected hashtag 'test', got: %v", results[0]["name"])
	}
	if results[0]["match_count"] != 42 {
		t.Errorf("expected match_count 42, got: %v", results[0]["match_count"])
	}
}

// TestHashtagController_GetHashtags_DBError verifies that GetHashtags returns an error
// when the underlying database query fails.
func TestHashtagController_GetHashtags_DBError(t *testing.T) {
	// Save the original function and restore it after the test.
	originalGetHashtagsFromDB := hashtagcontroller.GetHashtagsFromDB
	defer func() { hashtagcontroller.GetHashtagsFromDBB = originalGetHashtagsFromDB }()

	// Override GetHashtagsFromDB with a stub that uses *sql.DB as its first parameter.
	hashtagcontroller.GetHashtagsFromDBB = func(conn *sql.DB, hashtag string) ([]map[string]interface{}, error) {
		return []map[string]interface{}{
			{"name": hashtag, "match_count": 42},
		}, nil
	}

	hc := hashtagcontroller.NewHashtagController()
	hc.Database = &FakeDatabase{}

	// Call GetHashtags expecting an error.
	_, err := hc.GetHashtags("test")
	if err == nil {
		t.Fatalf("expected an error, got none")
	}
	if err.Error() != "db error" {
		t.Errorf("expected error 'db error', got: %v", err)
	}
}
