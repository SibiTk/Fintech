package handler

import "strings"

func IsValidCurrency(currency string) bool {
	switch strings.ToLower(currency) {
	case "saudi", "dubia", "india":
		return true
	default:
		return false
	}
}