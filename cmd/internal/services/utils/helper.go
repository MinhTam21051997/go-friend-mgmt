package utils

import (
	"errors"
	"fmt"
	"regexp"
)

func VaditationEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) < 3 && len(email) > 254 {
		return fmt.Errorf("email length should from %d to %d", 3, 254)
	}

	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}