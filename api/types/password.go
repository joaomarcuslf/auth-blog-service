package types

import (
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var BCRYPT_PATTERN = "\\A\\$2a?\\$\\d\\d\\$[./0-9A-Za-z] {53}"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func IsValidPassword(password string) bool {
	matched, _ := regexp.MatchString(BCRYPT_PATTERN, password)

	return matched
}

type Password struct {
	Hash string
}

func (p *Password) UnmarshalJSON(input []byte) error {
	strInput := strings.Trim(string(input), `"`)

	if IsValidPassword(strInput) {
		p.Hash = strInput
		return nil
	}

	hash, _ := HashPassword(strInput)

	p.Hash = hash
	return nil
}
