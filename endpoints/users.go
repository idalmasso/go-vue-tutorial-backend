package endpoints

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	commonLib "github.com/idalmasso/go_vue_tutorial_backend/common"
	mdb "github.com/idalmasso/go_vue_tutorial_backend/database"
)

//getUser will return single user
func getUser(w http.ResponseWriter, r *http.Request) {
	log.Println("getuser called")
	vars := mux.Vars(r)
	username, ok := vars["USERNAME"]
	if !ok {
		http.Error(w, "Cannot find username in request", http.StatusBadRequest)
		return
	}
	if user, err := mdb.FindUser(r.Context(), username); err!=nil {
		http.Error(w, "Cannot find user", http.StatusNotFound) 
	} else {
		var userReturn commonLib.UserAPI
		userReturn.Username=user.Username
		userReturn.Description = user.Description
		sendJSONResponse(w, userReturn)
	}

}
