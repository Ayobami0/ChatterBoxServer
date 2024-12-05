package main

import (
	"log"

	"github.com/Ayobami0/chatter_box_server/src/model"
	"github.com/Ayobami0/chatter_box_server/src/server"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	gorm_db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
  })
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
