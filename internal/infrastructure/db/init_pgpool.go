package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/WaritDev/private-fitness-backend/config"
)

func ProvideMariaDB(ctx context.Context, cfg *config.Config) *sql.DB {
	port := cfg.DBPort
	if (cfg.UseDocker && port == "") || port == "" {
		port = "3306"
	}

	params := cfg.DBParams
	if params == "" {
		params = "parseTime=true&charset=utf8mb4&loc=Asia%2FBangkok"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, port, cfg.DBName, params)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ open mysql dsn error: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("❌ MariaDB not reachable: %v", err)
	}

	log.Println("🫙 Connected to MariaDB")
	return db
}

func GetDB() *sql.DB {
	ctx := context.Background()
	cfg := config.ProvideConfig()
	return ProvideMariaDB(ctx, cfg)
}