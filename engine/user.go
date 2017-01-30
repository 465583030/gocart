package engine

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"sync"

	"github.com/alioygur/gocart/domain"
)

type (
	User interface {
		GetFromJWT(tokenStr string) (*domain.User, error)
		GenJWT(*domain.User) (string, error)
		GenPasswordResetToken(*domain.User) (string, error)
		Login(*LoginRequest) (*domain.User, error)
		Register(*RegisterRequest) (*domain.User, error)
		SendPasswordResetMail(*ForgotPasswordRequest) error
		ResetPassword(*ResetPasswordRequest) error
		Show(*ShowUserRequest) (*domain.User, error)
		Update(*UpdateUserRequest) error
	}

	user struct {
		jwt       JWTSignParser
		repo      UserRepository
		mailer    Mailer
		validator Validator
	}

	LoginRequest struct {
		Email    string
		Password string
	}

	RegisterRequest struct {
		Email     string
		Password  string
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		IsActive  *bool  `json:"-"`
	}

	ForgotPasswordRequest struct {
		Email   string
		BaseURL string
	}

	ResetPasswordRequest struct {
		Email    string
		Token    string
		Password string
	}

	ShowUserRequest struct {
		ID uint
	}

	UpdateUserRequest struct {
		ID uint `json:"-"`
		RegisterRequest
	}
)

var (
	userInstance User
	userOnce     sync.Once
)

func (f *factory) NewUser() User {
	userOnce.Do(func() {
		userInstance = &user{
			jwt:       f.jwt,
			repo:      f.NewUserRepository(),
			mailer:    f.NewMail(),
			validator: f.v,
		}
	})
	return userInstance
}

func (u *user) Login(r *LoginRequest) (*domain.User, error) {
	f := []*Filter{NewFilter("email", Equal, r.Email)}
	usr, err := u.repo.OneBy(f)
	if err != nil {
		if err == ErrNoRows {
			return nil, ErrWrongCredentials
		}
		return nil, err
	}

	if !usr.IsCredentialsVerified(r.Password) {
		return nil, ErrWrongCredentials
	}

	if !*usr.IsActive {
		return nil, ErrInActiveUser
	}
	return usr, nil
}

func (u *user) Register(r *RegisterRequest) (*domain.User, error) {
	// validation
	if err := u.validator.CheckEmail(r.Email); err != nil {
		return nil, err
	}
	if err := checkPassowrd(u.validator, r.Password); err != nil {
		return nil, err
	}

	// check for email
	f := []*Filter{NewFilter("email", Equal, r.Email)}
	exists, err := u.repo.ExistsBy(f)
	if err != nil {
		return nil, err
	} else if exists {
		return nil, ErrEmailExists
	}

	if r.IsActive == nil {
		r.IsActive = boolPtr(false)
	}

	var usr domain.User
	usr.FirstName = r.FirstName
	usr.LastName = r.LastName
	usr.Email = r.Email
	usr.SetPassword(r.Password)
	usr.IsActive = r.IsActive
	usr.IsAdmin = boolPtr(false)

	if err := u.repo.Add(&usr); err != nil {
		return nil, err
	}

	if err := u.mailer.SendWelcomeMail(usr.Email); err != nil {
		return nil, err
	}

	return &usr, nil
}

func (u *user) GetFromJWT(tokenStr string) (*domain.User, error) {
	claims, err := u.jwt.Parse(tokenStr, os.Getenv("SECRET_KEY"))
	if err != nil {
		return nil, err
	}
	idStr, ok := claims["userID"].(string)
	if !ok {
		return nil, fmt.Errorf("userID can't get from token claims, token: %s", tokenStr)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	return u.repo.OneBy(idFilter(uint(id)))
}

func (u *user) GenJWT(du *domain.User) (string, error) {
	id := strconv.Itoa(int(du.ID))
	claims := map[string]interface{}{
		"userID": id,
		"exp":    time.Now().Add(time.Hour * 6).Unix(),
	}

	return u.jwt.Sign(claims, os.Getenv("SECRET_KEY"))
}

func (u *user) GenPasswordResetToken(du *domain.User) (string, error) {
	claims := map[string]interface{}{"email": du.Email, "exp": time.Now().Add(time.Hour * 5).Unix()}
	return u.jwt.Sign(claims, du.Password)
}

func (u *user) CheckPasswordResetToken(tokenStr string, du *domain.User) error {
	claims, err := u.jwt.Parse(tokenStr, du.Password)
	if err != nil {
		return err
	}

	email, ok := claims["email"].(string)
	if !ok {
		return fmt.Errorf("email can't get from token claims, token: %s", tokenStr)
	}

	if email != du.Email {
		return fmt.Errorf("token's email and user's email aren't equal: %s", tokenStr)
	}

	return nil
}

func (u *user) SendPasswordResetMail(r *ForgotPasswordRequest) error {
	if r.BaseURL == "" {
		r.BaseURL = os.Getenv("PASSWORD_RESET_URL")
	}
	resetURL, err := url.Parse(r.BaseURL)
	if err != nil {
		return err
	}

	f := []*Filter{NewFilter("email", Equal, r.Email)}
	usr, err := u.repo.OneBy(f)
	if err != nil {
		return err
	}

	tokenString, err := u.GenPasswordResetToken(usr)
	if err != nil {
		return err
	}

	q := resetURL.Query()
	q.Set("token", tokenString)
	resetURL.RawQuery = q.Encode()

	return u.mailer.SendPasswordResetLink(r.Email, resetURL.String())
}

func (u *user) ResetPassword(r *ResetPasswordRequest) error {
	// validation
	if err := checkPassowrd(u.validator, r.Password); err != nil {
		return err
	}

	f := []*Filter{NewFilter("email", Equal, r.Email)}
	usr, err := u.repo.OneBy(f)
	if err != nil {
		return err
	}

	if err := u.CheckPasswordResetToken(r.Token, usr); err != nil {
		return err
	}

	var uusr domain.User
	uusr.ID = usr.ID
	uusr.SetPassword(r.Password)

	return u.repo.Update(&uusr)
}

func (u *user) Show(r *ShowUserRequest) (*domain.User, error) {
	f := idFilter(r.ID)
	return u.repo.OneBy(f)
}

func (u *user) Update(r *UpdateUserRequest) error {
	usr, err := u.repo.OneBy(idFilter(r.ID))
	if err != nil {
		return err
	}

	var user domain.User
	user.ID = r.ID
	if r.FirstName != "" {
		user.FirstName = r.FirstName
	}
	if r.LastName != "" {
		user.LastName = r.LastName
	}
	if r.Email != "" && r.Email != usr.Email {
		// validate
		if err := u.validator.CheckEmail(r.Email); err != nil {
			return err
		}
		// someone else exists with this email?
		var f []*Filter
		f = append(f, NewFilter("email", Equal, r.Email))
		exists, err := u.repo.ExistsBy(f)
		if err != nil {
			return err
		}
		if exists {
			return ErrEmailExists
		}
		user.Email = r.Email
	}
	if r.Password != "" {
		if err := checkPassowrd(u.validator, r.Password); err != nil {
			return err
		}
		user.SetPassword(r.Password)
	}

	return u.repo.Update(&user)
}

func checkPassowrd(v Validator, p string) error {
	return v.CheckStringLen(p, 4, 8, "Password")
}
