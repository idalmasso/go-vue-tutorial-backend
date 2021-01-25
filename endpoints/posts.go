package endpoints

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	commonLib "github.com/idalmasso/go_vue_tutorial_backend/common"
	mdb "github.com/idalmasso/go_vue_tutorial_backend/database"
)

//addPost will get in Body a post with ONLY username and post-> need to add the others and save it
func addPost(w http.ResponseWriter, r *http.Request) {
	log.Println("addPost called")
	var actualPost commonLib.Post
	err:=json.NewDecoder(r.Body).Decode(&actualPost)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isUsernameContextOk(actualPost.Username, r){
		http.Error(w, "Cannot post for another user", http.StatusUnauthorized)
		return 
	}
	if actualPost, err = mdb.AddSinglePost(r.Context(), actualPost); err!=nil{
		//Care! THis is not the actual "nice" way to do this... I should create some new Error types 
		http.Error(w, "Error inserting data", http.StatusInternalServerError)
	}else{
		sendJSONResponse(w,actualPost)
	}
}
//deletePost removes the post that is being passed. Get the id from the query
func deletePost(w http.ResponseWriter, r *http.Request) {
	log.Println("deletePost called")
	vars := mux.Vars(r)
	idString, ok := vars["POST_ID"]
	if !ok {
		http.Error(w, "Cannot find ID", http.StatusBadRequest)
		return
	}
	if err:=mdb.DeletePost(r.Context(), idString); err!=nil{
		if strings.HasPrefix(err.Error(), "Not found") {
			http.Error(w,"Cannot find the id", http.StatusNotFound)
		}else{
			http.Error(w, "Internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
//addComment will get the comment in the body, and the id in the query
func addComment(w http.ResponseWriter, r *http.Request) {
	log.Println("addComment called")
	vars := mux.Vars(r)
	idString, ok := vars["POST_ID"]
	if !ok {
		http.Error(w, "Cannot find ID", http.StatusBadRequest)
		return
	}
	var actualComment commonLib.Comment
	err := json.NewDecoder(r.Body).Decode(&actualComment)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isUsernameContextOk(actualComment.Username, r){
		http.Error(w, "Cannot comment post for another user", http.StatusUnauthorized)
		return 
	}
	if post, err := mdb.AddComment(r.Context(), actualComment, idString); err!=nil{
		if strings.HasPrefix(err.Error(), "Not found") {
			http.Error(w,"Cannot find the id", http.StatusNotFound)
		}else{
			http.Error(w, "Internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}	else	{
		sendJSONResponse(w, post)
	}
}
//getPosts will return all the posts actually in the array
func getPosts(w http.ResponseWriter, r *http.Request) {
	log.Println("getPosts called")
	if posts, err:= mdb.GetAllPosts(r.Context());err!=nil{
		http.Error(w, "Cannot read from DB", http.StatusInternalServerError)
	}else{
		sendJSONResponse(w, posts)
	}
}

func isUsernameContextOk(username string, r *http.Request) bool {
	usernameCtx, ok:=context.Get(r, "username").(string)
	if !ok{
		return false
	}
	if usernameCtx!=username{
		return false
	}
	return true
}
