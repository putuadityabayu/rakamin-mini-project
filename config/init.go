package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ENV                       string
	JWT_ISSUER                = "https://rakamin.com"
	APPLICATION_NAME          = "Rakamin Mini Project"
	JWT_SIGNATURE_KEY         = []byte("the secret of kalimdor")
	LOGIN_EXPIRATION_DURATION = time.Duration(1) * time.Hour
	JWT_SIGNING_METHOD        = jwt.SigningMethodHS256
	DB                        *gorm.DB
)

func setupDatabase() {
	var err error
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")

	URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	DB, err = gorm.Open(mysql.Open(URL), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err.Error())
	}
}

func Initialization() {
	jwtToken := os.Getenv("JWT_SIGNATURE_KEY")
	JWT_SIGNATURE_KEY = []byte(jwtToken)
	setupDatabase()
}
