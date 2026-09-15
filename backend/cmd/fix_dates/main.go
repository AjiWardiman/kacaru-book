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

	res := db.Exec("UPDATE bookings SET booked_at = created_at WHERE booked_at < '2000-01-01'")
	fmt.Printf("Updated zero booked_at rows: %d\n", res.RowsAffected)
}
