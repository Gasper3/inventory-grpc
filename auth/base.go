package auth

import "golang.org/x/crypto/bcrypt"

type User struct {
	Username       string
	HashedPassword string
	Role           string
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
}
