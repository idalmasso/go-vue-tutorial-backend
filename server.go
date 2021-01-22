package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/idalmasso/go_vue_tutorial_backend/endpoints"

	"github.com/gorilla/mux"
)

func main(){
	r := mux.NewRouter()
	r=endpoints.AddRouterEndpoints(r)
	fs := http.FileServer(http.Dir("./dist"))
	r.PathPrefix("/").Handler(fs)
	
	http.Handle("/",r)
	fmt.Println("Listening")	
	log.Panic(
		http.ListenAndServe(":3000", nil),
	)
}
