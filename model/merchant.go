package model

import (
	"database/sql"
	"github.com/dchest/uniuri"
	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"main/db"
	"main/utils"
	"strings"
	"time"
)

var (
	CharSet = []byte("abcdefghij0123456789")
	NumSet  = []byte("0123456789")
)

type HeaderInfo struct {
	UserName  string `json:"user_name"`
	AuthToken string `json:"auth_token"`
	IsAdmin   bool   `json:"is_admin"`
}

type Merchant struct {
	//UID              string    `json:"uid"`
	UserName            string    `json:"user_name"`
	BrandName           string    `json:"brand_name"`
	CreatedAt           time.Time `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	IsVerified          bool      `json:"is_verified"`
	IsActive            bool      `json:"is_active"`
	Bio                 string    `json:"bio"`
	FollowersCount      int64     `json:"followers_count,omitempty"`
	ProductViewCount    int64     `json:"product_view_count,omitempty"`
	Phone               string    `json:"phone"`
	Email               string    `json:"email"`
	Address             string    `json:"address"`
	PhotoUrl            string    `json:"photo_url"`
	BackgroundPhotoUrl  string    `json:"background_photo_url"`
	MerchantDisplayName string    `json:"merchant_display_name"`
	BrandHandle         string    `json:"brand_handle"`
	//PANNumber        string    `json:"pan_number"`
	//AadharNumber     string    `json:"aadhar_number"`
	//GSTNumber        string    `json:"gst_number"`
	//BankAccount      string    `json:"bank_account,omitempty"`
	//IfscCode         string    `json:"ifsc_code,omitempty"`
}
type FetchMerchantInfoRequest struct {
	HeaderInfo
}

type FetchMerchantInfoResponse struct {
	Merchant
}

type FetchProductForMerchantRequest struct {
	HeaderInfo
}

type FetchProductForMerchantResponse struct {
	Products []ProductObject `json:"products"`
	Merchant `json:"merchant_info"`
}

type FetchOrderForMerchantRequest struct {
	HeaderInfo
}

type FetchOrderForMerchantResponse struct {
	Orders   []ProductOrder `json:"orders"`
	Merchant `json:"merchant_info"`
	Message  string `json:"message"`
}

type UpdateMerchantInfoRequest struct {
	HeaderInfo
	BrandName          string `json:"brand_name"`
	IsActive           bool   `json:"is_active"`
	Bio                string `json:"bio"`
	Phone              string `json:"phone"`
	Address            string `json:"address"`
	PhotoUrl           string `json:"photo_url"`
	BackgroundPhotoUrl string `json:"background_photo_url"`
	BrandHandle        string `json:"brand_handle"`
}

type VerificationDocuments struct {
	PANIdUrl    string `json:"pan_id_url"`
	AadharIdUrl string `json:"aadhar_id_url"`
	GSTIdUrl    string `json:"gst_id_url"`
	BankIdUrl   string `json:"bank_id_url"`
}

type DocumentUploadRequest struct {
	HeaderInfo
	VerificationDocuments

	PANIdStatus    string `json:"pan_id_status,omitempty"`
	AadharIdStatus string `json:"aadhar_id_status,omitempty"`
	BankIdStatus   string `json:"bank_id_status,omitempty"`
	GSTIdStatus    string `json:"gst_id_status,omitempty"`
}

func FetchHeaderInfo(ctx echo.Context) HeaderInfo {
	header := HeaderInfo{}
	header.UserName = ctx.Request().Header.Get("UserName")
	header.AuthToken = ctx.Request().Header.Get("AuthToken")
	return header
}

func GenerateRandomString(size int) string {
	return uniuri.NewLenChars(size, CharSet)
}

func GenerateRandomNumberString(size int) string {
	return uniuri.NewLenChars(size, NumSet)
}
func (req *UpdateMerchantInfoRequest) UpdateMerchantInfo(username string, dbConn *sql.DB) error {
	log.Infof("Inside UpdateMerchantInfo method")

	if req.UserName != username && req.IsAdmin == false {
		return utils.ErrBadRequest
	}
	merchantInfoRequest := FetchMerchantInfoRequest{
		req.HeaderInfo,
	}
	merchantInfoResponse, err := merchantInfoRequest.FetchMerchantDetailsByID(req.UserName, dbConn)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Asia/Kolkata")
	brandName := utils.Ternary(req.BrandName, "", merchantInfoResponse.BrandName, req.BrandName)
	bio := utils.Ternary(req.Bio, "", merchantInfoResponse.Bio, req.Bio)
	address := utils.Ternary(req.Address, "", merchantInfoResponse.Address, req.Address)
	photoUrl := utils.Ternary(req.PhotoUrl, "", merchantInfoResponse.PhotoUrl, req.PhotoUrl)
	backgroundPhotourl := utils.Ternary(req.BackgroundPhotoUrl, "", merchantInfoResponse.BackgroundPhotoUrl, req.BackgroundPhotoUrl)
	phone := utils.Ternary(req.Phone, "", merchantInfoResponse.Phone, req.Phone)
	isActiveString := utils.TernaryBool(req.IsActive)
	brandHandle := req.BrandHandle
	if brandHandle != merchantInfoResponse.BrandHandle {
		brandHandle = CheckHandleAvailability("", []string{brandHandle}, dbConn)
		if brandHandle == "" {
			return utils.ErrBrandHandleAlreadyExist
		}
	}

	updateQueryGenerator := utils.CreateUpdateFieldStatement([]string{"Name", "Bio", "Address", "PhotoUrl", "BackgroundPhotoUrl", "Phone", "Active", "UpdatedAt", "BrandHandle"}, []string{brandName, bio, address, photoUrl, backgroundPhotourl, phone, isActiveString, time.Now().In(loc).Format("2006-01-02 15:04:05"), brandHandle})
	if updateQueryGenerator == "" {
		log.Info("Unable to create Update query statement")
		return utils.ErrBadRequest
	}
	queryStmt := "Update Locklly.Merchant Set " + updateQueryGenerator + " where UserName='" + req.UserName + "'"
	update, err := dbConn.Query(queryStmt)
	defer update.Close()
	if err != nil {
		log.Errorf("Error while updating DB for merchant details inside model.UpdateMerchantInfo %v", err)
		return err
	}
	return nil
}

func (req *DocumentUploadRequest) UploadVerificationDocuments(dbConn *sql.DB) error {

	merchantInfoRequest := FetchMerchantInfoRequest{
		req.HeaderInfo,
	}
	merchantInfoResponse, err := merchantInfoRequest.FetchMerchantDetailsByID(req.UserName, dbConn)
	if err != nil {
		if strings.HasPrefix(err.Error(), "User Doesn't Exist") {
			return utils.ErrUnauthorized
		}
		return err
	}

	if !(req.PANIdStatus == "0" || req.PANIdStatus == "1") {
		req.PANIdStatus = "0"
	}
	if !(req.AadharIdStatus == "0" || req.AadharIdStatus == "1") {
		req.AadharIdStatus = "0"
	}
	if !(req.BankIdStatus == "0" || req.BankIdStatus == "1") {
		req.BankIdStatus = "0"
	}
	if !(req.GSTIdStatus == "0" || req.GSTIdStatus == "1") {
		req.GSTIdStatus = "0"
	}

	updateQueryStmt := ""
	if merchantInfoResponse.UserName == "" {
		return utils.ErrBadRequest
	} else if merchantInfoResponse.IsVerified {
		return nil
	} else {
		queryStmt := "Select Username from Locklly.Documents where Username='" + req.UserName + "'"
		result, err := dbConn.Query(queryStmt)
		if err != nil {

		} else {
			if result.Next() {
				updateQueryStmt = "Update Locklly.Documents set PANIdURL='" + req.PANIdUrl + "', PANIdVerified='" + req.PANIdStatus + "', AadharIdURL='" + req.AadharIdUrl + "', AadharIdVerified='" + req.AadharIdStatus + "', BankIdURL='" + req.BankIdUrl + "', BankIdVerified='" + req.BankIdStatus + "', GSTIdURL='" + req.GSTIdUrl + "', GSTIdVerified='" + req.GSTIdStatus + "' where Username='" + req.UserName + "'"
				update, err := dbConn.Query(updateQueryStmt)
				defer update.Close()
				if err != nil {
					log.Errorf("Error while Updating Docs in DB %v", err)
					return err
				}
			} else {
				updateQueryStmt = "Insert into Locklly.Documents (Username, PANIdURL, PANIdVerified, AadharIdURL, AadharIdVerified, BankIdURL, BankIdVerified, GSTIdURL, GSTIdVerified) values(?,?,?,?,?,?,?,?,?)"
				insert, err := dbConn.Query(updateQueryStmt, req.UserName, req.PANIdUrl, req.PANIdStatus, req.AadharIdUrl, req.AadharIdStatus, req.BankIdUrl, req.BankIdStatus, req.GSTIdUrl, req.GSTIdStatus)
				defer insert.Close()
				if err != nil {
					log.Errorf("Error while Inserting Docs in DB %v", err)
					return err
				}
			}
		}

	}

	return nil
}

func (req *FetchMerchantInfoRequest) FetchMerchantDetailsByID(username string, dbConn *sql.DB) (FetchMerchantInfoResponse, error) {
	log.Info("Inside FetchMerchantDetailsByID method")
	response := FetchMerchantInfoResponse{}

	//if req.UserName != username {
	//	log.Infof("User Unauthorized")
	//	return response, utils.ErrUnauthorized
	//}

	queryStmt := "Select Merchant.Username, Merchant.Bio, Merchant.Active, Merchant.PhotoUrl, Merchant.BackgroundPhotoUrl, Merchant.Name, Merchant.Verified, Merchant.CreatedAt, Merchant.UpdatedAt, Merchant.FollowersCount, Merchant.ProductViewCount, Merchant.Address, Merchant.Phone, Merchant.Email, Merchant.BrandHandle, Merchant.MerchantName from Locklly.Merchant where Username='" + username + "'"
	results, err := dbConn.Query(queryStmt)
	defer results.Close()
	if err != nil {
		return response, err
	} else {
		createdAt := ""
		updatedAt := ""
		if results.Next() {
			err = results.Scan(&response.UserName, &response.Bio, &response.IsActive, &response.PhotoUrl, &response.BackgroundPhotoUrl, &response.BrandName, &response.IsVerified, &createdAt, &updatedAt, &response.FollowersCount, &response.ProductViewCount, &response.Address, &response.Phone, &response.Email, &response.BrandHandle, &response.MerchantDisplayName)
			if createdAt != "" {
				response.CreatedAt, err = utils.ParseStringToTime(createdAt)
				if err != nil {
					log.Errorf("Error parsing merchant creation time %v", err)
					return FetchMerchantInfoResponse{}, err
				}
			}
			if updatedAt != "" {
				response.UpdatedAt, err = utils.ParseStringToTime(updatedAt)
				if err != nil {
					log.Errorf("Error parsing merchant updation time %v", err)
					return FetchMerchantInfoResponse{}, err
				}
			}

		} else {
			return response, utils.ErrUserDoesNotExist
		}
	}
	return response, nil
}

func (req *FetchProductForMerchantRequest) FetchProductForMerchantId(dbConn *sql.DB, merchantId string) (FetchProductForMerchantResponse, error) {
	log.Info("Inside FetchProductForMerchantId Method")
	response := FetchProductForMerchantResponse{}

	merchantInfoRequest := FetchMerchantInfoRequest{
		HeaderInfo: req.HeaderInfo,
	}
	merchantInfoResponse, err := merchantInfoRequest.FetchMerchantDetailsByID(merchantId, dbConn)
	if err != nil {
		if strings.HasPrefix(err.Error(), "User Doesn't Exist") {
			return response, utils.ErrUserNotRegistered
		}
		return response, err
	}

	response.Merchant = merchantInfoResponse.Merchant
	queryStmt := "Select Product.ID, Product.Title, Product.Description, Product.CreatedBy, Product.CreatedAt, Product.Category, Product.PhotoUrls, Product.Likes, Product.Rating, Product.Price, Product.CountPeopleRated, Product.Sizes, Product.Tags, Product.Events, Product.InterestingFact, Product.UpdatedAt, Product.State from Locklly.Product where CreatedBy='" + merchantId + "' order by CreatedAt desc"
	results, err := dbConn.Query(queryStmt)
	defer results.Close()
	if err != nil {
		log.Errorf("Error while making query for fetching products for product ids %v", err)
		return response, err
	}

	for results.Next() {
		product := ProductObject{}
		createdAt := ""
		updatedAt := ""
		photos := ""
		categories := ""
		tags := ""
		events := ""
		availableSizes := ""

		err = results.Scan(&product.Id, &product.Title, &product.Description, &product.CreatedBy, &createdAt, &categories, &photos, &product.Likes, &product.Rating, &product.Price, &product.CountPeopleRated, &availableSizes, &tags, &events, &product.InterestingFact, &updatedAt, &product.State)
		if err != nil {
			log.Errorf("Error while scanning product data from db %v", err)
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
		product.Categories = utils.FetchListFromString(categories, "##")

		product.Likes, err = FetchProductLikes(product.Id)
		if err != nil {
			log.Errorf("Unable to fetch likes for product %v", err)
		}
		response.Products = append(response.Products, product)
	}
	return response, nil

}

func (req *FetchOrderForMerchantRequest) FetchOrderForMerchantId(dbConn *sql.DB, merchantId string) (FetchOrderForMerchantResponse, error) {
	log.Info("Inside FetchProductForMerchantId Method")

	// Open and defer close db connections
	redisConn := db.GetRedisConnFromPool()
	defer redisConn.Close()

	response := FetchOrderForMerchantResponse{}

	merchantInfoRequest := FetchMerchantInfoRequest{
		HeaderInfo: req.HeaderInfo,
	}
	merchantInfoResponse, err := merchantInfoRequest.FetchMerchantDetailsByID(merchantId, dbConn)
	if err != nil {
		if strings.HasPrefix(err.Error(), "User Doesn't Exist") {
			return response, utils.ErrUserNotRegistered
		}
		return response, err
	}
	response.Merchant = merchantInfoResponse.Merchant

	orderIds, err := redis.Strings((redisConn).Do("ZRANGE", utils.ORDERS_FOR_MERCHANT+merchantId, 0, -1))
	if err != nil {
		return FetchOrderForMerchantResponse{Message: "Error!"}, err
	}
	if len(orderIds) == 0 {
		log.Infof("No Product Order for the Merchant")
		response.Message = "No Product Order for the Merchant"
		return response, nil
	}

	orderDetails, err := FetchOrdersByOrderIds(dbConn, orderIds)

	for _, orderId := range orderIds {
		if orderDetail, ok := orderDetails[orderId]; ok {
			response.Orders = append(response.Orders, orderDetail)
		} else {
			log.Errorf("Unable to find order details in the map %v", err)
			return FetchOrderForMerchantResponse{}, err
		}
	}
	response.Message = "Success!"
	return response, nil

}
