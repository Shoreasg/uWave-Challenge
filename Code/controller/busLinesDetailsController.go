package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
		Bearing      int16   `json:"bearing"`
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

	if len(response.Data) == 0 {
		return c.Status(500).JSON(fiber.Map{"Error": "No bus locations and details are avaliable now"})
	}

	// Map the data from the response to the BusLocations model
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

	var updatedBusLine models.BusLinesDetails
	updateResult, err := mgm.Coll(&models.BusLinesDetails{}).UpdateOne(c.Context(), bson.M{"busLineId": busLineId}, bson.M{"$set": bson.M{"busLocations": busLocations}})
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	if updateResult.ModifiedCount == 0 {
		return c.Status(500).JSON(fiber.Map{"Error": "No documents found with the specified busLineId"})
	}

	return c.JSON(updatedBusLine)
}
