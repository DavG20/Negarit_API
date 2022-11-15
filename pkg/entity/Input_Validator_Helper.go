package entity

import (
	"net/mail"
	"strconv"
	"strings"
)

func ValidateUserName(userName string) bool {
	trim := func() bool {
		name := strings.Trim(userName, " ")
		return len(name) < 4
	}

	checkNum := func() bool {
		_, err := strconv.Atoi(userName)
		return err == nil
	}
	if trim() || checkNum() {
		return false
	}
	return true

}

func ValidatePassword(password string) bool {
	pass := strings.Trim(password, " ")
	if len(password) < 8 {
		return false
	}
	if len(pass) == 0 {
		return false
	}
	return true
}

func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}

	_, err := mail.ParseAddress(email)

	return err == nil
}
