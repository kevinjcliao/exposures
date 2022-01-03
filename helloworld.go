// [START gae_go111_app]

// Sample helloworld is an App Engine app.
package main

// [START import].
import (
	"context"
	"database/sql"
	"fmt"
	"helloworld/ent"
	"helloworld/ent/checkin"
	"helloworld/ent/migrate"
	"helloworld/ent/user"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// [END import]
// [START main_func]

// mustGetEnv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

func initSocketConnectionPool() (*sql.DB, error) {
	var (
		dbUser                 = mustGetenv("DB_USER") // e.g. 'my-db-user'
		dbPwd                  = mustGetenv("DB_PASS") // e.g. 'my-db-password'
		instanceConnectionName = mustGetenv(
			"INSTANCE_CONNECTION_NAME",
		) // e.g. 'project:region:instance'
		dbName    = mustGetenv("DB_NAME") // e.g. 'my-database'
		env       = mustGetenv("EXPOSURES_ENV")
		socketDir = "/cloudsql"
	)

	var dbURI string
	if env == "TEST" {
		dbTCPHost := "34.134.62.245"
		dbPort := "5432"
		dbURI = fmt.Sprintf(
			"host=%s user=%s password=%s port=%s database=%s",
			dbTCPHost,
			dbUser,
			dbPwd,
			dbPort,
			dbName,
		)
	} else {
		dbURI = fmt.Sprintf(
			"user=%s password=%s database=%s host=%s/%s",
			dbUser,
			dbPwd,
			dbName,
			socketDir,
			instanceConnectionName,
		)
	}
	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("pgx", dbURI)

	return dbPool, err
}

func main() {
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Fatalf("sql.Open: %v", err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	drv := entsql.OpenDB(dialect.Postgres, db)

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

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/zHwTHvytNXpesVFJwhtDuaFKyeFwLNaA", smsHandler(context.Background(), client))
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
func handlePositiveCase(ctx context.Context, client *ent.Client, from string) {
	fourteenDaysAgo := time.Now().Unix() - (14 * 24 * 60 * 60)
	threeHours := int64(3 * 60 * 60)
	checkins := client.Checkin.Query().
		Where(
			checkin.And(
				checkin.HasSenderWith(user.PhoneNumberEQ(from)),
				checkin.CheckinTimeGT(fourteenDaysAgo),
			),
		).
		AllX(ctx)
	for _, positiveCheckin := range checkins {
		similarCheckins := client.Checkin.Query().Where(
			checkin.EventIDEQ(positiveCheckin.EventID),
		).
			Where(
				checkin.HasSenderWith(user.PhoneNumberNEQ(from)),
			).
			Where(
				checkin.Or(
					checkin.CheckinTimeGT(positiveCheckin.CheckinTime-threeHours),
					checkin.CheckinTimeLT(positiveCheckin.CheckinTime+threeHours),
				),
			).
			AllX(ctx)
		fmt.Println("Similar checkins: ", similarCheckins)
	}
}

// [END Handle Positive Notification].

// [START smsHandler].
func smsHandler(ctx context.Context, entClient *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Failed to parse request.")
		}
		body := r.Form.Get("Body")
		from := r.Form.Get("From")
		client := twilio.NewRestClientWithParams(twilio.RestClientParams{
			Username: mustGetenv("TWILIO_SID"),
			Password: mustGetenv("TWILIO_API_KEY"),
		})

		if body == "POSITIVE" {
			handlePositiveCase(ctx, entClient, from)
			return
		}
		params := &openapi.CreateMessageParams{}
		params.SetTo(from)
		params.SetFrom("+1(417)668-2737")
		rsvpMessage := "Thanks for checking in. You can let your friends/other attendees know if you tested positive by replying POSITIVE."
		params.SetBody(rsvpMessage)
		user, err := entClient.User.Query().Where(user.PhoneNumberEQ(from)).Only(ctx)
		if err != nil {
			switch err.(type) {
			case *ent.NotFoundError:
				user = entClient.User.Create().SetPhoneNumber(from).SaveX(ctx)

			default:
				fmt.Println("Huh")
			}
		}
		_, err = entClient.Checkin.Create().
			SetCheckinTime(time.Now().Unix()).
			SetEventID(body).
			SetSender(user).
			Save(ctx)
		if err != nil {
			fmt.Printf("Error creating checkin: %v", err)
		}

		_, err = client.ApiV2010.CreateMessage(params)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("SMS Sent Successfully!")
		}
	}
}

// [END smsHandler]

// [START indexHandler]

// indexHandler responds to requests with our greeting.
type IndexData struct {
	Uuid string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	t := template.Must(template.ParseFiles("templates/index.html"))
	err := t.Execute(w, IndexData{Uuid: uuid.New().String()})
	if err != nil {
		log.Println("Failed to parse index template.")
	}
}

// [END indexHandler]
// [END gae_go111_app]
