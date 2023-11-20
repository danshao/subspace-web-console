package helpers

import (
	"testing"
	"fmt"
)

func TestHashPassword(t *testing.T) {
	PASSWORD := "MyAwesomePassword"
	PASSWORD_WRONG := "WrongPassword"
	fmt.Println(PASSWORD)
	passwordHash, _ := HashPassword(PASSWORD)
	fmt.Println(passwordHash)
	if !CheckPasswordHash(PASSWORD, passwordHash) {
		t.Error("Password not match.")
	}

	if CheckPasswordHash(PASSWORD_WRONG, passwordHash) {
		t.Error("Password should not match.")
	}
}