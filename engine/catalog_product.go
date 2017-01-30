package engine

import "github.com/alioygur/gocart/domain"

type (
	AddProductRequest struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Price       *float64 `json:"price"`
		IsActive    *bool    `json:"isActive"`
		Image       string   `json:"image"`
		Categories  []uint   `json:"categories"`
	}

	ShowProductRequest struct {
		ID uint
		// IncludeInactive include inactive products
		IncludeInactive bool
	}

	ListProductsRequest struct {
		//InCategories list products that are in these categories
		InCategories []uint
		// IncludeInactive include inactive products
		IncludeInactive bool
		Limit           int
		Offset          int
	}

	UpdateProductRequest struct {
		ID uint
		*AddProductRequest
	}

	DeleteProductRequest struct {
		ID uint
	}
)

func (c *catalog) AddProduct(r *AddProductRequest) (*domain.Product, error) {
	// validate
	if err := checkProductTitle(c.validator, r.Title); err != nil {
		return nil, err
	}
	if err := checkProductDesc(c.validator, r.Description); err != nil {
		return nil, err
	}

	var p domain.Product
	p.Title = r.Title
	p.Description = r.Description
	p.Price = r.Price
	if r.IsActive == nil {
		p.IsActive = boolPtr(false)
	}

	if r.Image != "" {
		img, err := c.imgRepo.FirstOrInit(r.Image)
		if err != nil {
			return nil, err
		}
		p.Image = img
	}

	q := NewQuery("").Filter("id", In, r.Categories)

	cats, err := c.repo.FindCategories(q)
	if err != nil {
		return nil, err
	}

	for _, c := range cats {
		p.AddCategory(c)
	}

	return &p, c.repo.AddProduct(&p)
}

func (c *catalog) ShowProduct(r *ShowProductRequest) (*domain.Product, error) {
	f := idFilter(r.ID)
	if r.IncludeInactive == false {
		f = append(f, NewFilter("is_active", Equal, true))
	}

	return c.repo.OneProductBy(f)
}

func (c *catalog) ListProducts(r *ListProductsRequest) ([]*domain.Product, error) {
	q := NewQuery("")
	if r != nil {
		if r.IncludeInactive == false {
			q.Filter("is_active", Equal, true)
		}

		if r.Limit > 0 {
			q.Slice(r.Offset, r.Limit)
		} else {
			q.Slice(0, 100)
		}

		if r.InCategories != nil {
			return c.repo.FindProductsInCategories(r.InCategories, q)
		}
	}

	return c.repo.FindProducts(q)
}

func (c *catalog) UpdateProduct(r *UpdateProductRequest) error {
	// get original product from db
	f := idFilter(r.ID)
	op, err := c.repo.OneProductBy(f)
	if err != nil {
		return err
	}

	var p domain.Product

	if r.Title != "" {
		if err := checkProductTitle(c.validator, r.Title); err != nil {
			return err
		}
		p.Title = r.Title
	}
	if r.Description != "" {
		if err := checkProductDesc(c.validator, r.Description); err != nil {
			return err
		}
		p.Description = r.Description
	}
	if r.Price != nil {
		p.Price = r.Price
	}
	if r.IsActive != nil {
		p.IsActive = r.IsActive
	}

	if r.Image != "" {
		img, err := c.imgRepo.FirstOrInit(r.Image)
		if err != nil {
			return err
		}

		if op.ImageID == 0 || op.ImageID != img.ID {
			p.Image = img
		}
	}

	q := NewQuery("").Filter("id", In, r.Categories)

	cs, err := c.repo.FindCategories(q)
	if err != nil {
		return err
	}

	for _, c := range cs {
		p.AddCategory(c)
	}

	return c.repo.UpdateProduct(&p)
}

func (c *catalog) DeleteProduct(r *DeleteProductRequest) error {
	f := idFilter(r.ID)
	return c.repo.DeleteProductBy(f)
}

func checkProductTitle(v Validator, title string) error {
	return v.CheckStringLen(title, 2, 255, "Title")
}

func checkProductDesc(v Validator, desc string) error {
	return v.CheckStringLen(desc, 2, 1024, "Description")
}
