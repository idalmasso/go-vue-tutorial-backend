package endpoints

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	commonLib "github.com/idalmasso/go_vue_tutorial_backend/common"
	mdb "github.com/idalmasso/go_vue_tutorial_backend/database"
	"golang.org/x/crypto/bcrypt"
)

//getTokenUserPassword returns a jwt token for a user if the password is ok
func getTokenUserPassword(w http.ResponseWriter, r *http.Request) {
	var u commonLib.UserPassword
	err:=json.NewDecoder(r.Body).Decode(&u)
	if err!=nil{
		http.Error(w, "cannot decode username/password struct", http.StatusBadRequest)
		return
	}
	//here I have a user!
	//Now check if exists 
	 user, err:= mdb.FindUser(r.Context(), u.Username);
	 if err!=nil{
		http.Error(w, "Cannot find the username", http.StatusNotFound)
		return
	}
	
	err=bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(u.Password))
	if err!=nil{
		return
	}
	token, err:=createToken(u.Username)
	if err!=nil{
		http.Error(w, "Cannot create token", http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, struct {Token string `json:"token"`}{ token })

	
	
}

func createUser(w http.ResponseWriter, r *http.Request){
	log.Println("createUser called")
	var user commonLib.UserPassword
	var u commonLib.UserDB
	err := json.NewDecoder(r.Body).Decode(&user)
	if err!=nil{
		http.Error(w, "Cannot decode request", http.StatusBadRequest)
		return
	}
	//here I have a user!
	//Now check if exists 
	 _, err = mdb.FindUser(r.Context(), user.Username);
	 if err==nil{
		http.Error(w,"User already exists", http.StatusBadRequest)
		return
	}
	u.Username=user.Username
	//fOr now empty desc... Maybe another time we'll add some way to update the user
	u.Description = ""
	//If I'm here-> add user and return a token
	if value, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err!=nil{
		http.Error(w, "Cannot generate password hash: "+err.Error(), http.StatusInternalServerError)
		return
	}else{
		u.PasswordHash=value
		
	}
	if _, err= mdb.AddUser(r.Context(), u); err!=nil{
		http.Error(w, "Cannot insert user in database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	token, err:=createToken(u.Username)
	if err!=nil{
		http.Error(w, "Cannot create token", http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, struct {Token string `json:"token"`}{ token })
}

func getTokenByToken(w http.ResponseWriter, r *http.Request){
	//Here I already have the token checked... Just get the username from Request context
	username, ok :=context.Get(r,"username").(string)
	if !ok{
		http.Error(w, "Cannot check username", http.StatusInternalServerError)
		return
	}
	token, err:=createToken(username)
	if err!=nil{
		http.Error(w, "Cannot create token", http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, struct {Token string}{ token })
}

func createToken(username string) (string, error) {
  var err error
  //Creating Access Token
  
  atClaims := jwt.MapClaims{}
  atClaims["authorized"] = true
  atClaims["username"] = username
  atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	secret:= getSecret()
  token, err := at.SignedString([]byte(secret))
  if err != nil {
     return "", err
  }
  return token, nil
}


