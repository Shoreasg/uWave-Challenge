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
	app.Post("/seed-bus-lines-details", controller.SeedBusLineDetailsData)
	app.Get("/bus-lines", controller.GetBusLines)
	app.Get("/bus-lines-details/:busLineId", controller.GetBusLinesDetails)
	app.Get("/bus-line/:busLineId/bus-stop/:busStopId/time", controller.CalculateArrivalTime)
	app.Listen(port)

}
