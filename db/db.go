package db

import (
	"context"
	"database/sql"
	"exposures/ent"
	"exposures/ent/migrate"
	"exposures/env"
	"fmt"
	"log"
	"os"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func buildSqlDb() (*sql.DB, error) {
	var (
		dbUser = env.MustGetenv("DB_USER") // e.g. 'my-db-user'
		dbPwd  = env.MustGetenv("DB_PASS") // e.g. 'my-db-password'
		dbName = env.MustGetenv("DB_NAME") // e.g. 'my-database'
		dbHost = env.MustGetenv("DB_HOST")
		dbPort = env.MustGetenv("DB_PORT")
		env    = env.MustGetenv("EXPOSURES_ENV")
	)

	var dbURI string
	if env == "TEST" {
		file, err := os.Create("sqlite-database.db")
		if err != nil {
			log.Fatalf("Failed creating db: %v", err.Error())
		}
		file.Close()
		log.Println("Created SQLite Database for test.")
		sqldb, err := sql.Open("sqlite3", "./sqlite-database.db?_foreign_keys=ON")
		if err != nil {
			log.Fatalf("Failed initiating db connection: %v", err.Error())
		}
		return sqldb, err
	}
	dbURI = fmt.Sprintf(
		"user=%s password=%s database=%s host=%s port=%s sslmode=require",
		dbUser,
		dbPwd,
		dbName,
		dbHost,
		dbPort,
	)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("pgx", dbURI)

	return dbPool, err
}

func BuildEntClient() (*ent.Client, error) {
	db, err := buildSqlDb()
	var env = env.MustGetenv("EXPOSURES_ENV")

	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	dl := dialect.Postgres
	if env == "TEST" {
		dl = dialect.SQLite
	}
	drv := entsql.OpenDB(dl, db)

	client := ent.NewClient(ent.Driver(drv))
	if err != nil {
		return nil, err
	}

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true)); err != nil {
		return nil, err
	}
	return client, nil
}
