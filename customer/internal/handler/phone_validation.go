package handler

import (
	"errors"
	"strings"
)

// ValidatePhoneNumber  phone number is exactly 10 digits
func ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return errors.New("phone number is required")
	}

	// Remove common separators:
	PhoneNum := strings.ReplaceAll(phone, " ", "")
	PhoneNum = strings.ReplaceAll(PhoneNum, "-", "")
	PhoneNum = strings.ReplaceAll(PhoneNum, "(", "")
	PhoneNum = strings.ReplaceAll(PhoneNum, ")", "")
	PhoneNum = strings.ReplaceAll(PhoneNum, "+", "")

	if len(PhoneNum) != 10 {
		return errors.New("phone number must be exactly 10 digits")
	}

	for _, ch := range PhoneNum {
		if ch < '0' || ch > '9' {
			return errors.New("phone number must contain only digits")
		}
	}
	return nil
}
