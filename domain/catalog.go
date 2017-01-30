package domain

import "time"

type (
	Order struct {
		Model
		UserID          uint       `json:"-"`
		StatusID        uint       `json:"-"`
		PaymentMethodID uint       `json:"-"`
		PatmentDetails  string     `json:"-"`
		Total           float64    `json:"total"`
		CustomerNote    string     `json:"customerNote"`
		DeliveryTime    *time.Time `json:"deliveryTime"`

		Status        *OrderStatus   `json:"status"`
		Address       *OrderAddress  `json:"address"`
		PaymentMethod *PaymentMethod `json:"paymentMethod"`
		Products      []OrderProduct `json:"items"`
	}

	OrderProduct struct {
		Model
		OrderID   uint    `json:"-"`
		ProductID uint    `json:"-"`
		Qty       uint    `json:"qty"`
		Price     float64 `json:"price"`
		Total     float64 `json:"total"`
		TaxRate   float64 `json:"taxRate"`
		Options   string  `json:"options"`

		Product *Product `json:"product"`
	}

	OrderAddress struct {
		Model
		OrderID uint `json:"-"`
		AddressBody
	}

	OrderHistory struct {
		Model
		OrderID  uint
		UserID   uint
		StatusID uint16
		Note     string
	}

	OrderStatus struct {
		Model
		Name        string `json:"name"`
		Description string `json:"description"`
		SortNumber  uint16 `json:"-"`
		Status      bool   `json:"-"`
	}

	PaymentMethod struct {
		Model
		Name        string `json:"name"`
		Description string `json:"description"`
		SortNumber  uint16 `json:"-"`
		Status      bool   `json:"-"`
	}

	AddressBody struct {
		Name        string `json:"name"`
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		Tel         string `json:"tel"`
		Tel2        string `json:"tel2"`
		Email       string `json:"email"`
		Address     string `json:"address"`
		City        string `json:"city"`
		District    string `json:"district"`
		Description string `json:"description"`
	}

	Address struct {
		Model
		UserID  uint `json:"-"`
		Default bool `json:"default"`
		AddressBody
	}

	Product struct {
		Model       `storm:"inline"`
		Title       string   `json:"title" fako:"product_name"`
		Description string   `json:"description" gorm:"size:1024" fako:"paragraph"`
		Price       *float64 `json:"price"`
		IsActive    *bool    `json:"isActive"`

		Categories []*Category `gorm:"many2many:pivot_product_category" json:"categories,omitempty"`
		Image      *Image      `json:"defaultImage,omitempty"`
		ImageID    uint        `json:"-"`
	}

	Category struct {
		Model
		Title       string `json:"title" fako:"title"`
		Description string `json:"description" gorm:"size:1024" fako:"paragraph"`
		IsActive    bool   `json:"isActive"`

		Image    *Image    `json:"image,omitempty"`
		ImageID  uint      `json:"-"`
		Products []Product `gorm:"many2many:pivot_product_category" json:"products,omitempty"`
	}

	Image struct {
		Model
		PublicID     string `json:"publicId" gorm:"unique_index"`
		ResourceType string `json:"resourceType"`
	}
)

func (o *Order) SetTotal() {
	for _, op := range o.Products {
		o.Total += op.Total
	}
}

func (op *OrderProduct) GetTotal() float64 {
	return float64(op.Qty) * op.Price
}

func (op *OrderProduct) SetTotal() {
	op.Total = float64(op.Qty) * op.Price
}

func (p *Product) AddCategory(c *Category) {
	for _, v := range p.Categories {
		if c.ID == v.ID {
			return
		}
	}
	p.Categories = append(p.Categories, c)
}
