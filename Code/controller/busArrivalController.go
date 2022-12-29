package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Shoreasg/uWave-Challenge/models"
	"github.com/alouche/go-geolib"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func CalculateArrivalTime(c *fiber.Ctx) error {
	// Parse the bus line and bus stop ID from the request parameters
	busLineId := c.Params("busLineId")
	busStopId := c.Params("busStopId")

	var busLine models.BusLines
	if err := mgm.Coll(&models.BusLinesDetails{}).First(bson.M{"busLineId": busLineId}, &busLine); err != nil {
		return c.Status(400).JSON(fiber.Map{"Error": "Invalid busLineId"})
	}
	// we we get the latest location of the buses by calling our own API
	// /bus-lines-details/{busLineId}
	http.Get("http://localhost:3000/bus-lines-details/" + busLineId)

	// Query the database for the details of the specified bus line
	var busLineDetails models.BusLinesDetails
	if err := mgm.Coll(&models.BusLinesDetails{}).FindOne(c.Context(), bson.M{"busLineId": busLineId}).Decode(&busLineDetails); err != nil {
		return c.Status(400).JSON(fiber.Map{"Error": "Invalid busLineId"})
	}

	// Find the bus stop in the bus line details
	var busStopLat, busStopLng float64
	for _, busStop := range busLineDetails.BusLineBusStops {
		if busStop.Id == busStopId {
			busStopLat = busStop.Lat
			busStopLng = busStop.Lng
			break
		}
	}

	if busStopLat == 0 && busStopLng == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Invalid bus stop ID"})
	}
	// Assume the speed of the bus in kilometers per hour
	const speed = 50

	// Calculate the arrival time for each bus
	var arrivalTimes []string
	var distanceBetweenBusAndBusStop []string
	for _, location := range busLineDetails.BusLocations {
		// Calculate the distance between the bus and the bus stop using the Haversine formula
		distance := geolib.HaversineDistance(busStopLat, busStopLng, location.Lat, location.Lng)

		// Calculate the arrival time in hours using the assumed speed and the distance
		arrivalTime := distance / speed

		//distance between bus and bus Stop rounded to 2 decimal places
		roundedDistance := fmt.Sprintf("%.2f", distance)

		// Convert the arrival time from hours to seconds
		arrivalTimeInSeconds := int(arrivalTime * 3600)

		//To convert an integer number of units to a Duration, multiply:
		arrivalDuration := time.Duration(arrivalTimeInSeconds) * time.Second

		//convert it to mins and seconds
		arrivalTimeString := arrivalDuration.String()

		//append the arrival times into the slice and also append the distance between to the slice

		arrivalTimes = append(arrivalTimes, arrivalTimeString)
		distanceBetweenBusAndBusStopString := fmt.Sprintf("%v km", roundedDistance)
		distanceBetweenBusAndBusStop = append(distanceBetweenBusAndBusStop, distanceBetweenBusAndBusStopString)
	}

	// Create a BusArrivalDetails object for each bus
	var busArrivalDetails []models.BusArrivalDetails
	for i, location := range busLineDetails.BusLocations {
		busArrivalDetail := models.BusArrivalDetails{
			BusLineId:                    busLineId,
			DistanceBetweenBusAndBusStop: distanceBetweenBusAndBusStop[i],
			BusStopId:                    busStopId,
			BusArrivalTime:               arrivalTimes[i],
			BusVehiclePlate:              location.VehiclePlate,
		}
		//append it to slice
		busArrivalDetails = append(busArrivalDetails, busArrivalDetail)
	}

	// Create a BusArrival object to store the bus arrival details
	busArrival := models.BusArrival{
		BusArrivalDetails: busArrivalDetails,
	}

	// Save the BusArrival object to the database. In the future, can use it to analyze and create even more accurate data
	err := mgm.Coll(&models.BusArrival{}).Create(&busArrival)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Error": "Error creating new records in DB"})
	}

	// Return the calculated arrival times
	return c.JSON(fiber.Map{"data": busArrivalDetails})
}
