package handler


import (
	"errors"
	"strings"
)

// Valid currencies for accounts.
var allowedCurrencies = map[string]struct{}{
	"saudi": {},
	"dubia": {},
	"india": {},
}

var allowedAccountTypes = map[string]struct{}{
	"savings":  {},
	"current":  {},
	"business": {},
}

// ValidateAccountFields checks account fields before creation.
func ValidateAccountFields(accountType, currency, status string, availableBalance, pendingBalance int64) error {
	if _, ok := allowedAccountTypes[strings.ToLower(accountType)]; !ok {
		return errors.New("account type must be one of: savings, current, business")
	}
	if _, ok := allowedCurrencies[strings.ToLower(currency)]; !ok {
		return errors.New("currency must be one of: saudi, dubia, india")
	}
	if status == "" {
		return errors.New("account status is required")
	}
	if availableBalance < 0 {
		return errors.New("available balance cannot be negative")
	}
	if pendingBalance < 0 {
		return errors.New("pending balance cannot be negative")
	}
	return nil
}
