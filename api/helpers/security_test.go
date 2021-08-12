package helpers

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash1, err := HashPassword("123456")
	hash2, _ := HashPassword("123456")

	if err != nil {
		t.Error("HashPassword 01 failed")
	} else {
		t.Log("HashPassword 01 passed")
	}

	if hash1 != hash2 {
		t.Log("HashPassword 02 passed")
	} else {
		t.Error("HashPassword 02 failed")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	hash1, _ := HashPassword("123456")

	if CheckPasswordHash("123456", hash1) {
		t.Log("CheckPasswordHash 01 passed")
	} else {
		t.Error("CheckPasswordHash 01 failed")
	}

	if !CheckPasswordHash("12345", hash1) {
		t.Log("CheckPasswordHash 01 passed")
	} else {
		t.Error("CheckPasswordHash 01 failed")
	}
}
