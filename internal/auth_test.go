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
