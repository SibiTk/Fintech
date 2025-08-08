package handler
import (
	"errors"
	"strings"
)
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	email = strings.TrimSpace(email)
	atIdx := strings.Index(email, "@")
	if atIdx <= 0  {
		return errors.New("email must have @ symbol ")
	}
	if len(email) > 20 {
		return errors.New("email is too long, max 20 characters allowed")
	}
return nil
}