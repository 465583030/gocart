package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alioygur/gocart/adapters/web"
	"github.com/alioygur/gocart/engine"
	"github.com/alioygur/gocart/providers"
	"github.com/alioygur/gocart/providers/gormdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	sess, err := gorm.Open("mysql", "root:ali@/gocart?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	sess.LogMode(true)

	// Setup storage factory
	var sf engine.StorageFactory
	sf = gormdb.NewStorage(sess)

	// Setup service dependencies
	var (
		validator  engine.Validator
		mailSender engine.MailSender
		jwt        engine.JWTSignParser
	)

	validator = providers.NewValidator()
	mailSender = providers.NewFakeMail()
	jwt = providers.NewJWT()

	f := engine.New(sf, mailSender, validator, jwt)

	log.Printf("server starting port: %s", os.Getenv("PORT"))
	if err := http.ListenAndServe(":5000", web.NewWebAdapter(f, nil)); err != nil {
		log.Fatal(err)
	}
}
