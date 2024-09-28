package usecases

import (
	"testing"
)

func TestNewHealthCheckInteractor(t *testing.T) {
	// Given
	interactor := NewHealthCheckInteractor()

	// Then
	if interactor == nil {
		t.Errorf("nil이 아님을 기대했으나 nil이 반환")
	}
}

func TestPerformHealthCheck(t *testing.T) {
	// Given
	interactor := NewHealthCheckInteractor()

	expectedStatus := "Healthy"
	expectedMessage := "Success"

	// When
	status := interactor.PerformHealthCheck()

	// Then
	if status.Status != expectedStatus {
		t.Errorf("상태가 '%s'이길 기대했으나 '%s'가 반환.", expectedStatus, status.Status)
	}
	if status.Message != expectedMessage {
		t.Errorf("메시지가 '%s'이길 기대했으나 '%s'가 반환.", expectedMessage, status.Message)
	}
}
