package endpoints

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)
type comment struct {
	ID   		 int      `json:"id"`
	Username string   `json:"username"`
	Post string    		`json:"post"`
	Date time.Time 		`json:"date"`	
}
//post Struct is used as post structure...
type post struct {
	ID   		 int       `json:"id"`
	Username string    `json:"username"`
	Post string    		 `json:"post"`
	Date time.Time 		 `json:"date"`
	Comments []comment `json:"comments"`
}
//This is my "database", in memory... Will be changed in a real database in future...
var posts []post=make([]post, 0)
//need an index for the array... When I'll delete the posts the index will have to go on...
var index int=1

//addPost will get in Body a post with ONLY username and post-> need to add the others and save it
func addPost(w http.ResponseWriter, r *http.Request) {
	log.Println("addPost called")
	var actualPost post
	err:=json.NewDecoder(r.Body).Decode(&actualPost)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	actualPost.ID = index
	index++
	actualPost.Date=time.Now()
	if actualPost.Comments== nil{
		actualPost.Comments=make([]comment, 0)
	}
	posts=append(posts, actualPost)
	sendJSONResponse(w,actualPost)
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
	id, err := strconv.Atoi(idString)
	if err!=nil{
		http.Error(w, "Cannot convert the id value to string", http.StatusBadRequest)
		return
	}
	for i:=0;i<len(posts);i++{
		if posts[i].ID==id {
			posts[i]=posts[len(posts)-1]
			posts=posts[:len(posts)-1]
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	http.Error(w, "Cannot find the requested id", http.StatusNotFound)
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
	id, err := strconv.Atoi(idString)
	if err!=nil{
		http.Error(w, "Cannot convert the id value to string", http.StatusBadRequest)
		return
	}
	var actualComment comment
	err = json.NewDecoder(r.Body).Decode(&actualComment)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i:=0;i<len(posts);i++{
		if posts[i].ID==id {
			//Now I have the post
			var commMax int=0
			for comm:=0;comm<len(posts[i].Comments);comm++{
				if commMax<posts[i].Comments[comm].ID{
					commMax=posts[i].Comments[comm].ID
				}
			}
			actualComment.ID=commMax+1
			actualComment.Date=time.Now()
			posts[i].Comments = append(posts[i].Comments, actualComment)
			sendJSONResponse(w, posts[i])
			return
		}
	}
	//If I'm here, there is no post with the id searched... 
	http.Error(w, "Cannot find a post with the selected id", http.StatusNotFound)
}
//getPosts will return all the posts actually in the array
func getPosts(w http.ResponseWriter, r *http.Request) {
	log.Println("Get post called")
	sendJSONResponse(w, posts)
}
