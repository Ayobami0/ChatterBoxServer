package main

import (
	"log"
	"os"

	"github.com/Ayobami0/chatter_box_server/src/model"
	"github.com/Ayobami0/chatter_box_server/src/server"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := os.Getenv("ENVIRONMENT")

	log_file := os.Getenv("LOG_FILE")

	file, err := os.Create(log_file)

	if err != nil {
		log.Fatal("unable to create log file")
	}

	log.SetOutput(file)

	var log_level logger.LogLevel

	if env != "PRODUCTION" {
		log_level = logger.Info
	} else {
		log_level = logger.Error
	}

	dbLogger := logger.New(
		log.Default(),
		logger.Config{
			LogLevel: log_level,
		},
	)

	var gorm_db *gorm.DB

	if env != "PRODUCTION" {
		gorm_db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
			Logger: dbLogger,
		})
	} else {
		gorm_db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
			Logger: dbLogger,
		})
	}

	if err != nil {
		panic("failed to connect database")
	}

	gorm_db.AutoMigrate(&model.Conversation{}, &model.Message{}, &model.User{}, &model.Request{})

	dbs := map[string]interface{}{
		"gorm": gorm_db,
	}

	config := server.NewConfig(dbs)

	a, err := server.App(config)

	if err != nil {
		log.Fatal(err)
	}

	a.Start(":1323")
}
