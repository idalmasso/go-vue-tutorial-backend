package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mdb "github.com/idalmasso/go_vue_tutorial_backend/database"
	"github.com/idalmasso/go_vue_tutorial_backend/endpoints"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main(){
	err := godotenv.Load(".env")

  if err != nil {
    log.Fatalf("Error loading .env file")
  }

	mdb.Connect()
	r := mux.NewRouter()
	r=endpoints.AddRouterEndpoints(r)
	fs := http.FileServer(http.Dir("./dist"))
	r.PathPrefix("/").Handler(fs)
	
	http.Handle("/",&corsRouterDecorator{r})
	fmt.Println("Listening")	
	log.Panic(
		http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), nil),
	)
}


type corsRouterDecorator struct {
	R *mux.Router
}

func (c *corsRouterDecorator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, PATCH")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	}
		// Stop here if its Preflighted OPTIONS request, I just add an OK 
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}


	c.R.ServeHTTP(w, r)
}
