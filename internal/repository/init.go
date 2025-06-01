package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"startfront-backend/pkg/config"
)

var db *sqlx.DB

func InitDB(cfg *config.DBConfig) {
	fmt.Printf("✅ Loaded DB config: %+v\n", config.DB)

	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		cfg.User, cfg.Password, cfg.DBName, cfg.Host, cfg.Port, cfg.SSLMode,
	)

	var err error
	db, err = sqlx.Open("postgres", dsn)
	if err != nil {
		panic("❌ Failed to connect to database: " + err.Error())
	}

	if err = db.Ping(); err != nil {
		panic("❌ Failed to ping database: " + err.Error())
	}

	fmt.Println("✅ Connected to database:", cfg.DBName)
}
