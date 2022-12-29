package models

import "github.com/kamva/mgm/v3"

type BusLinesDetails struct {
	mgm.DefaultModel `bson:",inline"`
	BusLineId        string            `json:"busLineId" bson:"busLineId"`
	BusLineBusStops  []BusLineBusStops `json:"busStops" bson:"busStops"`
	BusLinePaths     []BusLinePaths    `json:"busPaths" bson:"busPaths"`
	BusLocations     []BusLocations    `json:"busLocations" bson:"busLocations"`
}

type BusLineBusStops struct {
	Id          string  `json:"Id" bson:"Id"`
	Lat         float64 `json:"Lat" bson:"Lat"`
	Lng         float64 `json:"Lng" bson:"Lng"`
	BusStopName string  `json:"busStopName" bson:"busStopName"`
}

type BusLinePaths struct {
	Lat float64 `json:"Lat" bson:"Lat"`
	Lng float64 `json:"Lng" bson:"Lng"`
}

type BusLocations struct {
	Bearing      float64 `json:"bearing" bson:"bearing"`
	Lat          float64 `json:"Lat" bson:"Lat"`
	Lng          float64 `json:"Lng" bson:"Lng"`
	CrowdLevel   string  `json:"crowdLevel" bson:"crowdLevel"`
	VehiclePlate string  `json:"vehiclePlate" bson:"vehiclePlate"`
}
