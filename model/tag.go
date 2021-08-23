package model

import (
	"github.com/labstack/gommon/log"
	"main/db"
	"main/utils"
)

func FetchTagForProduct(productId string) ([]string, error) {
	// Open and defer close DB connection
	dbConn := db.DBConn()
	//defer dbConn.Close()

	tagList := []string{}
	tagString := ""
	queryCategory := "Select Tags from Locklly.Product where ID='" + productId + "'"
	results, err := dbConn.Query(queryCategory)
	defer results.Close()
	if err != nil {
		log.Errorf("Error while fetching category for product from DB %v", err)
		return tagList, err
	}
	if results.Next() {
		err := results.Scan(&tagString)
		if err != nil {
			log.Errorf("Error while scanning for categories for product id")
		}
	}
	tagList = utils.FetchListFromString(tagString, "##")
	return tagList, nil
}
