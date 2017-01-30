package domain

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type (
	contextKey string

	// User model
	User struct {
		Model
		FirstName string `json:"firstName" fako:"first_name"`
		LastName  string `json:"lastName" fako:"last_name"`
		Email     string `json:"email" fako:"email_address" storm:"unique"`
		Password  string `json:"password" fako:"simple_password"`
		IsActive  *bool  `json:"isActive"`
		IsAdmin   *bool  `json:"isAdmin"`
	}
)

var userContextKey contextKey = "user"

// SetPassword sets user's password
func (u *User) SetPassword(p string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	u.Password = string(hashedPassword)
}

// IsCredentialsVerified matches given password with user's password
func (u *User) IsCredentialsVerified(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}

func (u *User) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// UserFromContext gets user from context
func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userContextKey).(*User)
	return u, ok
}

// UserMustFromContext gets user from context. if can't make panic
func UserMustFromContext(ctx context.Context) *User {
	u, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		panic("user can't get from request's context")
	}
	return u
}
