package namespaces

import (
	"testing"
)

func TestValidateName(t *testing.T) {
	checker := TelegramChecker{}

	// Valid name
	err := checker.ValidateName("valid_name")
	if err != nil {
		t.Errorf("expected no error for valid name, but got %v", err)
	}

	// Name too short
	err = checker.ValidateName("a")
	if err == nil || err.Error() != "Name must be at least 5 characters long" {
		t.Errorf("expected error for short name, but got %v", err)
	}

	// Name with invalid characters
	err = checker.ValidateName("invalid!")
	if err == nil || err.Error() != "Name may consist only of a-z, 0-9, and underscores" {
		t.Errorf("expected error for invalid characters, but got %v", err)
	}

	// Name with valid length but invalid characters
	err = checker.ValidateName("name$with$")
	if err == nil || err.Error() != "Name may consist only of a-z, 0-9, and underscores" {
		t.Errorf("expected error for invalid characters in name, but got %v", err)
	}
}
