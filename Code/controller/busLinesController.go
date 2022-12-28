package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Shoreasg/uWave-Challenge/models"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gofiber/fiber/v2"
)

type busLinesResponse struct {
	Data []struct {
		ID        string `json:"id"`
		FullName  string `json:"fullName"`
		ShortName string `json:"shortName"`
	} `json:"payload"`
}

func SeedData(c *fiber.Ctx) error {
	// check if there is any exisitng buslines in DB
	var busLinesInDB []models.BusLines
	if err := mgm.Coll(&models.BusLines{}).SimpleFind(&busLinesInDB, bson.M{}); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if len(busLinesInDB) == 0 { //if db doesn't contain any busline, call the API and seed the data
		resp, err := http.Get("https://test.uwave.sg/busLines")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"Error": "Error Making get Request"})
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"Error": "Error readng response body"})
		}

		var response busLinesResponse
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Error unmarshalling JSON data:", err)
			return c.Status(500).SendString("Error unmarshalling JSON data")
		}

		var busLines []models.BusLines

		for _, d := range response.Data {
			busLines = append(busLines, models.BusLines{
				BusLineId:        d.ID,
				BusLineName:      d.FullName,
				BusLineShortName: d.ShortName,
			})
		}

		for i := 0; i < len(busLines); i++ {
			if err := mgm.Coll(&busLines[i]).Create(&busLines[i]); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
		}

		return c.JSON(fiber.Map{"Success": "Seed data successfully"})
	} else {
		//check if the number of busLines in DB is lesser than the API, if yes, means there is an additional new busLines. Insert the new data
		var busLinesInDB []models.BusLines
		if err := mgm.Coll(&models.BusLines{}).SimpleFind(&busLinesInDB, bson.M{}); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		resp, err := http.Get("https://test.uwave.sg/busLines")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"Error": "Error Making get Request"})
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"Error": "Error readng response body"})
		}

		var response busLinesResponse
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Error unmarshalling JSON data:", err)
			return c.Status(500).SendString("Error unmarshalling JSON data")
		}

		var busLines []models.BusLines

		for _, d := range response.Data {
			busLines = append(busLines, models.BusLines{
				BusLineId:        d.ID,
				BusLineName:      d.FullName,
				BusLineShortName: d.ShortName,
			})
		}
		if len(busLines) > len(busLinesInDB) {
			if err := mgm.Coll(&models.BusLines{}).Drop(c.Context()); err != nil {
				return c.Status(500).SendString(err.Error())
			}

			for i := 0; i < len(busLines); i++ {
				if err := mgm.Coll(&busLines[i]).Create(&busLines[i]); err != nil {
					return c.Status(500).SendString(err.Error())
				}
			}
			return c.JSON(fiber.Map{"Success": "new data seeded successfully"})
		} else {
			return c.JSON(fiber.Map{"Success": "No new data needed to be seeded"})
		}

	}
}

func GetBusLines(c *fiber.Ctx) error {
	var busLines []models.BusLines
	if err := mgm.Coll(&models.BusLines{}).SimpleFind(&busLines, bson.M{}); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{"success": busLines})
}
