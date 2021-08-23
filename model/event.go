package model

import (
	"github.com/labstack/gommon/log"
	"main/db"
	"main/utils"
)

type Event string


func ToEventList(events []string) []Event {
	eventList := []Event{}
	for _, event := range events {
		eventList = append(eventList, Event(event))
	}
	return eventList
}

func FetchEventForProduct(productId string) ([]string, error) {

	// Open and defer close DB connection
	dbConn := db.DBConn()
	//defer dbConn.Close()

	eventList := []string{}
	eventString := ""
	queryCategory := "Select Events from Locklly.Product where ID='" + productId + "'"
	results, err := dbConn.Query(queryCategory)
	defer results.Close()
	if err != nil {
		log.Errorf("Error while fetching category for product from DB %v", err)
		return eventList, err
	}
	if results.Next() {
		err := results.Scan(&eventString)
		if err != nil {
			log.Errorf("Error while scanning for categories for product id")
		}
	}
	eventList = utils.FetchListFromString(eventString, "##")
	return eventList, nil
}
