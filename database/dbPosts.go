package database

import (
	"context"
	"fmt"
	"log"
	"time"

	commonLib "github.com/idalmasso/go_vue_tutorial_backend/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//GetAllPosts returns all posts made by everyone (not a real world scenario...)
func GetAllPosts(c context.Context) ([]commonLib.Post,error){
	var posts [] commonLib.Post
	cursor, err:=PostCollection.Find(c, bson.M{})
	if err != nil {
		log.Printf("Error while getting cursor, Reason: %v\n", err)
		return nil, err
	}
	
	for cursor.Next(c) {
		var post commonLib.Post
		cursor.Decode(&post)
		posts = append(posts, post)
	}
	return posts, nil
}	
//AddSinglePost adds a single post to the db
func AddSinglePost(c context.Context, post commonLib.Post) (commonLib.Post, error) {
	post.Date=time.Now()
	if post.Comments== nil{
		post.Comments=make([]commonLib.Comment, 0)
	}
	if result, err := PostCollection.InsertOne(c, post); err!=nil{
		return post, err
	}else{
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
				post.ID = oid
				return post, nil
			} else {
				return post, fmt.Errorf("Cannot get id from results")
			}
	}
}
//AddComment adds a single comment to the post with id postID to the  db. Returns the newly post and also the error code to be used
func AddComment(c context.Context,comment commonLib.Comment, postID string  ) (commonLib.Post,  error) {
	comment.Date = time.Now()
	var post commonLib.Post
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	objID, err := primitive.ObjectIDFromHex(postID) 
	if err != nil {
		return post, fmt.Errorf("Cannot decode postID")
	}
	comment.ID = primitive.NewObjectID()
	pushValue := bson.M{"$push": bson.M{"comments": comment}}
	result := PostCollection.FindOneAndUpdate(c, bson.M{"_id": objID}, pushValue, &opt)
	if result.Err() != nil {
		return post, fmt.Errorf("Cannot find and update post, reason: ", result.Err().Error())
	}
	if err := result.Decode(&post); err != nil {
		return post, fmt.Errorf("Not found post, reason: ", err)
	}
	return post, nil
}

//DeletePost deletes a single post with id postID from the db
func DeletePost(c context.Context, postID string ) (error) {
	objID, err := primitive.ObjectIDFromHex(postID) 
	if err != nil {
		return fmt.Errorf("Cannot decode postID")
	}
	result, err := PostCollection.DeleteOne(c, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("Error, cannot delete ID in request" + err.Error())
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("Not found post")
	}
	return nil
}
