package gormdb

import (
	"errors"

	"github.com/alioygur/gocart/engine"
	"github.com/jinzhu/gorm"
)

type (
	notFoundErrChecker struct{}
)

var (
	errNotFound = errors.New("not found")
)

func (e *notFoundErrChecker) IsNotFoundErr(err error) bool {
	return err == gorm.ErrRecordNotFound
}

func handleErr(err error) error {
	if err == gorm.ErrRecordNotFound {
		return engine.ErrNoRows
	}
	return err
}
