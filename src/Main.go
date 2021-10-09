////////////////ansh shukla api for instagram using golang and mongodb for database storage//////////////////////
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

///////////structure//////////////

// Users should have the following attributes
// Id{
// "Name":"",
// "Email":"",
// "Password":""}

type Person struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" bson:"username,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"-,omitempty" bson:"-,omitempty"` //the value "-" ensures we are not able to see the password value in the get operation hence ensuring security :)
}

// Posts should have the following Attributes. All fields are mandatory unless marked optional:
// Id
// {"Caption":"",
// "Img":"",
// "Posted_timestamp":""}
type Post struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Posted_timestamp string             `json:"Posted_timestamp,omitempty" bson:"Posted_timestamp,omitempty"`
	Caption          string             `json:"Caption,omitempty" bson:"Caption,omitempty"`
	Img              string             `json:"img,omitempty" bson:"img,omitempty"`
	Authid           string             `json:"-,omitempty" bson:"-,omitempty"`
}

///////////////// users /////////////////////////

//POST users info in database instauserdatabase//
func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("instauserdatabase").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

//GET users info in database instauserdatabase//
func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []Person
	collection := client.Database("instauserdatabase").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(people)
}

//
func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Person
	collection := client.Database("instauserdatabase").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(person)
}

////////////////// post /////////////////////////

//POST the insta post details in database instauserdatabase//
func CreatePostEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var post Post
	_ = json.NewDecoder(request.Body).Decode(&post)
	collection := client.Database("instauserdatabase").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, post)
	json.NewEncoder(response).Encode(result)
}

//GET the insta post details from the database instauserdatabase//
func GetPostEndpoint(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Function Ran..")
	response.Header().Set("content-type", "application/json")
	var post []Post
	collection := client.Database("instauserdatabase").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": " + err.Error() + " }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var posts Post
		cursor.Decode(&posts)
		post = append(post, posts)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": " + err.Error() + " }`))
		return
	}
	json.NewEncoder(response).Encode(post)
}

//GET all the posts of a user using auth id as the "key"//
func Getusersposts(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	Auth, _ := params["id"]
	var post []Post
	collection := client.Database("instauserdatabase").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, Post{Authid: Auth})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": " + err.Error() + " }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var posts Post
		cursor.Decode(&posts)
		post = append(post, posts)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": " + err.Error() + " }`))
		return
	}
	json.NewEncoder(response).Encode(post)
}

//////////////////////////////////////////
func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	///////users functions: get,post//////
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}", GetPersonEndpoint).Methods("GET")
	///////posts functions : get,post//////
	router.HandleFunc("/posts", CreatePostEndpoint).Methods("POST")
	router.HandleFunc("/person/posts/{id}", Getusersposts).Methods("GET")
	router.HandleFunc("/posts/{id}", GetPostEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}
