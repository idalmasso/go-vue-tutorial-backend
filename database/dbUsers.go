package database

import (
	"context"
	"fmt"

	commonLib "github.com/idalmasso/go_vue_tutorial_backend/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//FindUser finds a user from the usercollection
func FindUser(c context.Context, username string) (commonLib.UserDB,error){
	var result commonLib.UserDB

	err := UserCollection.FindOne(c, bson.M{"username": username}).Decode(&result)
	if err != nil {
			return result, err
	}
	return result, nil
}

//AddUser add a user to the usercollection
func AddUser(c context.Context, user commonLib.UserDB) (commonLib.UserDB,error){
	if result, err := UserCollection.InsertOne(c, user); err!=nil{
		return user, err
	}else{
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
				user.ID = oid
				return user, nil
			} else {
				return user, fmt.Errorf("Cannot get id from results")
			}
	}
}
//EditUserDescription edits a single user description (From patch)
func EditUserDescription(c context.Context, user commonLib.UserDB) (commonLib.UserDB,error){
	var updatedUser commonLib.UserDB
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	
	result := UserCollection.FindOneAndUpdate(c, bson.M{"username": user.Username}, bson.M{"$set": bson.M{"description": user.Description}} ,&opt);
	if  result.Err()!=nil{
		return user,  result.Err()
	}
	if err := result.Decode(&updatedUser); err != nil {
		return user, fmt.Errorf("Cannot decode: %w", err)
	}	
	return updatedUser, nil
	
}

//AddAuthenticationToken adds the authenticationToken to the user
func AddAuthenticationToken(c context.Context, user commonLib.UserDB) (commonLib.UserDB, error){
		var updatedUser commonLib.UserDB
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	
	result := UserCollection.FindOneAndUpdate(c, bson.M{"username": user.Username}, bson.M{"$set": bson.M{"AuthenticationToken": user.AuthenticationToken}} ,&opt);
	if  result.Err()!=nil{
		return user,  result.Err()
	}
	if err := result.Decode(&updatedUser); err != nil {
		return user, fmt.Errorf("Cannot decode: %w", err)
	}	
	return updatedUser, nil
}
