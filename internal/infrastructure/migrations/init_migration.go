package migrations

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/WaritDev/private-fitness-backend/config"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Migration struct {
	Title string
	Up    func(*pgxpool.Pool) error
	Down  func(*pgxpool.Pool) error
}

var Migrations []*Migration

func MigrateSchema() {
	ctx := context.Background()
	cfg := config.ProvideConfig()

	pgPool := db.GetPgPool()
	if pgPool == nil {
		panic("pgPool is nil")
	}

	_, err := pgPool.Exec(
		ctx,
		fmt.Sprintf(
			`CREATE SCHEMA IF NOT EXISTS %s;`, cfg.Schema,
		),
	)
	if err != nil {
		panic(err)
	}

	log.Printf("🏗️ Schema %s migrated successfully", cfg.Schema)
}

func MakeMigration(migrationFilename *string) {
	formattedFilename := time.Now().Format("20060102150405") + "_" + *migrationFilename + ".go"
	filepath := "internal/infrastructure/migrations/" + formattedFilename

	migrationVarName := strcase.ToLowerCamel(strings.ReplaceAll(*migrationFilename, "_", " "))

	formattedTemplate := strings.Replace(migrationTemplate, "<migration_name>", migrationVarName, -1)
	formattedTemplate = strings.Replace(formattedTemplate, "<filename>", formattedFilename, -1)

	_, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filepath, []byte(formattedTemplate), 0)
	if err != nil {
		panic(err)
	}

	log.Printf("📁 Migration file %s created successfully", formattedFilename)
}

const migrationTemplate = `
package migrations
import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)
func init() {
	Migrations = append(Migrations, <migration_name>)
}
var <migration_name> = &Migration{
	Title: "<filename>",
	Up: func(pgPool *pgxpool.Pool) error {
		_, err := pgPool.Exec(context.Background(),` + "`" + `
		` + "`" + `)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(pgPool *pgxpool.Pool) error {
		_, err := pgPool.Exec(context.Background(), ` + "`" + `
		` + "`" + `)
		if err != nil {
			return err
		}
		return nil
	},
}
`

func MigrateUp() {
	ctx := context.Background()

	pgPool := db.GetPgPool()
	if pgPool == nil {
		panic("pgPool is nil")
	}

	_, err := pgPool.Exec(
		ctx,
		`CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (title)
		);`,
	)
	if err != nil {
		if err.Error() == "ERROR: no schema has been selected to create in (SQLSTATE 3F000)" {
			log.Fatalln("🚨 Schema not found. Please run `go run main.go -migrate:schema` to create schema")
		}
		panic(err)
	}

	migrations, err := pgPool.Query(
		ctx,
		`SELECT title FROM migrations ORDER BY title;`,
	)
	if err != nil {
		panic(err)
	}
	defer migrations.Close()

	// executedMigrations := make(map[string]bool)
	executedMigrations := make(map[string]interface{})
	for migrations.Next() {
		var title string
		if err := migrations.Scan(&title); err != nil {
			panic(err)
		}
		executedMigrations[title] = nil
	}

	for _, migration := range Migrations {
		_, exists := executedMigrations[migration.Title]
		if !exists {
			log.Printf("🚀 Migrating up %s ...", migration.Title)
			err := migration.Up(pgPool)
			if err != nil {
				panic(err)
			}

			args := pgx.NamedArgs{
				"title": migration.Title,
			}
			_, err = pgPool.Exec(
				ctx,
				`INSERT INTO migrations (title) VALUES (@title);`,
				args,
			)
			if err != nil {
				panic(err)
			}
		}
	}
}

func MigrateDown(step int) {
	ctx := context.Background()

	pgPool := db.GetPgPool()

	currentMigrations, err := pgPool.Query(
		ctx,
		`SELECT title FROM migrations ORDER BY title DESC;`,
	)
	if err != nil {
		panic(err)
	}
	defer currentMigrations.Close()

	executedMigrations := []string{}
	for currentMigrations.Next() {
		var title string
		if err := currentMigrations.Scan(&title); err != nil {
			panic(err)
		}
		executedMigrations = append(executedMigrations, title)
	}

	if step > len(executedMigrations) {
		step = len(executedMigrations)
	}

	for i := 0; i < step; i++ {
		for _, migration := range Migrations {
			if migration.Title == executedMigrations[i] {
				log.Printf("🔙 Migrating down %s ...", migration.Title)
				err := migration.Down(pgPool)
				if err != nil {
					panic(err)
				}

				args := pgx.NamedArgs{
					"title": migration.Title,
				}
				_, err = pgPool.Exec(
					ctx,
					`DELETE FROM migrations WHERE title = @title;`,
					args,
				)
				if err != nil {
					panic(err)
				}
				break
			}
		}
	}
}

func MigrateReset() {
	ctx := context.Background()

	pgPool := db.GetPgPool()
	if pgPool == nil {
		panic("pgPool is nil")
	}

	_, err := pgPool.Exec(
		ctx,
		fmt.Sprintf(
			`DROP SCHEMA IF EXISTS %s CASCADE;`,
			config.ProvideConfig().Schema,
		),
	)
	if err != nil {
		panic(err)
	}

	log.Printf("🔥 Schema %s reset successfully", config.ProvideConfig().Schema)
}