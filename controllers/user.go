package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/airbornharsh/go-mongodb/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	client *mongo.Client
}

func NewUserController(c *mongo.Client) *UserController {
	return &UserController{c}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	Id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": Id}
	u := models.User{}

	err := uc.client.Database("go-mongodb").Collection("users").FindOne(context.Background(), filter).Decode(&u)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := models.User{}

	json.NewDecoder(r.Body).Decode(&u)

	u.ID = primitive.NewObjectID()

	inserted, err := uc.client.Database("go-mongodb").Collection("users").InsertOne(context.Background(), u)

	if err != nil {
		if writeErr, ok := err.(mongo.WriteException); ok {
			for _, we := range writeErr.WriteErrors {
				if we.Code == 11000 {
					log.Printf("Duplicate key error: %v\n", we)
					// Handle the duplicate key error here, for example, inform the user.
				} else {
					log.Printf("Write error: %v\n", we)
				}
			}
			http.Error(w, "Error while inserting data", http.StatusInternalServerError)
			return
		}
	}
	fmt.Println("Inserted with Id: ", inserted.InsertedID)

	uj, err := json.Marshal(u)

	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	Id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": Id}

	result, err := uc.client.Database("go-mongodb").Collection("users").DeleteOne(context.Background(), filter)

	if err == nil {
		w.WriteHeader(404)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted User", result.DeletedCount, "\n")
}
