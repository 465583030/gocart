package gormdb

import (
	"github.com/alioygur/gocart/domain"
	"github.com/alioygur/gocart/engine"
	"github.com/jinzhu/gorm"
)

type (
	userRepository struct {
		sess *gorm.DB
		notFoundErrChecker
	}
)

func NewUserRepository(sess *gorm.DB) engine.UserRepository {
	return &userRepository{sess: sess}
}

func (ur *userRepository) Add(u *domain.User) error {
	return handleErr(ur.sess.Create(u).Error)
}
func (ur *userRepository) OneBy(f []*engine.Filter) (*domain.User, error) {
	var u domain.User
	return &u, handleErr(translateFilter(ur.sess, f).First(&u).Error)
}

func (ur *userRepository) ExistsBy(f []*engine.Filter) (bool, error) {
	var n uint
	err := translateFilter(ur.sess.Table("users"), f).Count(&n).Error
	return n > 0, handleErr(err)
}

func (ur *userRepository) Update(u *domain.User) error {
	return handleErr(ur.sess.Model(u).Update(u).Error)
}
