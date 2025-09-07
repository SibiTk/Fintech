package handler

import (
	"errors"
	"strings"
	"time"
	"unicode"
)


func isAlphaSpace(s string) bool {
	for _, r := range s {
		if !(unicode.IsLetter(r) || r == ' ') {
			return false
		}
	}
	return true
}

func ValidateCustomerFields(
	customerNumber, firstName, lastName, dateOfBirth, status, kycStatus string,
) error {

	if strings.TrimSpace(customerNumber) == "" {
		return errors.New("customer number is required")
	}
	if strings.TrimSpace(firstName) == "" {
		return errors.New("first name is required")
	}
	if strings.TrimSpace(lastName) == "" {
		return errors.New("last name is required")
	}
	if strings.TrimSpace(status) == "" {
		return errors.New("status is required")
	}
	if strings.TrimSpace(kycStatus) == "" {
		return errors.New("kyc status is required")
	}

	// First and last names: only letters and spaces
	if !isAlphaSpace(firstName) {
		return errors.New("first name must contain only letters and spaces")
	}
	if !isAlphaSpace(lastName) {
		return errors.New("last name must contain only letters and spaces")
	}

	// Date of Birth: must be in DD-MM-YYYY format if provided
	if dob := strings.TrimSpace(dateOfBirth); dob != "" {
		if _, err := time.Parse("02-01-2006", dob); err != nil {
			return errors.New("date of birth must be in DD-MM-YYYY format")
		}
	}

	// KYC Status: only "active", "inprogress", "failed"
	allowedKYC := map[string]struct{}{
		"active":     {},
		"inprogress": {},
		"failed":     {},
	}
	kycVal := strings.ToLower(strings.TrimSpace(kycStatus))
	if _, ok := allowedKYC[kycVal]; !ok {
		return errors.New("kyc status must be one of: active, inprogress, failed")
	}

	// Status: only "pending", "complete", "failed"
	allowedStatus := map[string]struct{}{
		"pending":  {},
		"complete": {},
		"failed":   {},
	}
	statusVal := strings.ToLower(strings.TrimSpace(status))
	if _, ok := allowedStatus[statusVal]; !ok {
		return errors.New("status must be one of: pending, complete, failed")
	}
	

	return nil
}
