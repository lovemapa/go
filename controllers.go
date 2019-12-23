package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type user struct {
	Name string `json:name`
	Age  int    `json:age`
	City string `json:city`
}

var userCollection = db().Database("goTest").Collection("trainers")

// Create Profile or Signup

func createProfile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var person user
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		fmt.Print(err)
	}
	insertResult, err := userCollection.InsertOne(context.TODO(), person)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult)
	json.NewEncoder(w).Encode(insertResult.InsertedID)

}

// Get Profile of a particular User by Name

func getUserProfile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var body user
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	var result primitive.M
	err := userCollection.FindOne(context.TODO(), bson.D{{"name", body.Name}}).Decode(&result)
	if err != nil {

		fmt.Println(err)

	}
	json.NewEncoder(w).Encode(result)

}

//Update Profile of User

func updateProfile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	type updateBody struct {
		Name string `json:"name"`
		City string `json:"city"`
	}
	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	filter := bson.D{{"name", body.Name}}
	after := options.After
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}
	update := bson.D{{"$set", bson.D{{"city", body.City}}}}
	updateResult := userCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)
}

//Delete Profile of User

func deleteProfile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"]

	_id, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		fmt.Printf(err.Error())
	}
	opts := options.Delete().SetCollation(&options.Collation{})
	res, err := userCollection.DeleteOne(context.TODO(), bson.D{{"_id", _id}}, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted %v documents\n", res.DeletedCount)
	json.NewEncoder(w).Encode(res.DeletedCount)

}

func getAllUsers(w http.ResponseWriter, r *http.Request) {

	var results []primitive.M
	cur, err := userCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {

		fmt.Println(err)

	}
	for cur.Next(context.TODO()) {

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)
	}
	cur.Close(context.TODO())
	json.NewEncoder(w).Encode(results)
}
