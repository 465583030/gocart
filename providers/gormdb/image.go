package gormdb

import (
	"github.com/alioygur/gocart/domain"

	"sync"

	"github.com/alioygur/gocart/engine"
)

type (
	imageRepository struct {
		sync.RWMutex
		images []*domain.Image
	}
)

func newImageRepository() engine.ImageRepository {
	return &imageRepository{}
}

func (r *imageRepository) FirstOrInit(name string) (*domain.Image, error) {
	r.RLock()
	defer r.RUnlock()
	for _, img := range r.images {
		if img.PublicID == name {
			return img, nil
		}
	}
	return &domain.Image{PublicID: name}, nil
}
