package endpoints

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//getUser will return single user
func getUser(w http.ResponseWriter, r *http.Request) {
	log.Println("getuser called")
	vars := mux.Vars(r)
	user, ok := vars["USERNAME"]
	if !ok {
		http.Error(w, "Cannot find username in request", http.StatusBadRequest)
		return
	}
	if _, ok :=users[user]; ok{
		sendJSONResponse(w, struct{Username string `json:"username"`; Description string `json:"description"` }{user ,  ""})
		return
	}
	http.Error(w, "Cannot find user", http.StatusNotFound)
}
