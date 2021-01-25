package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//PostCollection is the actual collection for posts
var PostCollection *mongo.Collection
//UserCollection is the actual collection for users
var UserCollection *mongo.Collection

//SetCollections sets the db and correct collection
func SetCollections(c *mongo.Database) {
	PostCollection = c.Collection("posts")
	UserCollection = c.Collection("users")
}

//Connect connects to the DB and also set up the controllers
func Connect() {
	// Database Config
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)

	//Set up a context required by mongo.Connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	//To close the connection at the end
	defer cancel()

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}
	db := client.Database("postsDB")
	SetCollections(db)
	return
}


