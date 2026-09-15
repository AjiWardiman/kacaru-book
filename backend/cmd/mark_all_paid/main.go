package main

import (
	"fmt"
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=admin dbname=tourism_booking port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB conn error: %v", err)
	}

	res := db.Exec("UPDATE bookings SET status = 'paid' WHERE status = 'pending'")
	fmt.Printf("Updated pending bookings to paid: %d\n", res.RowsAffected)
}
