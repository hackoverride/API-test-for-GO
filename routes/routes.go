package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "api.reservation.oslofjord.com/database"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handlers struct {
	Client *mongo.Client
}

func (routes *Handlers) NewStay(w http.ResponseWriter, r *http.Request) {
	var filter, option interface{}

	/* Set header to expect a json file */
	w.Header().Set("Content-Type", "application/json")
	// filter  gets all document,

	filter = bson.D{primitive.E{Key: "rotation", Value: bson.D{primitive.E{Key: "$lt", Value: 70}}}}
	option = bson.D{primitive.E{Key: "_id", Value: 0}} // do not get the _id field.
	cursor, err := db.Query(routes.Client, r.Context(), "test",
		"comments", filter, option)

	if err != nil {
		// Defaulted error handeling
		panic(err)
	}

	/* bson.M = object with keys and values (bson.D is splitted object with {key: '', value: ''}, {key: '', value: ''}) */
	var results []bson.M
	if err := cursor.All(r.Context(), &results); err != nil {
		panic(err)
	}

	for _, doc := range results {
		fmt.Println(doc)
	}

	/* Returns the json file */
	json.NewEncoder(w).Encode(results)
}

func (routes *Handlers) GetBookingById(w http.ResponseWriter, r *http.Request) {
	/* Route variable = id */
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprintf(w, "You've requested the booking with id: %s", id)
}
