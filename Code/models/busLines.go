package models

import "github.com/kamva/mgm/v3"

type BusLines struct {
	mgm.DefaultModel `bson:",inline"`
	BusLineId        string `json:"id" bson:"id"`
	BusLineName      string `json:"fullName" bson:"fullName"`
	BusLineShortName string `json:"shortName" bson:"shortName"`
}
