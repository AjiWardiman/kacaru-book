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

	var bookings []map[string]interface{}
	db.Raw("SELECT id, booking_code, status, total_price, booked_at FROM bookings ORDER BY id DESC LIMIT 10").Scan(&bookings)
	for _, b := range bookings {
		fmt.Printf("%+v\n", b)
	}
}
