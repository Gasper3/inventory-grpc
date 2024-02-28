package auth

import (
	"testing"
	"time"
)

const secretKey = "secret-test"

func TestGenerateAndVerify(t *testing.T) {
	user := User{Username: "Test", Role: "test"}
	manager := JwtManager{SecretKey: secretKey, Duration: 5 * time.Minute}

	token, err := manager.GenerateKey(&user)
	if err != nil {
		t.Errorf("Token generation failed: %v", err)
	}

	claims, err := manager.Verify(token)
	if err != nil {
		t.Errorf("Token verification failed: %v", err)
	}

	if claims.Role != user.Role || claims.Username != user.Username {
		t.Errorf("Claims are incorrect: role -> %v got %v , username -> %v got %v", user.Role, claims.Role, user.Username, claims.Username)
	}
}
