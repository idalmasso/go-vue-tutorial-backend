package commmon

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//UserDB struct will contain the info about authentication. Password won't be saved, only password hash.
type UserDB struct{
	ID		primitive.ObjectID  `bson:"_id,omitempty"`
	Username string  
	PasswordHash []byte 
	Description string 
}
//UserAPI is a struct used to be passed back to the user (so, no username and password hash needed here!)
type UserAPI struct{
	Username string `json:"username"`
	Description string `json:"description"`
}
//UserPassword type that will be used to decode the user login & signup requests. No Password will be saved, in any way!
type UserPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
//Comment is the actual data of Comment type
type Comment struct {
	ID   		 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string   `json:"username"`
	Post string    		`json:"post"`
	Date time.Time 		`json:"date"`	
}

//Post Struct is used as post structure...
type Post struct {
	ID   		 primitive.ObjectID       `json:"id" bson:"_id,omitempty"`
	Username string    `json:"username"`
	Post string    		 `json:"post"`
	Date time.Time 		 `json:"date"`
	Comments []Comment `json:"comments"`
}


