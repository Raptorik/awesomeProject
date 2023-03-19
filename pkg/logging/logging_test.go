package logging

import (
	"bytes"
	"testing"
)

type MockWriter struct {
	buffer bytes.Buffer
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	return m.buffer.Write(p)
}
func TestInit(t *testing.T) {}

func TestGetLogger(t *testing.T) {
	logger := GetLogger()
	if logger == nil {
		t.Fatal("GetLogger() returned a nil logger")
	}

	t.Run("log message", func(t *testing.T) {
		mockWriter := &MockWriter{}
		e.Logger.SetOutput(mockWriter)

		testMessage := "Test message"
		logger.Info(testMessage)

		logOutput := mockWriter.buffer.String()
		if logOutput == "" {
			t.Fatal("Logger did not write to the provided writer")
		}

		if !bytes.Contains([]byte(logOutput), []byte(testMessage)) {
			t.Errorf("Logger output does not contain the expected message: %s", testMessage)
		}
	})
}
