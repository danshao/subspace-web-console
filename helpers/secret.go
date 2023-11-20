package helpers

import (
	"math/rand"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	r        = rand.New(rand.NewSource(time.Now().UnixNano())) // seed random generation
	mustHave = []func(rune) bool{
		unicode.IsUpper,
		unicode.IsLower,
		unicode.IsDigit,
	}
)

const (
	alphabetLower = "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numerals      = "0123456789"
	passwordChars = alphabetLower + alphabetUpper + numerals
)

// GeneratePassword returns a random password of given length
func GeneratePassword(strlen int) string {
	var (
		result = ""
		valid  = false
	)

	for valid == false {
		for i := 0; i < strlen; i++ {
			index := r.Intn(len(passwordChars))
			result += passwordChars[index : index+1]
		}

		// validate password conformity
		if passwordOK(result) {
			valid = true
		} else {
			result = ""
		}
	}

	return result
}

// HashPassword generates a bcrypt hashed password using a default cost of 2^10
func HashPassword(password string) (string, error) {
	//The cost too large will spend more time, Go default is 10 (2^10), Ruby devise is 11 (2^11)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash checks a password against its bcrypt hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

/*
 *	Private Functions
 */
func passwordOK(p string) bool {
	for _, testRune := range mustHave {
		found := false
		for _, r := range p {
			if testRune(r) {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
