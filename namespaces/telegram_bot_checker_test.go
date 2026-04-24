package namespaces

import (
	"testing"
)

func TestTelegramBotValidateName(t *testing.T) {
	checker := TelegramBotChecker{}

	// Valid name
	err := checker.ValidateName("validbot")
	if err != nil {
		t.Errorf("expected no error for valid name, but got %v", err)
	}

	err = checker.ValidateName("valid")
	if err == nil {
		t.Errorf("without bot at the end invalid %v", err)
	}

	err = checker.ValidateName("vbot")
	if err == nil {
		t.Errorf("to short %v", err)
	}
}
