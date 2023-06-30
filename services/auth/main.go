package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/scraparty/scraparty/auth/routes"
	scrapartydb "github.com/scraparty/scraparty-db"
)

func main() {
	db, err := scrapartydb.Connect()

	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	r.Get("/auth/github", func(w http.ResponseWriter, r *http.Request) {
		err := routes.Login(w, r, db)

		if err != nil {
			fmt.Println(err)

			http.Error(w, "There was an error authenticating you. Please try again later.", http.StatusInternalServerError)
		}
	})

	r.Get("/auth/github/callback", func(w http.ResponseWriter, r *http.Request) {
		err := routes.Callback(w, r, db)

		if err != nil {
			fmt.Println(err)

			http.Error(w, "There was an error authenticating you. Please try again later.", http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":3000", r)
}
