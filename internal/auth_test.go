package auth

import (
	"fmt"
	"testing"
)

func TestHasshPassword(t *testing.T) {
	password := "testPassword"

	hassh, err := HashPassword(password)
	if err != nil {
		fmt.Printf("Error :%s", err)
	}
	fmt.Printf("Hasshed Password: %s\n", hassh)
}

func TestPasswordLong(t *testing.T){
	password := "testPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPassword"
	_, err := HashPassword(password)
	if err == nil{ // If the hasshing succedding err would be nil, so failing the test
		t.Errorf("Should not have worked")
	}
}

func TestCheckPasswordHassh(t *testing.T) {
	hasshedPassword := "$2a$10$8u.iZ8zIGCJe1j3NIKI0eua3iS0JUV06qVP9pj/34p14dn36r1ZaK"
	password := "testPassword"

	err := CheckPasswordHash(password, hasshedPassword)
	if err != nil {
		fmt.Printf("Error :%s\n", err)
	}
}
