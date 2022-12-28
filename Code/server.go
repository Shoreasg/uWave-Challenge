package main

import (
	"fmt"
	"os"

	"github.com/Shoreasg/uWave-Challenge/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	loadEnv := godotenv.Load() //try to load env file
	if loadEnv != nil {        //if loading fail, exit the app
		panic("Error loading .env file")
	}
	mongoUri := os.Getenv("MONGO_URI")
	//connect to DB instance

	connectDB := mgm.SetDefaultConfig(nil, "uWave", options.Client().ApplyURI(mongoUri))

	if connectDB != nil { //if connect fail, exit the app
		panic(connectDB)
	}

	fmt.Printf("Successfully load env file and connect to DB")

	app := fiber.New() // initalize new fiber app
	port := os.Getenv("PORT")
	app.Post("/seed-bus-lines", controller.SeedData)
	app.Get("/bus-lines", controller.GetBusLines)
	app.Listen(port)

}
