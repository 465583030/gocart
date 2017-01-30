package gormdb

import (
	"testing"

	"github.com/alioygur/gocart/domain"
	"github.com/alioygur/gocart/engine"

	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func openDB() (*gorm.DB, error) {
	return gorm.Open("mysql", "root:ali@/gocart?charset=utf8&parseTime=True&loc=Local")
}

func initCatalogRepository() engine.CatalogRepository {
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	return NewCatalogRepository(db)
}

func TestReflect(t *testing.T) {
	c := initCatalogRepository()
	p := exampleProduct()

	if err := c.AddProduct(p); err != nil {
		log.Fatal(err)
	}
	log.Println(p)
}

type ali struct {
	IsActive *bool
}

func TestAli(t *testing.T) {
}

func exampleProduct() *domain.Product {
	var p domain.Product
	p.Title = "Example product"
	p.Description = "example desc"
	// p.Price = 5.9
	// p.IsActive = true
	return &p
}

type user struct {
	isAdmin *bool
	age     *int
}

type a int

func TestOnce(t *testing.T) {
	var u *user
	log.Println(u)
}
