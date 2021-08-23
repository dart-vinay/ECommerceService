package model

import (
	"github.com/garyburd/redigo/redis"
	"github.com/labstack/gommon/log"
	"main/db"
	"main/utils"
	"strings"
)

// Primary Category
const (
	KIDS  = "KIDS"
	WOMEN = "WOMEN"
	MEN   = "MEN"
)

//Secondary Category
const (
	SHIRT = "SHIRT"
)
type EmptyResponse struct {

}

type Category struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ParentId    []string `json:"parent_id"` //Saved as List of Category ID separated by ## in DB
	PhotoUrl    string   `json:"photo_url"`
}

type CreateCategoryListingRequest struct {
	HeaderInfo
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ParentId    []string `json:"parent_id"`
	PhotoUrl    string   `json:"photo_url"`
}

type FetchCategoryRequest struct {
	HeaderInfo
}

type FetchAllCategoriesRequest struct {
	HeaderInfo
}

type FetchProductsForCategoryRequest struct {
	HeaderInfo
}

type UpdateCategoryInfoRequest struct {
	HeaderInfo
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ParentId    []string `json:"parent_id"`
	PhotoUrl string 	`json:"photo_url"`
}

type FetchCategoryDetailsResponse struct {
	Category
}

type FetchAllCategoryResponse struct {
	Categories []Category `json:"categories"`
}

type FetchProductsForCategoryResponse struct {
	Products []Product `json:"products"`
}

// Create New Category
func (req *CreateCategoryListingRequest) CreateCategory() error {
	log.Infof("Inside CreateCategory method")

	// Open and defer close DB connection
	dbConn := db.DBConn()
	//defer dbConn.Close()

	categoryId := strings.ToUpper(req.Name)
	parentId := ""
	for _, id := range req.ParentId {
		parentId = parentId + id + "##"
	}
	parentId = strings.Trim(parentId, "##")

	queryStmt := "Insert into Locklly.Category (ID, Name,  Description, ParentId, PhotoUrl) values(?,?,?,?,?)"
	insert, err := dbConn.Query(queryStmt, categoryId, req.Name, req.Description, parentId, req.PhotoUrl)
	defer insert.Close()
	if err != nil {
		log.Errorf("Error while creating category info %v", err)
		return err
	}
	return nil
}

// Fetch All Categories
func (req *FetchAllCategoriesRequest) FetchAllCategory() (FetchAllCategoryResponse, error) {
	log.Infof("Inside FetchAllCategory method")

	// Open and defer close DB connection
	dbConn := db.DBConn()
	//defer dbConn.Close()

	response := FetchAllCategoryResponse{}

	queryStmt := "Select Category.ID, Category.Name, Category.Description, Category.ParentId, Category.PhotoUrl from Locklly.Category"
	results, err := dbConn.Query(queryStmt)
	defer results.Close()
	if err != nil {
		log.Errorf("Error while creating category info %v", err)
		return response, err
	}
	for results.Next() {
		category := Category{}
		parentIdString := ""
		err = results.Scan(&category.Id, &category.Name, &category.Description, &parentIdString, &category.PhotoUrl)
		if err != nil {
			log.Errorf("Error while reading Category values from DB %v", err)
			continue
		}
		category.ParentId = utils.FetchListFromString(parentIdString, "##")
		response.Categories = append(response.Categories, category)
	}
	return response, nil
}

// Fetch Category Info
func (req *FetchCategoryRequest) FetchCategoryByCategoryId(categoryId string) (Category, error) {
	log.Infof("Inside FetchCategoryByCategoryId method")
	result, err := FetchCategoriesByCategoryIds([]string{categoryId})
	return result[categoryId], err
}

// Update Category Info
func (req *UpdateCategoryInfoRequest) UpdateCategoryInfo(categoryId string) error {
	log.Infof("Inside UpdateCategoryInfo method")

	// Open and defer close DB connection
	dbConn := db.DBConn()
	//defer dbConn.Close()

	parentId := ""
	for _, id := range req.ParentId {
		parentId = parentId + id + "##"
	}
	parentId = strings.Trim(parentId, "##")

	queryStmt := "Update Locklly.Category SET Name='" + req.Name + "' , Description='" + req.Description + "', ParentId='" + parentId + "', PhotoUrl='"+req.PhotoUrl+"' where ID='" + categoryId + "'"

	update, err := dbConn.Query(queryStmt)
	defer update.Close()
	if err != nil {
		log.Errorf("Error while updating Category info %v", err)
		return err
	}
	return nil
}

// Fetch Categories by Category IDs
func FetchCategoriesByCategoryIds(categories []string) (map[string]Category, error) {
	log.Infof("Inside FetchCategoriesByCategoryIds method")

	// Open and defer close DB connection
	dbConn := db.DBConn()
	//defer dbConn.Close()

	response := make(map[string]Category)
	categoryQueryString := ""
	for _, catValue := range categories {
		categoryQueryString = categoryQueryString + "'" + catValue + "'" + ","
	}
	categoryQueryString = strings.Trim(categoryQueryString, ",")
	queryStmt := "Select Category.ID, Category.Name, Category.Description, Category.ParentId, Category.PhotoUrl from Locklly.Category where ID in (" + categoryQueryString + ")"
	results, err := dbConn.Query(queryStmt)
	defer results.Close()
	if err != nil {
		log.Errorf("Error while making query for fetching categories for category id %v", err)
		return response, err
	}

	for results.Next() {
		category := Category{}
		parentIdString := ""
		err = results.Scan(&category.Id, &category.Name, &category.Description, &parentIdString, &category.PhotoUrl)
		if err != nil {
			log.Errorf("Error while reading Category values from DB %v", err)
			continue
		}
		category.ParentId = utils.FetchListFromString(parentIdString, "##")
		response[category.Id] = category
	}
	return response, err
}

//Fetch Products for Category
func (req *FetchProductsForCategoryRequest) FetchProductsForCategory(categoryId string) (FetchProductsForCategoryResponse, error) {
	response := FetchProductsForCategoryResponse{}

	redisConn := db.GetRedisConnFromPool()
	defer redisConn.Close()

	productIds, err := redis.Strings(redisConn.Do("SMEMBERS", utils.CATEGORY+categoryId))
	if err != nil {
		log.Errorf("Error while fetching product ids for a category %v", err)
		return response, err
	}
	if len(productIds) == 0 {
		return FetchProductsForCategoryResponse{}, nil
	}

	products, err := FetchProductsByProductIds(productIds, false)
	if err != nil {
		log.Errorf("Error returned while Fetching Products for Category: %v, Error: %v", categoryId, err)
		return response, err
	}

	for _, val := range products {
		response.Products = append(response.Products, val)
	}
	return response, nil
}
