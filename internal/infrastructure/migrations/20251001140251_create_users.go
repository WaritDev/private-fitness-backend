
package migrations
import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)
func init() {
	Migrations = append(Migrations, createUsers)
}
var createUsers = &Migration{
	Title: "20251001140251_create_users.go",
	Up: func(pgPool *pgxpool.Pool) error {
		_, err := pgPool.Exec(context.Background(),`
			CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(100) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT now()
		);
		`)
		if err != nil {
			return err
		}
		return nil
	},
	Down: func(pgPool *pgxpool.Pool) error {
		_, err := pgPool.Exec(context.Background(), `
		`)
		if err != nil {
			return err
		}
		return nil
	},
}
