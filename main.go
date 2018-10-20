package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/Henry-Sarabia/craft"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

var (
	crafter = &craft.Crafter{}
)

func init() {
	var err error
	crafter, err = craft.NewFromFiles("templates.json", "classes.json", "materials.json", "qualities.json", "details.json")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.Use(handlers.RecoveryHandler())
	r.Use(handlers.CORS())

	api := r.PathPrefix("/api/v1/").Subrouter()
	api.HandleFunc("/item", itemHandler)

	os.Setenv("PORT", "8081")
	port, err := getPort()
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, r),
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	i, err := crafter.NewItem()
	if err != nil {
		http.Error(w, "cannot generate item", http.StatusInternalServerError)
		log.Println(err)
	}
	fmt.Println(i)
	fmt.Println(i.Description)
	render.JSON(w, r, i)
}

// getPort returns the port from the $PORT environment variable as a string.
// Returns an error if $PORT is not set.
func getPort() (string, error) {
	p := os.Getenv("PORT")
	if p == "" {
		return "", errors.New("$PORT must be set")
	}

	return p, nil
}
