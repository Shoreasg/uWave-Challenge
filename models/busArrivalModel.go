package models

import "github.com/kamva/mgm/v3"

type BusArrival struct {
	mgm.DefaultModel  `bson:",inline"`
	BusArrivalDetails []BusArrivalDetails `json:"busArrivalDetails" bson:"busArrivalDetails"`
}

type BusArrivalDetails struct {
	BusLineId                    string `json:"busLineId" bson:"busLineId"`
	BusStopId                    string `json:"busStopId" bson:"busStopId"`
	DistanceBetweenBusAndBusStop string `json:"distance" bson:"distance"`
	BusArrivalTime               string `json:"arrivalTime" bson:"arrivalTime"`
	BusVehiclePlate              string `json:"vehiclePlate" bson:"vehiclePlate"`
}
