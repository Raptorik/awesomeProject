package apperror_test

import (
	"errors"
	"testing"

	"awesomeProject/internal/apperror"
)

func TestNewAppError(t *testing.T) {
	err := errors.New("some error")
	msg := "some message"
	devMsg := "some developer message"
	code := "US-0000001"

	appErr := apperror.NewAppError(err, msg, devMsg, code)

	if appErr.Err != err {
		t.Errorf("Expected Err to be %v, but got %v", err, appErr.Err)
	}
	if appErr.Message != msg {
		t.Errorf("Expected Message to be %v, but got %v", msg, appErr.Message)
	}
	if appErr.DeveloperMessage != devMsg {
		t.Errorf("Expected DeveloperMessage to be %v, but got %v", devMsg, appErr.DeveloperMessage)
	}
	if appErr.Code != code {
		t.Errorf("Expected Code to be %v, but got %v", code, appErr.Code)
	}
}

func TestSystemError(t *testing.T) {
	err := errors.New("some error")
	expectedCode := "US-000000"

	appErr := apperror.SystemError(err)

	if appErr.Err != err {
		t.Errorf("Expected Err to be %v, but got %v", err, appErr.Err)
	}
	if appErr.Message != "internal system error" {
		t.Errorf("Expected Message to be 'internal system error', but got %v", appErr.Message)
	}
	if appErr.DeveloperMessage != err.Error() {
		t.Errorf("Expected DeveloperMessage to be %v, but got %v", err.Error(), appErr.DeveloperMessage)
	}
	if appErr.Code != expectedCode {
		t.Errorf("Expected Code to be %v, but got %v", expectedCode, appErr.Code)
	}
}
