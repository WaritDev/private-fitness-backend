package main

import (
	"flag"
	"log"
	"os"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/migrations"
)

func main() { initMigrationCommand() }

func initMigrationCommand() {
	migrateSchema := flag.Bool("migrate:schema", false, "Run database schema migration (create database/schema if needed)")
	makeMigrationFile := flag.Bool("migrate:make", false, "Create new migration file")
	migrationFilename := flag.String("name", "", "Migration file name")
	migrateUp := flag.Bool("migrate:up", false, "Apply all up migrations")
	migrateDown := flag.Bool("migrate:down", false, "Apply down migrations by step")
	migrationDownStep := flag.Int("step", 0, "Down migration step count")
	migrationReset := flag.Bool("migrate:reset", false, "Reset database (down all then up all)")
	flag.Parse()

	switch {
	case *migrateSchema:
		log.Println("Migrating schema ...")
		migrations.MigrateSchema()
	case *makeMigrationFile:
		if *migrationFilename == "" {
			log.Fatalln("🚨 Missing name. Use: go run cmd/migration/main.go -migrate:make -name=<snake_case>")
		}
		log.Println("Creating new migration file ...")
		migrations.MakeMigration(migrationFilename)
	case *migrateUp:
		log.Println("Migrating up ...")
		migrations.MigrateUp()
	case *migrateDown:
		if *migrationDownStep == 0 {
			log.Fatalln("🚨 Missing step. Use: go run cmd/migration/main.go -migrate:down -step=<N>")
		}
		log.Println("Migrating down ...")
		migrations.MigrateDown(*migrationDownStep)
	case *migrationReset:
		log.Println("Resetting migration ...")
		migrations.MigrateReset()
	default:
		log.Println("No migration command specified.")
	}
	os.Exit(0)
}