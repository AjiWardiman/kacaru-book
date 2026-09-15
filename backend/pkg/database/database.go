// Package database bertugas membuka dan mengelola koneksi ke database
// PostgreSQL yang dipakai seluruh aplikasi menggunakan ORM GORM.
package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"tourism-booking/internal/config"
)

// Connect membuka koneksi (connection pool) ke database PostgreSQL
// menggunakan GORM berdasarkan pengaturan di cfg.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	// Susun connection string (DSN) sesuai format yang dipahami driver PostgreSQL GORM.
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database GORM: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil instance sql.DB dari GORM: %w", err)
	}

	// Batasi jumlah koneksi connection pool.
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database GORM (ping gagal): %w", err)
	}

	return db, nil
}
