package main

import (
	"BMSTU_IU5_53B_rip/internal/app/ds"
	"BMSTU_IU5_53B_rip/internal/app/dsn"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	env := dsn.FromEnv()
	fmt.Println("!   !   !   DB Connection String:", env)

	db, err := gorm.Open(postgres.Open(env), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	if err := db.AutoMigrate(
		&ds.DeliveryItem{},
		&ds.DeliveryRequest{},
		&ds.User{},
		&ds.Item_request{},
	); err != nil {
		fmt.Println("Migration error:", err)
		panic("cant migrate db")
	}
}
