package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/airbornharsh/go-mongodb/controllers"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "gopkg.in/mgo.v2"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	r := httprouter.New()
	uc := controllers.NewUserController(getSession())
	r.POST("/user", uc.CreateUser)
	r.GET("/user/:id", uc.GetUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	http.ListenAndServe("localhost:8080", r)
}

func getSession() *mongo.Client {
	opts := options.Client().ApplyURI(goDotEnvVariable("MONGODB_URI"))

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to mongodb")
	return client
}
