package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // Or your driver
	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type Config struct {
	Database DBConfig `yaml:"database"`
}

var DBConfigData *DBConfig
var DB *sql.DB // ✅ this is the real sql.DB

func InitDB() {
	data, err := os.ReadFile("configs/db.yaml")
	if err != nil {
		panic("Failed to read db.yaml: " + err.Error())
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic("Failed to parse db.yaml: " + err.Error())
	}

	DBConfigData = &cfg.Database

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		DBConfigData.Host, DBConfigData.Port, DBConfigData.User,
		DBConfigData.Password, DBConfigData.DBName, DBConfigData.SSLMode)

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		panic("Failed to open DB connection: " + err.Error())
	}

	if err := DB.Ping(); err != nil {
		panic("Failed to ping DB: " + err.Error())
	}

	fmt.Println("✅ Connected to database")
}
