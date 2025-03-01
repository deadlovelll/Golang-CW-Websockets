package basecontroller

import (
	"os"
	"os/exec"
	"testing"

	basecontroller "messenger_engine/controllers/base_controller"
	"messenger_engine/modules/database/database"
)

// MockDatabase simulates a database instance for testing.
type MockDatabase struct{}

// TestNewBaseController_Success tests successful initialization.
func TestNewBaseController_Success(t *testing.T) {
	mockDB := &database.Database{} // Create a mock database instance
	baseCtrl := basecontroller.NewBaseController(mockDB)

	if baseCtrl.Database != mockDB {
		t.Errorf("Expected database instance to be assigned, got nil or different instance")
	}
}

// TestNewBaseController_NilDatabase tests behavior when passing nil.
func TestNewBaseController_NilDatabase(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		_ = basecontroller.NewBaseController(nil)
		return
	}

	cmd := exec.Command("go", "test", "-run=TestNewBaseController_NilDatabase")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	if exitError, ok := err.(*exec.ExitError); ok && !exitError.Success() {
		// Expected failure
		return
	}
	t.Fatalf("Expected process to exit with failure, but it did not")
}

// TestGetDatabase tests GetDatabase method.
func TestGetDatabase(t *testing.T) {
	mockDB := &database.Database{}
	baseCtrl := basecontroller.NewBaseController(mockDB)

	if baseCtrl.GetDatabase() != mockDB {
		t.Errorf("GetDatabase() returned unexpected database instance")
	}
}
