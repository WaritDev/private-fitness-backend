package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db"
	"github.com/iancoleman/strcase"
)

type Migration struct {
	Title string
	Up    func(*sql.DB) error
	Down  func(*sql.DB) error
}

var Migrations []*Migration

func MigrateSchema() {
	log.Printf("🏗️ MariaDB uses database '%s' already. No schema step required.", "<db from DSN>")
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
	if err = os.WriteFile(filepath, []byte(formattedTemplate), 0); err != nil {
		panic(err)
	}

	log.Printf("📁 Migration file %s created successfully", formattedFilename)
}

const migrationTemplate = "package migrations\n" +
"import (\n" +
"\t\"context\"\n" +
"\t\"database/sql\"\n" +
")\n" +
"\n" +
"func init() {\n" +
"\tMigrations = append(Migrations, <migration_name>)\n" +
"}\n" +
"\n" +
"var <migration_name> = &Migration{\n" +
"\tTitle: \"<filename>\",\n" +
"\tUp: func(db *sql.DB) error {\n" +
"\t\t_, err := db.ExecContext(context.Background(), " + "`" + "\n" +
"\t\t-- WRITE YOUR 'UP' SQL HERE (MariaDB)\n" +
"\t\t" + "`" + ")\n" +
"\t\tif err != nil {\n" +
"\t\t\treturn err\n" +
"\t\t}\n" +
"\t\treturn nil\n" +
"\t},\n" +
"\tDown: func(db *sql.DB) error {\n" +
"\t\t_, err := db.ExecContext(context.Background(), " + "`" + "\n" +
"\t\t-- WRITE YOUR 'DOWN' SQL HERE (MariaDB)\n" +
"\t\t" + "`" + ")\n" +
"\t\tif err != nil {\n" +
"\t\t\treturn err\n" +
"\t\t}\n" +
"\t\treturn nil\n" +
"\t},\n" +
"}\n"

func MigrateUp() {
	ctx := context.Background()
	sqldb := db.GetDB()
	if sqldb == nil {
		panic("db is nil")
	}

	_, err := sqldb.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS migrations (
			id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`)
	if err != nil {
		panic(err)
	}

	rows, err := sqldb.QueryContext(ctx, `SELECT title FROM migrations ORDER BY title;`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	executed := make(map[string]struct{})
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			panic(err)
		}
		executed[title] = struct{}{}
	}

	for _, m := range Migrations {
		if _, ok := executed[m.Title]; ok {
			continue
		}
		log.Printf("🚀 Migrating up %s ...", m.Title)
		if err := m.Up(sqldb); err != nil {
			panic(err)
		}
		if _, err := sqldb.ExecContext(ctx, `INSERT INTO migrations (title) VALUES (?)`, m.Title); err != nil {
			panic(err)
		}
	}
}

func MigrateDown(step int) {
	ctx := context.Background()
	sqldb := db.GetDB()
	if sqldb == nil {
		panic("db is nil")
	}

	rows, err := sqldb.QueryContext(ctx, `SELECT title FROM migrations ORDER BY title DESC;`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	executed := []string{}
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			panic(err)
		}
		executed = append(executed, title)
	}

	if step > len(executed) {
		step = len(executed)
	}

	for i := 0; i < step; i++ {
		target := executed[i]
		for _, m := range Migrations {
			if m.Title == target {
				log.Printf("🔙 Migrating down %s ...", m.Title)
				if err := m.Down(sqldb); err != nil {
					panic(err)
				}
				if _, err := sqldb.ExecContext(ctx, `DELETE FROM migrations WHERE title = ?`, m.Title); err != nil {
					panic(err)
				}
				break
			}
		}
	}
}

// ล้างทั้งฐาน (ดรอปทุกตารางใน DB ปัจจุบัน)
func MigrateReset() {
	ctx := context.Background()
	sqldb := db.GetDB()
	if sqldb == nil {
		panic("db is nil")
	}

	log.Printf("🔥 Reset current MariaDB database (drop all tables)")
	_, _ = sqldb.ExecContext(ctx, `SET FOREIGN_KEY_CHECKS=0;`)

	// ดึงรายการตารางทั้งหมด
	tables := []string{}
	rows, err := sqldb.QueryContext(ctx, `SHOW TABLES;`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			panic(err)
		}
		tables = append(tables, name)
	}

	for _, t := range tables {
		if t == "" {
			continue
		}
		if _, err := sqldb.ExecContext(ctx, fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", t)); err != nil {
			panic(err)
		}
	}

	_, _ = sqldb.ExecContext(ctx, `SET FOREIGN_KEY_CHECKS=1;`)
	log.Printf("✅ Reset completed.")
}