// [START gae_go111_app]

// Sample helloworld is an App Engine app.
package main

// [START import].
import (
	"context"
	"database/sql"
	"fmt"
	"helloworld/ent"
	"helloworld/ent/migrate"
	"helloworld/env"
	"helloworld/requesthandlers"
	"log"
	"net/http"
	"os"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

// [END import]
// [START main_func]

func initSocketConnectionPool() (*sql.DB, error) {
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

func SmsEndpoint() string {
	return "/zHwTHvytNXpesVFJwhtDuaFKyeFwLNaA"
}

func main() {
	var env = env.MustGetenv("EXPOSURES_ENV")
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Fatalf("sql.Open: %v", err)
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
		log.Fatalf("Could not open Postgres: %v", err)
	}
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true)); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	users := client.User.Query().AllX(context.Background())
	checkins := client.Checkin.Query().AllX(context.Background())
	log.Printf("Started exposures with these users: %s and these checkins: %s", users, checkins)

	http.HandleFunc("/", requesthandlers.IndexHandler)
	http.HandleFunc(SmsEndpoint(), requesthandlers.SmsHandler(context.Background(), client))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
	// [END setting_port]
}

// [END main_func]

// [START Handle Positive Notification].
