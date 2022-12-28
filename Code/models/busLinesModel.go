package models

import "github.com/kamva/mgm/v3"

type BusLines struct {
	mgm.DefaultModel `bson:",inline"`
	BusLineId        string `json:"busLineId" bson:"busLineId"`
	BusLineName      string `json:"busLineName" bson:"busLineName"`
	BusLineShortName string `json:"busLineShortName" bson:"busLineShortName"`
}
