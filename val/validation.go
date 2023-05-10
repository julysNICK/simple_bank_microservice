package val

import (
	"fmt"
	"net/mail"
	"regexp"
)


var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)


func ValidateString(value string, minLength int, maxLength int) error {
	if len(value) < minLength || len(value) > maxLength {
		return fmt.Errorf("length must be between %d and %d", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 50); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("username must be alphanumeric example: july_1997")
	}

	return nil

}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 50); err != nil {
		return err
	}

	if !isValidFullName(value) {
		return fmt.Errorf("full_name contain only letters and spaces")
	}

	return nil

}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 50)
}

func ValidEmail(value string) error {
	if err := ValidateString(value, 3, 50); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("invalid email address")
	}

	return nil
}

