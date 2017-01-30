package engine

import (
	"github.com/alioygur/gocart/domain"
)

type (
	notFoundErrChecker interface {
		IsNotFoundErr(error) bool
	}

	UserRepository interface {
		Add(*domain.User) error
		OneBy([]*Filter) (*domain.User, error)
		ExistsBy([]*Filter) (bool, error)
		Update(*domain.User) error
	}

	CatalogRepository interface {
		notFoundErrChecker
		AddProduct(*domain.Product) error
		OneProductBy([]*Filter) (*domain.Product, error)
		FindProducts(*Query) ([]*domain.Product, error)
		FindProductsInCategories([]uint, *Query) ([]*domain.Product, error)
		UpdateProduct(*domain.Product) error
		DeleteProductBy([]*Filter) error
		FindCategories(*Query) ([]*domain.Category, error)
	}

	ImageRepository interface {
		FirstOrInit(string) (*domain.Image, error)
	}

	StorageFactory interface {
		NewUserRepository() UserRepository
		NewCatalogRepository() CatalogRepository
		NewImageRepository() ImageRepository
	}
)
