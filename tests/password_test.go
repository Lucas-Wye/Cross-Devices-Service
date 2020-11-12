package test

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestAuth(t *testing.T) {
	secret := "$2a$10$zVeDUQ6CdmzQK55iojloiecJEoHz2qW7AMvIb19JXQ/kRfRFe7s.O"
	password := "hello"
	if bcrypt.CompareHashAndPassword([]byte(secret), []byte(password)) != nil {
		t.Fail()
	}
}
