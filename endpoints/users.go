package endpoints

import (
	"encoding/json"
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
		sendJSONResponse(w, userReturn, http.StatusOK)
	}

}

func editUserDescription(w http.ResponseWriter, r *http.Request) {
	log.Println("editUserDescription called")
	vars := mux.Vars(r)
	username, ok := vars["USERNAME"]
		if !ok {
		http.Error(w, "Cannot find username in request", http.StatusBadRequest)
		return
	}
	var userAPI commonLib.UserAPI
	err:=json.NewDecoder(r.Body).Decode(&userAPI)
	if err!=nil{
		http.Error(w, "Not a valid json request", http.StatusBadRequest)
		return
	}
	if userAPI.Username!=username || !isUsernameContextOk(username, r){
		http.Error(w, "Cannot update a different user", http.StatusBadRequest)
		return
	}
	var user commonLib.UserDB
	user.Username = username
	user.Description = userAPI.Description
	if user, err := mdb.EditUserDescription(r.Context(), user); err!=nil {
		http.Error(w, "Cannot find user", http.StatusNotFound) 
	} else {
		userAPI.Description = user.Description
		sendJSONResponse(w, userAPI, http.StatusOK)
	}
}
