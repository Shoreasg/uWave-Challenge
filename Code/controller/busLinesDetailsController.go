package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Shoreasg/uWave-Challenge/models"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type busLinesDetailsResponse struct {
	Data []struct {
		ID       string             `json:"id"`
		BusStops []busStopsResponse `json:"busStops"`
		Path     [][]float64        `json:"path"`
	} `json:"payload"`
}

type busStopsResponse struct {
	Id          string  `json:"id"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	BusStopName string  `json:"name"`
}

type busLocationsResponse struct {
	Data []struct {
		Bearing      float64 `json:"bearing"`
		Lat          float64 `json:"lat"`
		Lng          float64 `json:"lng"`
		CrowdLevel   string  `json:"crowdLevel"`
		VehiclePlate string  `json:"vehiclePlate"`
	} `json:"payload"`
}

func SeedBusLineDetailsData(c *fiber.Ctx) error {

	resp, err := http.Get("https://test.uwave.sg/busLines")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Error": "Error Making get Request"})
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Error": "Error readng response body"})
	}

	var response busLinesDetailsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error unmarshalling JSON data:", err)
		return c.Status(500).JSON(fiber.Map{"Error": "Error unmarshalling JSON data"})
	}

	var busLinesDetails []models.BusLinesDetails

	for _, d := range response.Data {
		var busLineBusStops []models.BusLineBusStops
		for _, bs := range d.BusStops {
			busLineBusStops = append(busLineBusStops, models.BusLineBusStops{
				Id:          bs.Id,
				Lat:         bs.Lat,
				Lng:         bs.Lng,
				BusStopName: bs.BusStopName,
			})
		}

		var busLinePaths []models.BusLinePaths
		for _, p := range d.Path {
			busLinePaths = append(busLinePaths, models.BusLinePaths{
				Lat: p[0],
				Lng: p[1],
			})
		}

		busLinesDetails = append(busLinesDetails, models.BusLinesDetails{
			BusLineId:       d.ID,
			BusLineBusStops: busLineBusStops,
			BusLinePaths:    busLinePaths,
		})
	}

	for _, busLineDetails := range busLinesDetails {
		if err := mgm.Coll(&busLineDetails).Create(&busLineDetails); err != nil {
			return c.Status(500).JSON(fiber.Map{"Error": err.Error()})
		}
	}

	return c.JSON(fiber.Map{"Success": "Seed data successfully"})

}

func GetBusLinesDetails(c *fiber.Ctx) error {
	busLineId := c.Params("busLineId")

	// Check if user enter a valid busLineId
	var busLine models.BusLines
	if err := mgm.Coll(&models.BusLinesDetails{}).First(bson.M{"busLineId": busLineId}, &busLine); err != nil {
		return c.Status(400).JSON(fiber.Map{"Error": "Invalid busLineId"})
	}

	// Call the external API to get the data for the specified bus line
	resp, err := http.Get("https://test.uwave.sg/busPositions/" + busLineId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Error": "Error Making get Request"})
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Error": "Error readng response body"})
	}

	// Unmarshal the response data into a struct
	var response busLocationsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error unmarshalling JSON data:", err)
		return c.Status(500).JSON(fiber.Map{"Error": "Error unmarshalling JSON data"})
	}

	// if there is no buslocation at the time, just send the db document with empty locations which means there isn't bus at that time
	if len(response.Data) == 0 {
		var busLineDetails models.BusLinesDetails
		if err := mgm.Coll(&models.BusLinesDetails{}).FindOne(c.Context(), bson.M{"busLineId": busLineId}).Decode(&busLineDetails); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(fiber.Map{"Data": busLineDetails})
	}

	// Map the data from the response to the busLocations slice
	var busLocations []models.BusLocations
	for _, d := range response.Data {
		busLocations = append(busLocations, models.BusLocations{
			Bearing:      d.Bearing,
			Lat:          d.Lat,
			Lng:          d.Lng,
			CrowdLevel:   d.CrowdLevel,
			VehiclePlate: d.VehiclePlate,
		})
	}

	// update the DB with the latest bus locations and details
	var busLineDetails models.BusLinesDetails
	_, err = mgm.Coll(&models.BusLinesDetails{}).UpdateOne(c.Context(), bson.M{"busLineId": busLineId}, bson.M{"$set": bson.M{"busLocations": busLocations, "updated_at": time.Now()}})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Error": "Error updating bus locations"})
	}

	// use the find function and return the whole document
	if err := mgm.Coll(&models.BusLinesDetails{}).First(bson.M{"busLineId": busLineId}, &busLineDetails); err != nil {
		return c.Status(500).JSON(fiber.Map{"Error": "Error finding bus line details"})
	}

	return c.JSON(fiber.Map{"Data": busLineDetails})
}
