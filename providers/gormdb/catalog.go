package gormdb

import (
	"github.com/alioygur/gocart/domain"

	"github.com/alioygur/gocart/engine"
	"github.com/jinzhu/gorm"
)

type (
	catalogRepository struct {
		sess *gorm.DB
		notFoundErrChecker
	}
)

func NewCatalogRepository(sess *gorm.DB) engine.CatalogRepository {
	return &catalogRepository{sess: sess}
}

func (c *catalogRepository) AddProduct(p *domain.Product) error {
	return handleErr(c.sess.Create(p).Error)
}

func (c *catalogRepository) OneProductBy(f []*engine.Filter) (*domain.Product, error) {
	var p domain.Product
	return &p, handleErr(translateFilter(c.sess, f).First(&p).Error)
}

func (c *catalogRepository) FindProducts(q *engine.Query) ([]*domain.Product, error) {
	var ps []*domain.Product
	return ps, handleErr(translateQuery(c.sess, q).Find(&ps).Error)
}

func (c *catalogRepository) FindProductsInCategories(cs []uint, q *engine.Query) ([]*domain.Product, error) {
	var uids []uint
	// find products in those categories
	if err := c.sess.Table("pivot_product_category").
		Where("category_id in(?)", cs).
		Group("product_id").
		Pluck("product_id", &uids).Error; err != nil {
		return nil, handleErr(err)
	}

	if len(uids) == 0 {
		return nil, nil
	}

	q.Filters = append(q.Filters, engine.NewFilter("id", engine.In, uids))

	return c.FindProducts(q)
}

func (c *catalogRepository) UpdateProduct(p *domain.Product) error {
	return handleErr(c.sess.Model(p).Update(p).Error)
}

func (c *catalogRepository) DeleteProductBy(f []*engine.Filter) error {
	return handleErr(translateFilter(c.sess, f).Delete(&domain.Product{}).Error)
}

func (c *catalogRepository) FindCategories(q *engine.Query) ([]*domain.Category, error) {
	var cs []*domain.Category
	return cs, handleErr(translateQuery(c.sess, q).Find(&cs).Error)
}
