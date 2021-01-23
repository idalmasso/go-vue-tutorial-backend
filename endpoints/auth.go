package endpoints

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//User struct will contain the info about authentication. Password won't be saved!
type User struct{
	Username string `json:"username"`
	Password string `json:"password"`
}

var users map[string][]byte = make(map[string][]byte)
var idxUsers int =0

//getTokenUserPassword returns a jwt token for a user if the password is ok
func getTokenUserPassword(w http.ResponseWriter, r *http.Request) {
	var u User
	err:=json.NewDecoder(r.Body).Decode(&u)
	if err!=nil{
		http.Error(w, "cannot decode username/password struct", http.StatusBadRequest)
		return
	}
	//here I have a user!
	//Now check if exists 
	passwordHash, found:= users[u.Username]
	if !found{
		http.Error(w, "Cannot find the username", http.StatusNotFound)
	}
	
	err=bcrypt.CompareHashAndPassword(passwordHash, []byte(u.Password))
	if err!=nil{
		http.Error(w, "Invalid password!", http.StatusForbidden)
		return
	}
	token, err:=createToken(u.Username)
	if err!=nil{
		http.Error(w, "Cannot create token", http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, struct {Token string}{ token })

	
	
}

func createUser(w http.ResponseWriter, r *http.Request){
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err!=nil{
		http.Error(w, "Cannot decode request", http.StatusBadRequest)
		return
	}
	if _, found:= users[u.Username]; found{
		http.Error(w,"User already exists", http.StatusBadRequest)
		return
	}
	//If I'm here-> add user and return a token
	value, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	users[u.Username]=value
	token, err:=createToken(u.Username)
	if err!=nil{
		http.Error(w, "Cannot create token", http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, struct {Token string}{ token })
}


func createToken(userid string) (string, error) {
  var err error
  //Creating Access Token
  
  atClaims := jwt.MapClaims{}
  atClaims["authorized"] = true
  atClaims["user_id"] = userid
  atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	secret:=os.Getenv("ACCESS_SECRET")
	if secret==""{
		secret="sdmalncnjsdsmf"
	}
  token, err := at.SignedString([]byte(secret))
  if err != nil {
     return "", err
  }
  return token, nil
}
