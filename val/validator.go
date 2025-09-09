package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	inValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString

)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("string length must be between %d and %d", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value){
		return fmt.Errorf("username must only contain digits,lowercase alphanumeric characters and underscores")
	}
	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err!= nil {
        return err
    }
    if !inValidFullName(value){
        return fmt.Errorf("full name must only contain alphabetic characters and spaces")
    }
    return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value,5,100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value,3,200); err != nil{
		return err
	}
	if _,err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("invalid email address")
	}
	return nil
}


func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}


func ValidateSecretCode(value string) error {
	return ValidateString(value, 32, 128)
}