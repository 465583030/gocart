package gormdb

import (
	"sync"

	"github.com/alioygur/gocart/engine"
	"github.com/jinzhu/gorm"
)

type (
	storageFactory struct {
		sess *gorm.DB
	}
)

var (
	userRepositoryInstance    engine.UserRepository
	userRepositoryOnce        sync.Once
	imageRepositoryInstance   engine.ImageRepository
	imageRepositoryOnce       sync.Once
	catalogRepositoryInstance engine.CatalogRepository
	catalogRepositoryOnce     sync.Once
)

func NewStorage(session *gorm.DB) engine.StorageFactory {
	return &storageFactory{session}
}

func (sf *storageFactory) NewCatalogRepository() engine.CatalogRepository {
	catalogRepositoryOnce.Do(func() {
		catalogRepositoryInstance = NewCatalogRepository(sf.sess)
	})
	return catalogRepositoryInstance
}

func (sf *storageFactory) NewImageRepository() engine.ImageRepository {
	imageRepositoryOnce.Do(func() {
		imageRepositoryInstance = newImageRepository()
	})
	return imageRepositoryInstance
}

func (sf *storageFactory) NewUserRepository() engine.UserRepository {
	userRepositoryOnce.Do(func() {
		userRepositoryInstance = NewUserRepository(sf.sess)
	})
	return userRepositoryInstance
}
