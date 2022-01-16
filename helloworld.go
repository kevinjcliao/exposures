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
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
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
		dbUser = mustGetenv("DB_USER") // e.g. 'my-db-user'
		dbPwd  = mustGetenv("DB_PASS") // e.g. 'my-db-password'
		dbName = mustGetenv("DB_NAME") // e.g. 'my-database'
		dbHost = mustGetenv("DB_HOST")
		dbPort = mustGetenv("DB_PORT")
		env    = mustGetenv("EXPOSURES_ENV")
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

func main() {
	var env = mustGetenv("EXPOSURES_ENV")
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
	sendSms(from, "Thanks for letting us know. Please take care of yourself and those around you. We hope you recover soon.")
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
		phoneNumbers := make(map[string]bool)
		for _, similarCheckin := range similarCheckins {
			sender := similarCheckin.QuerySender().OnlyX(ctx)
			phoneNumbers[sender.PhoneNumber] = true
		}
		for phone, _ := range phoneNumbers {
			date := time.Unix(positiveCheckin.CheckinTime, 0)
			body := fmt.Sprintf("Someone at an event you attended on: %s tested positive for COVID-19. You should get tested.", date)
			sendSms(phone, body)
		}
	}
}

// [END Handle Positive Notification].

// [START smsHandler].
func sendSms(to string, body string) {
	if mustGetenv("EXPOSURES_ENV") != "PROD" {
		log.Println("Sending an SMS to: ", to, " with body: ", body)
		log.Println("Not actually sending an SMS because you're in test.")
		return
	}
	client := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username: mustGetenv("TWILIO_SID"),
		Password: mustGetenv("TWILIO_API_KEY"),
	})
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom("+1(417)668-2737")
	params.SetBody(body)
	_, err := client.ApiV2010.CreateMessage(params)
	if err != nil {
		log.Fatalf("Failed to send sms to: %s containing body: %s", to, body)
	}
}

func smsHandler(ctx context.Context, entClient *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Failed to parse request.")
		}
		body := r.Form.Get("Body")
		from := r.Form.Get("From")

		if body == "POSITIVE" {
			handlePositiveCase(ctx, entClient, from)
			return
		}

		rsvpMessage := "Thanks for checking in. You can let your friends/other attendees know if you tested positive by replying POSITIVE."
		sendSms(from, rsvpMessage)
		user, err := entClient.User.Query().Where(user.PhoneNumberEQ(from)).Only(ctx)
		if err != nil {
			switch err.(type) {
			case *ent.NotFoundError:
				user = entClient.User.Create().SetPhoneNumber(from).SaveX(ctx)

			default:
				log.Println("Huh")
			}
		}
		_, err = entClient.Checkin.Create().
			SetCheckinTime(time.Now().Unix()).
			SetEventID(body).
			SetSender(user).
			Save(ctx)
		if err != nil {
			log.Printf("Error creating checkin: %v", err)
		}

		if err != nil {
			log.Println(err.Error())
		} else {
			log.Println("SMS Sent Successfully!")
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
