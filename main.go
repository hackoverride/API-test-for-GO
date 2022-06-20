package main

import (
	"log"
	"net/http"
	"os"
	"time"

	db "api.reservation.oslofjord.com/database"
	routes "api.reservation.oslofjord.com/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {

	/* ## START - MongoConnection - START ## */
	mongoConnection := goDotEnvVariable("MONGO_STRING")
	client, ctx, cancel, err := db.Connect(mongoConnection)
	if err != nil {
		panic(err)
	}

	defer db.Close(client, ctx, cancel)

	db.Ping(client, ctx)
	/* ## END - MongoConnection - END ## */

	/* ## START - BookingRoutes - START ## */
	r := mux.NewRouter()
	h := routes.Handlers{
		Client: client,
	}

	r.HandleFunc("/booking/{id}", h.GetBookingById)
	/* # Subroutes : /tests */
	t := r.PathPrefix("/tests").Subrouter()
	t.HandleFunc("/", h.NewStay).Methods("GET")
	/* ## END - BookingRoutes - END ## */

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
