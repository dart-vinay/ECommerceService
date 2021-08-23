package model

import (
	"database/sql"
	"github.com/garyburd/redigo/redis"
	"github.com/labstack/gommon/log"
	"main/db"
	"main/utils"
	"strconv"
	"strings"
	"time"
)

type Product struct {
	Id               string     `json:"id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Categories       []Category `json:"categories,omitempty"`
	Tags             []string   `json:"tags"`
	Events           []Event    `json:"events"`
	CreatedBy        Merchant   `json:"created_by,omitempty"`
	CreatedAt        time.Time  `json:"created_at,omitempty"`
	UpdatedAt        time.Time  `json:"updated_at"`
	Price            float64    `json:"price"`
	Likes            int64      `json:"likes"`
	Rating           float64    `json:"rating"`
	Sizes            []string   `json:"sizes"`
	CountPeopleRated int64      `json:"count_people_rated"`
	Views            int64      `json:"views"`
	PhotoUrls        []string   `json:"photo_urls"`
	PrimaryURL       string     `json:"primary_url"`
	InterestingFact  string     `json:"interesting_fact"`

	// Internal Information
	State            string   `json:"state,omitempty"` // Published, Reviewed, Pending, Deleted etc.
	CategoryList     []string `json:"-"`
	EventList        []string `json:"-"`
	MerchantUserName string   `json:"-"`
}

type ProductObject struct {
	Id               string    `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Categories       []string  `json:"categories,omitempty"`
	Tags             []string  `json:"tags"`
	Events           []Event   `json:"events"`
	CreatedBy        string    `json:"created_by,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at"`
	Price            float64   `json:"price"`
	Likes            int64     `json:"likes"`
	Rating           float64   `json:"rating"`
	Sizes            []string  `json:"sizes"`
	CountPeopleRated int64     `json:"count_people_rated"`
	Views            int64     `json:"views"`
	PhotoUrls        []string  `json:"photo_urls"`
	PrimaryURL       string    `json:"primary_url"`
	InterestingFact  string    `json:"interesting_fact"`

	// Internal Information
	State string `json:"state,omitempty"` // Published, Reviewed, Pending, Deleted etc.
}

const (
	PUBLISHED = "PUBLISHED"
	REVIEWED  = "REVIEWED"
	PENDING   = "PENDING"
	DELETED   = "DELETED"
)

type CreateProductListingRequest struct {
	HeaderInfo
	Id              string   `json:"id"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Categories      []string `json:"categories"`
	PhotoUrls       []string `json:"photo_urls"`
	PrimaryURL      string   `json:"primary_url"`
	Sizes           []string `json:"sizes"`
	Price           float64  `json:"price"`
	Tags            []string `json:"tags"`
	Events          []string `json:"events"`
	InterestingFact string   `json:"interesting_fact,omitempty"`
}

type FetchProductDetailsRequest struct {
	HeaderInfo
	OnlyBasicInfo bool
}
type FetchProductDetailsResponse struct {
	Product
}

type ProductLikeRequest struct {
	HeaderInfo
	ProductId string `json:"product_id"`
	IsLiked   bool   `json:"is_liked"`
}

type ProductBookmarkRequest struct {
	HeaderInfo
	ProductId    string `json:"product_id"`
	IsBookmarked bool   `json:"is_bookmarked"`
}

type FetchProductLikedByUserRequest struct {
	HeaderInfo
}

type FetchProductBookmarkedByUser struct {
	HeaderInfo
}

type ProductLikedByUserResponse struct {
	Products []Product `json:"products"`
}

type ProductBookmarkedByUserResponse struct {
	Products []Product `json:"products"`
}

type UpdateProductInfoRequest struct {
	HeaderInfo
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Categories      []string `json:"categories"`
	Tags            []string `json:"tags"`
	Events          []string `json:"events"`
	PhotoUrls       []string `json:"photo_urls"`
	PrimaryURL      string   `json:"primary_url"`
	InterestingFact string   `json:"interesting_fact"`
	Sizes           []string `json:"sizes"`
	Price           float64  `json:"price"`
}

//Create a Product
func (req *CreateProductListingRequest) CreateProductListing() error {
	log.Infof("Inside CreateProductListing method")

	// Open and defer close DB connection
	dbConn := db.DBConn()
	//defer dbConn.Close()

	redisConn := db.GetRedisConnFromPool()
	defer redisConn.Close()

	productId := GenerateRandomString(8)
	loc, _ := time.LoadLocation("Asia/Kolkata")

	log.Infof("%v", time.Now().In(loc))
	if len(req.PhotoUrls) == 0 || req.Price == 0 {
		return utils.ErrBadRequest
	}
	if req.PrimaryURL == "" {
		req.PrimaryURL = req.PhotoUrls[0]
	}

	categories := utils.StringFromArray(req.Categories, "##")
	photos := utils.StringFromArray(utils.Unique(req.PhotoUrls), "##")
	events := utils.StringFromArray(utils.Unique(utils.ToScreamingSnakeArray(req.Events)), "##")
	tags := utils.StringFromArray(utils.Unique(req.Tags), "##")
	availableSizes := ""
	for _, size := range req.Sizes {
		if utils.ValidateSize(strings.ToUpper(size)) {
			availableSizes = availableSizes + size + "##"
		}
	}

	availableSizes = strings.Trim(availableSizes, "##")

	queryStmt := "Insert into Locklly.Product (ID, Title,  Description, CreatedBy, CreatedAt, UpdatedAt, PhotoUrls, PrimaryPhotoUrl, Category, Events, Tags, Sizes, Price, InterestingFact, State) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err := dbConn.Query(queryStmt, productId, req.Title, req.Description, req.UserName, time.Now().In(loc), time.Now().In(loc), photos, req.PhotoUrls[0], categories, events, tags, availableSizes, req.Price, req.InterestingFact, PENDING)
	if err != nil {
		log.Errorf("Error while inserting product in products table %v", err)
		return err
	}
	for _, category := range utils.Unique(req.Categories) {
		err := redisConn.Send("SADD", utils.CATEGORY+category, productId)
		if err != nil {
			log.Errorf("Error while adding product for the category in redis %v", err)
		}
	}
	for _, event := range utils.Unique(utils.ToScreamingSnakeArray(req.Events)) {
		err := redisConn.Send("SADD", utils.EVENT+event, productId)
		if err != nil {
			log.Errorf("Error while adding product for the event in redis %v", err)
		}
		err = redisConn.Send("SADD", utils.ALL_EVENTS, event)
		if err != nil {
			log.Errorf("Error while adding event in redis %v", err)
		}
	}
	for _, tag := range utils.Unique(req.Tags) {
		err := redisConn.Send("SADD", utils.TAG+tag, productId)
		if err != nil {
			log.Errorf("Error while adding product for the tag in redis %v", err)
		}
		err = redisConn.Send("SADD", utils.ALL_TAGS, tag)
		if err != nil {
			log.Errorf("Error while adding tag in redis %v", err)
		}
	}

	err = redisConn.Flush()
	if err != nil {
		log.Errorf("Error while flushing to Redis %v", err)
	}
	return nil
}

//Fetch Product Details
func (req *FetchProductDetailsRequest) FetchProductDetailById(productId string) (FetchProductDetailsResponse, error) {
	log.Info("Inside FetchProductDetailById method")
	response := FetchProductDetailsResponse{}
	if productInfo, err := FetchProductById(productId, req.OnlyBasicInfo); err != nil {
		log.Errorf("Error while trying to fetch product info for product id: %v %v", productId, err)
		return response, err
	} else {
		response.Product = productInfo
	}
	return response, nil
}

//func FetchAllProducts() ([]Product, error){
//
//}

func FetchProductById(productId string, onlyBasicInfo bool) (Product, error) {
	response, err := FetchProductsByProductIds([]string{productId}, onlyBasicInfo)
	if err != nil {
		log.Errorf("Error while fetching product info %v", err)
		return Product{}, err
	} else {
		return response[productId], nil
	}
}

func FetchProductsByProductIds(productIds []string, onlyBasicInfo bool) (map[string]Product, error) {
	log.Info("Inside FetchProductsByProductIds method")

	response := make(map[string]Product)

	// Open and defer close DB connection
	dbConn := db.DBConn()
	//defer dbConn.Close()

	productQueryString := utils.FetchInQueryStringFromArray(productIds)

	queryStmt := ""
	if onlyBasicInfo {
		queryStmt = "Select Product.ID, Product.Title, Product.Description, Product.CreatedBy, Product.PrimaryPhotoUrl, Product.Category, Product.Tags, Product.Events from Locklly.Product where ID in (" + productQueryString + ")"
	} else {
		queryStmt = "Select Product.ID, Product.Title, Product.Description, Product.CreatedBy, Product.CreatedAt, Product.Category, Product.PhotoUrls, Product.Likes, Product.Rating, Product.Price, Product.CountPeopleRated, Product.Sizes, Product.Tags, Product.Events, Product.InterestingFact, Product.UpdatedAt, Product.State from Locklly.Product where ID in (" + productQueryString + ")"
	}
	results, err := dbConn.Query(queryStmt)
	if err != nil {
		log.Errorf("Error while making query for fetching products for product ids %v", err)
		return response, err
	}
	for results.Next() {
		product := Product{}
		createdAt := ""
		updatedAt := ""
		photos := ""
		categories := ""
		tags := ""
		events := ""
		productMerchantId := ""
		availableSizes := ""

		if onlyBasicInfo {
			err := results.Scan(&product.Id, &product.Title, &product.Description, &productMerchantId, &product.PrimaryURL, &categories, &tags, &events)
			if err != nil {
				log.Infof("Error while scanning product data from db %v", err)
				return response, err
			}
			product.Tags = utils.FetchListFromString(tags, "##")
			product.EventList = utils.FetchListFromString(events, "##")
			product.CategoryList = utils.FetchListFromString(categories, "##")
			product.MerchantUserName = productMerchantId
			response[product.Id] = product
			continue
		}
		err = results.Scan(&product.Id, &product.Title, &product.Description, &productMerchantId, &createdAt, &categories, &photos, &product.Likes, &product.Rating, &product.Price, &product.CountPeopleRated, &availableSizes, &tags, &events, &product.InterestingFact, &updatedAt, &product.State)
		if err != nil {
			log.Errorf("Error while scanning product data from db %v", err)
			return response, err
		}
		product.CreatedAt, err = utils.ParseStringToTime(createdAt)
		if err != nil {
			log.Errorf("Error parsing time for createdAt or updatedAt, %v", err)
			return response, err
		}
		product.UpdatedAt, err = utils.ParseStringToTime(updatedAt)
		if err != nil {
			log.Errorf("Error parsing time for createdAt or updatedAt, %v", err)
			return response, err
		}
		product.PhotoUrls = utils.FetchListFromString(photos, "##")
		if len(product.PhotoUrls) > 0 {
			product.PrimaryURL = product.PhotoUrls[0]
		}
		product.Sizes = utils.FetchListFromString(availableSizes, "##")
		product.Tags = utils.FetchListFromString(tags, "##")
		product.Events = ToEventList(utils.FetchListFromString(events, "##"))
		categoryList := utils.FetchListFromString(categories, "##")
		categoryMap, err := FetchCategoriesByCategoryIds(categoryList)
		if err != nil {
			log.Errorf("Error While Fetching Categories for Ids inside FetchProductsByProductIds %v", err)
			product.Categories = []Category{}
			response[product.Id] = product
			continue
		}
		product.Likes, err = FetchProductLikes(product.Id)
		if err != nil {
			log.Errorf("Unable to fetch likes for product %v", err)
		}

		for _, val := range categoryMap {
			product.Categories = append(product.Categories, val)
		}
		if productMerchantId == "" {
			product.CreatedBy = Merchant{}
		} else {

			merchantInfoRequest := FetchMerchantInfoRequest{
				HeaderInfo{UserName: productMerchantId, AuthToken: "", IsAdmin: false},
			}
			merchantInfoResponse, err := merchantInfoRequest.FetchMerchantDetailsByID(productMerchantId, dbConn)
			if err != nil {
				if strings.HasPrefix(err.Error(), "User Doesn't Exist") {
					return response, utils.ErrUnauthorized
				}
				return response, err
			}
			product.CreatedBy = merchantInfoResponse.Merchant
		}

		//merchantInfoReq := FetchMerchantInfoForCustomerRequest{productMerchantId}
		//productMerchant, err := merchantInfoReq.FetchMerchantInfoForCustomer(dbConn)
		//if err != nil {
		//	log.Errorf("Error fetching merchant info for product id %v", product.Id)
		//	product.CreatedBy = MerchantInfoObject{}
		//} else {
		//	product.CreatedBy = productMerchant
		//}
		//if otherProductInfo {
		//	product.SameMerchantProducts, err = FetchBasicProductInfoForMerchant(product.CreatedBy.Id, dbConn, []string{product.Id})
		//	if err != nil {
		//		log.Error("Error while fetching Same Merchant Products")
		//	}
		//	product.RelatedProducts, err = FetchBasicProductInfoForCategory(categoryList, dbConn, []string{product.Id})
		//	if err != nil {
		//		log.Error("Error while fetching related Products")
		//	}
		//}
		response[product.Id] = product
	}
	return response, nil
}

func (req *UpdateProductInfoRequest) UpdateProductInfo(productId string) error {
	log.Infof("Inside UpdateProductInfo method")

	// Open and defer close DB connection
	dbConn := db.DBConn()
	//defer dbConn.Close()

	redisConn := db.GetRedisConnFromPool()
	defer redisConn.Close()

	if len(req.PhotoUrls) == 0 || req.Price == 0 {
		return utils.ErrBadRequest
	}
	if req.PrimaryURL == "" {
		req.PrimaryURL = req.PhotoUrls[0]
	}

	oldCategories, err := FetchCategoryForProduct(productId)
	oldEvents, err := FetchEventForProduct(productId)
	oldTags, err := FetchTagForProduct(productId)

	if err != nil {
		log.Errorf("Error encountered while fetching old data for product, %v", err)
		return err
	}

	loc, _ := time.LoadLocation("Asia/Kolkata")

	newTitle := req.Title
	newDesc := req.Description

	newPhotoUrls := utils.StringFromArray(utils.Unique(req.PhotoUrls), "##")
	newCategories := utils.StringFromArray(req.Categories, "##")
	newEvents := utils.StringFromArray(utils.Unique(utils.ToScreamingSnakeArray(req.Events)), "##")
	newTags := utils.StringFromArray(utils.Unique(req.Tags), "##")
	newPrice := strconv.Itoa(int(req.Price))
	newFact := req.InterestingFact

	availableSizes := ""
	for _, size := range req.Sizes {
		if utils.ValidateSize(strings.ToUpper(size)) {
			availableSizes = availableSizes + size + "##"
		}
	}
	availableSizes = strings.Trim(availableSizes, "##")
	updateFieldStatement := utils.CreateUpdateFieldStatement([]string{"Title", "Description", "PhotoUrls", "PrimaryPhotoUrl", "Events", "Category", "Sizes", "Tags", "Price", "InterestingFact", "UpdatedAt"}, []string{newTitle, newDesc, newPhotoUrls, req.PhotoUrls[0], newEvents, newCategories, availableSizes, newTags, newPrice, newFact, time.Now().In(loc).Format("2006-01-02 15:04:05")})
	queryStmt := "Update Locklly.Product Set " + updateFieldStatement + " where ID='" + productId + "'"
	_, err = dbConn.Query(queryStmt)
	if err != nil {
		log.Errorf("Error while updating product in products table %v", err)
		return err
	}
	for _, category := range oldCategories {
		err := redisConn.Send("SREM", utils.CATEGORY+category, productId)
		if err != nil {
			log.Errorf("Error while removing product for the category in redis %v", err)
		}
	}

	for _, category := range req.Categories {
		err := redisConn.Send("SADD", utils.CATEGORY+category, productId)
		if err != nil {
			log.Errorf("Error while adding product for the category in redis %v", err)
		}
	}

	for _, event := range oldEvents {
		err := redisConn.Send("SREM", utils.EVENT+event, productId)
		if err != nil {
			log.Errorf("Error while removing product for the event in redis %v", err)
		}
	}

	for _, event := range utils.ToScreamingSnakeArray(req.Events) {
		err := redisConn.Send("SADD", utils.EVENT+event, productId)
		if err != nil {
			log.Errorf("Error while adding product for the event in redis %v", err)
		}
	}

	for _, tag := range oldTags {
		err := redisConn.Send("SREM", utils.TAG+tag, productId)
		if err != nil {
			log.Errorf("Error while removing product for the tag in redis %v", err)
		}
	}

	for _, tag := range req.Tags {
		err := redisConn.Send("SADD", utils.TAG+tag, productId)
		if err != nil {
			log.Errorf("Error while adding product for the tag in redis %v", err)
		}
	}

	err = redisConn.Flush()
	if err != nil {
		log.Errorf("Error while flushing to Redis %v", err)
	}

	return nil
}

func (req *UpdateProductInfoRequest) Delete(productId string) error {
	log.Infof("Inside Delete Product method")

	// Open and defer close DB connection
	dbConn := db.DBConn()
	//defer dbConn.Close()

	redisConn := db.GetRedisConnFromPool()
	defer redisConn.Close()

	productRequest := FetchProductDetailsRequest{
		OnlyBasicInfo: true,
	}
	product, err := productRequest.FetchProductDetailById(productId)
	if err != nil {
		log.Errorf("Unable to fetch product %v for deletion %v", productId, err)
		return err
	}

	if product.Id == "" {
		return nil
	}

	if req.UserName == product.MerchantUserName || req.IsAdmin {
		err = product.Delete(dbConn, &redisConn)
		if err != nil {
			return err
		}
	} else {
		return utils.ErrBadRequest
	}

	return nil
}

func (product *Product) Delete(dbConn *sql.DB, redisConn *redis.Conn) error {

	// Delete Product from DB
	queryStmt := "Delete from Locklly.Product where Product.ID=?"
	delete, err := dbConn.Query(queryStmt, product.Id)
	defer delete.Close()
	if err != nil {
		log.Errorf("Error deleting product from DB %v", err)
	}

	for _, category := range product.CategoryList {
		err := (*redisConn).Send("SREM", utils.CATEGORY+category, product.Id)
		if err != nil {
			log.Errorf("Error while removing product for the category in redis %v", err)
			return err
		}
	}
	for _, event := range product.EventList {
		err := (*redisConn).Send("SREM", utils.EVENT+event, product.Id)
		if err != nil {
			log.Errorf("Error while removing product for the event in redis %v", err)
			return err
		}
	}

	for _, tag := range product.Tags {
		err := (*redisConn).Send("SREM", utils.TAG+tag, product.Id)
		if err != nil {
			log.Errorf("Error while removing product for the tag list in redis %v", err)
			return err
		}
	}

	log.Infof("Product %v deleted Successfully!", product.Id)

	return nil

}
func FetchCategoryForProduct(productId string) ([]string, error) {

	// Open and defer close DB connection
	dbConn := db.DBConn()
	//defer dbConn.Close()

	categoryList := []string{}
	categoryString := ""
	queryCategory := "Select Category from Locklly.Product where ID='" + productId + "'"
	results, err := dbConn.Query(queryCategory)
	defer results.Close()
	if err != nil {
		log.Errorf("Error while fetching category for product from DB %v", err)
		return categoryList, err
	}
	if results.Next() {
		err := results.Scan(&categoryString)
		if err != nil {
			log.Errorf("Error while scanning for categories for product id")
		}
	}
	categoryList = utils.FetchListFromString(categoryString, "##")
	return categoryList, nil
}

func FetchBasicProductInfoForMerchant(merchantId string, dbConn *sql.DB, removeIds []string) ([]Product, error) {
	response := []Product{}

	queryStmt := "Select ID, Title, Description, PrimaryPhotoUrl, Price from Locklly.Product where CreatedBy='" + merchantId + "' LIMIT 7"
	results, err := dbConn.Query(queryStmt)
	defer results.Close()
	if err != nil {
		log.Errorf("Error while fetching basic product info for merchant from DB %v", err)
		return response, err
	}
	for results.Next() {
		productInfo := Product{}
		photo := ""
		err := results.Scan(&productInfo.Id, &productInfo.Title, &productInfo.Description, &photo, &productInfo.Price)
		if err != nil {
			log.Errorf("Error while scanning for basic product info for merchant %v", err)
			response = append(response, productInfo)
			continue
		}
		if utils.ExistsInArray(removeIds, productInfo.Id) {
			continue
		}
		productInfo.PhotoUrls = []string{photo}
		productInfo.PrimaryURL = photo
		response = append(response, productInfo)
	}
	return response, nil
}

func FetchBasicProductInfoForCategory(categoryIds []string, dbConn *sql.DB, removeIds []string) ([]Product, error) {
	response := []Product{}
	redisConn := db.GetRedisConnFromPool()
	defer redisConn.Close()

	finalProdIds := make(map[string]int)
	for _, categoryId := range categoryIds {
		err := redisConn.Send("SMEMBERS", utils.CATEGORY+categoryId)
		if err != nil {
			log.Error("Error while making query to redis inside FetchBasicProductInfoForCategory")
		}
	}

	err := redisConn.Flush()
	if err != nil {
		log.Errorf("Error while flusing send queries to redis")
		return response, err
	}

	for i, _ := range categoryIds {
		log.Infof("Processing %vth categoryid", i)
		prodIds, err := redis.Strings(redisConn.Receive())
		if err != nil {
			log.Errorf("Error while receiving from redis")
		}
		for _, prodId := range prodIds {
			if utils.ExistsInArray(removeIds, prodId) {
				continue
			}
			finalProdIds[prodId] = 1
		}
	}
	productQueryString := utils.FetchInQueryStringFromMapKeys(finalProdIds)
	if productQueryString == "" {
		return response, nil
	}
	queryStmt := "Select ID, Title, Description, PrimaryPhotoUrl, Price from Locklly.Product where ID in (" + productQueryString + ")  LIMIT 7"
	results, err := dbConn.Query(queryStmt)
	defer results.Close()
	if err != nil {
		log.Errorf("Error while fetching basic product info for same Category from DB %v", err)
		return response, err
	}
	for results.Next() {
		productInfo := Product{}
		photo := ""
		err := results.Scan(&productInfo.Id, &productInfo.Title, &productInfo.Description, &photo, &productInfo.Price)
		if err != nil {
			log.Errorf("Error while scanning for basic product info for merchant %v", err)
			response = append(response, productInfo)
			continue
		}
		productInfo.PhotoUrls = []string{photo}
		productInfo.PrimaryURL = photo
		response = append(response, productInfo)
	}
	return response, nil
}

func FetchProductListView(productIds []string) []Product {

	// Open and defer close DB connection
	dbConn := db.DBConn()
	//defer dbConn.Close()
	products := []Product{}
	productQueryString := utils.FetchInQueryStringFromArray(productIds)
	if productQueryString == "" {
		return products
	}
	queryStmt := "Select ID, Title, Description, PrimaryPhotoUrl, Price from Locklly.Product where ID in (" + productQueryString + ")  LIMIT 7"
	results, err := dbConn.Query(queryStmt)
	defer results.Close()
	if err != nil {
		log.Errorf("Error while fetching basic product info for same Category from DB %v", err)
		return products
	}
	for results.Next() {
		productInfo := Product{}
		photo := ""
		err := results.Scan(&productInfo.Id, &productInfo.Title, &productInfo.Description, &photo, &productInfo.Price)
		if err != nil {
			log.Errorf("Error while scanning for basic product info for merchant %v", err)
			products = append(products, productInfo)
			continue
		}
		productInfo.PhotoUrls = []string{photo}
		productInfo.PrimaryURL = photo
		products = append(products, productInfo)
	}
	return products
}

func (req *ProductLikeRequest) LikeProduct() error {
	redisConn := db.GetRedisConnFromPool()
	defer redisConn.Close()

	if req.IsLiked == true {
		exists, err := redisConn.Do("ZRANK", utils.PRODUCT_LIKES+req.ProductId, req.UserName)
		if err != nil {
			log.Errorf("Error while trying to fetch like by user: %v for product: %v. Error : %v", req.UserName, req.ProductId, err)
			return err
		}
		if exists != nil {
			return nil
		}
		// Like Count
		err = redisConn.Send("ZINCRBY", "PRODUCT_LIKES", 1, req.ProductId)
		if err != nil {
			log.Errorf("Error increasing like for the product with id %v", req.ProductId)
		}
		// Users Liked the Product
		err = redisConn.Send("ZADD", utils.PRODUCT_LIKES+req.ProductId, time.Now().Unix(), req.UserName)
		if err != nil {
			log.Errorf("Error adding like for the product with id %v by user %v", req.ProductId, req.UserName)
		}

		// Products Liked By User
		err = redisConn.Send("ZADD", utils.PRODUCTS_LIKED_BY_USER+req.UserName, time.Now().Unix(), req.ProductId)
		if err != nil {
			log.Errorf("Error adding product with is %v in user liked list %v", req.ProductId, req.UserName)
		}
	} else {
		exists, err := redisConn.Do("ZRANK", utils.PRODUCT_LIKES+req.ProductId, req.UserName)
		if err != nil {
			log.Errorf("Error while trying to fetch like by user: %v for product: %v. Error : %v", req.UserName, req.ProductId, err)
			return err
		}
		if exists == nil {
			return nil
		}
		err = redisConn.Send("ZINCRBY", "PRODUCT_LIKES", -1, req.ProductId)
		if err != nil {
			log.Errorf("Error decreasing like for the product with id %v", req.ProductId)
		}
		err = redisConn.Send("ZREM", utils.PRODUCT_LIKES+req.ProductId, time.Now().Unix(), req.UserName)
		if err != nil {
			log.Errorf("Error removing like for the product with id %v by user %v", req.ProductId, req.UserName)
		}
		err = redisConn.Send("ZREM", utils.PRODUCTS_LIKED_BY_USER+req.UserName, time.Now().Unix(), req.ProductId)
		if err != nil {
			log.Errorf("Error removind product with is %v in user liked list %v", req.ProductId, req.UserName)
		}
	}
	err := redisConn.Flush()
	if err != nil {
		log.Errorf("Error while flushing likes to redis %v", err)
		return err
	}
	return nil
}

func (req *ProductBookmarkRequest) BookmarkProduct() error {
	redisConn := db.GetRedisConnFromPool()
	defer redisConn.Close()

	if req.IsBookmarked {
		// Products Bookmarked By User
		err := redisConn.Send("ZADD", utils.PRODUCTS_BOOKMARKED_BY_USER+req.UserName, time.Now().Unix(), req.ProductId)
		if err != nil {
			log.Errorf("Error adding product with is %v in user bookmarked list %v", req.ProductId, req.UserName)
			redisConn.Flush()
			return err
		}
	} else {
		// Products Bookmarked By User
		err := redisConn.Send("ZREM", utils.PRODUCTS_BOOKMARKED_BY_USER+req.UserName, time.Now().Unix(), req.ProductId)
		if err != nil {
			log.Errorf("Error removing product with is %v in user bookmarked list %v", req.ProductId, req.UserName)
			redisConn.Flush()
			return err
		}
	}
	err := redisConn.Flush()
	if err != nil {
		log.Errorf("Error while flushing likes to redis %v", err)
		return err
	}
	return nil
}

func (req *FetchProductLikedByUserRequest) FetchProductsLikedByUser() (ProductLikedByUserResponse, error) {
	response := ProductLikedByUserResponse{}
	redisConn := db.GetRedisConnFromPool()
	defer redisConn.Close()

	likedProductsForUser, err := redis.Strings(redisConn.Do("ZREVRANGE", utils.PRODUCTS_LIKED_BY_USER+req.UserName, 0, -1))
	if err != nil {
		log.Errorf("Error fetching liked products for user with username %v due to error %v", req.UserName, err)
		return response, err
	}

	response.Products = FetchProductListView(likedProductsForUser)
	return response, nil
}

func (req *FetchProductBookmarkedByUser) FetchProductsBookmarkedByUser() (ProductBookmarkedByUserResponse, error) {
	response := ProductBookmarkedByUserResponse{}
	redisConn := db.GetRedisConnFromPool()
	defer redisConn.Close()

	bookmarkedProductsForUser, err := redis.Strings(redisConn.Do("ZREVRANGE", utils.PRODUCTS_BOOKMARKED_BY_USER+req.UserName, 0, -1))
	if err != nil {
		log.Errorf("Error fetching liked products for user with username %v due to error %v", req.UserName, err)
		return response, err
	}

	response.Products = FetchProductListView(bookmarkedProductsForUser)
	return response, nil
}

// TODO
func UpdateProductLikesInDB() error {
	redisConn := db.GetRedisConnFromPool()
	defer redisConn.Close()

	productsWithLikes, err := redis.Strings(redisConn.Do("ZREVRANGE", "PRODUCT_LIKES", 0, -1, "WITHSCORES"))
	if err != nil {
		log.Errorf("Error while fetching product likes from redis %v", err)
		return err
	}
	for _, val := range productsWithLikes {
		log.Infof("%v", val)
	}
	return nil
}

func FetchProductLikes(productId string) (int64, error) {
	redisConn := db.GetRedisConnFromPool()
	defer redisConn.Close()
	exists, err := redisConn.Do("ZRANK", "PRODUCT_LIKES", productId)
	if err != nil {
		log.Errorf("Error fetching product likes : %v", err)
		return 0, err
	}
	if exists == nil {
		return 0, nil
	}
	likeCount, err := redis.Int64(redisConn.Do("ZSCORE", "PRODUCT_LIKES", productId))
	if err != nil {
		log.Errorf("Error while fetching product likes from redis %v", err)
		return 0, err
	}
	return likeCount, nil
}

func ValidateProductIds(productIds []string, dbConn *sql.DB) bool {
	if len(productIds) == 0 {
		return true
	}

	productQueryString := utils.FetchInQueryStringFromArray(productIds)
	count := 0
	queryStmt := "Select count(*) As Count from Locklly.Product where ID in (" + productQueryString + ")"
	err := dbConn.QueryRow(queryStmt).Scan(&count)
	if err != nil {
		log.Infof("Error while fetching row count for products %v", err)
		return false
	}
	if count == len(productIds) {
		return true
	}
	return false
}

// Get offer for the product

// Get merchant for the product

// Add Category for the Product (Automatically adds product to the parent category)
// Calls for parent category for a category

// Get category for the product
