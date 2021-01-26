package database

import (
	"context"
	"log"
	"os"
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
	PostCollection = c.Collection(os.Getenv("POSTS_COLLECTION_NAME"))
	UserCollection = c.Collection(os.Getenv("USER_COLLECTION_NAME"))
}

//Connect connects to the DB and also set up the controllers
func Connect() {
	// Database Config
	clientOptions := options.Client().ApplyURI(os.Getenv("CONNECTION_STRING"))
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
	db := client.Database(os.Getenv("POSTS_DB_NAME"))
	SetCollections(db)
	return
}


