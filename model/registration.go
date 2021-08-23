package model

import (
	"context"
	"database/sql"
	"errors"
	"firebase.google.com/go/auth"
	"github.com/labstack/gommon/log"
	"main/db"
	"main/utils"
	"regexp"
	"strings"
	"time"
)

type RegisterMerchantRequest struct {
	AuthToken       string `json:"auth_token"`
	Email           string `json:"email"`
	DisplayName     string `json:"display_name"`
	PhotoUrl        string `json:"photo_url"`
	IsEmailVerified bool   `json:"is_email_verified"`
	//MerchantHandle string `json:"merchant_handle"`
	//BrandName      string `json:"brand_name"`
}

type UserAuthInfo struct {
	AuthToken string `json:"auth_token"`
	UserName  string `json:"user_name"`
}

type RegistrationResponse struct {
	UserAuthInfo
	RegistrationStatus   string `json:"registration_status"`
	ErrorMessage         string `json:"error_message"`
	EmailVerificationReq bool   `json:"email_verification_req"`
	NewUser              bool   `json:"new_user"`
	Email                string `json:"email"`
	DisplayName          string `json:"display_name"`
	PhotoUrl             string `json:"photo_url"`
	BrandHandle          string `json:"brand_handle"`
}

func (req *RegisterMerchantRequest) RegisterUser() (RegistrationResponse, error) {

	dbConn := db.DBConn()
	//defer dbConn.Close()

	userObj, isExisting, err := req.UserExists(dbConn)
	if err != nil {
		return RegistrationResponse{}, err
	}
	uId, err := req.GetFirebaseUser()
	if req.AuthToken != uId {
		return RegistrationResponse{}, utils.ErrUserNotRegistered
	}
	if isExisting {
		log.Infof("User already Exists")

		if err != nil {
			return RegistrationResponse{}, err
		}

		if uId == "" {
			// TODO : Need to be decided
		}

		userAuthInfo := UserAuthInfo{uId, userObj.UserName}
		response := RegistrationResponse{
			RegistrationStatus:   "Registered",
			UserAuthInfo:         userAuthInfo,
			ErrorMessage:         "",
			EmailVerificationReq: req.IsEmailVerified,
			NewUser:              false,
			Email:                req.Email,
			DisplayName:          userObj.MerchantDisplayName,
			PhotoUrl:             userObj.PhotoUrl,
			BrandHandle:          userObj.BrandHandle,
		}
		return response, nil

	} else {
		log.Infof("User can be created")
		//if uId == "" {
		//	uId, err = req.CreateFirebaseUser()
		//	if err != nil {
		//		log.Errorf("Error in Register User Function in model/registration.go while creating firebase user")
		//		return RegistrationResponse{}, err
		//	}
		//}

		userName, err := GenerateRandomUserName(dbConn)
		if err != nil {
			log.Errorf("Error while generating username for merchant %v", err)
			return RegistrationResponse{}, err
		}
		brandHandle := GenerateBranchHandleFromEmailId(req.Email, dbConn)

		err = req.CreateDBUser(userName, brandHandle, dbConn)
		if err != nil {
			log.Errorf("Error in Register User Function in model/registration.go while creating DB user")
			return RegistrationResponse{}, err
		}

		userAuthInfo := UserAuthInfo{uId, userName}
		response := RegistrationResponse{
			RegistrationStatus:   "Registered",
			UserAuthInfo:         userAuthInfo,
			ErrorMessage:         "",
			EmailVerificationReq: req.IsEmailVerified,
			NewUser:              true,
			Email:                req.Email,
			DisplayName:          req.DisplayName,
			PhotoUrl:             req.PhotoUrl,
			BrandHandle:          brandHandle,
		}
		return response, nil
	}

	return RegistrationResponse{}, nil
}

func (req *RegisterMerchantRequest) GetFirebaseUser() (string, error) {
	_, app := db.GetFirebaseClient()
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Errorf("Error creating Firebase Client")
		return "", utils.ErrServerError
	}

	if req.Email == "" {
		return "", utils.ErrBadRequest
	}

	u, err := client.GetUserByEmail(context.Background(), req.Email)
	if err != nil && strings.Contains(err.Error(), "cannot find user") {
		return "", nil
	} else if err != nil {
		log.Errorf("Error Getting user using EmailId. Error is : %v", err)
		return "", err
	}
	return u.UID, nil
}

func (req RegisterMerchantRequest) UserExists(dbConn *sql.DB) (Merchant, bool, error) {

	response := Merchant{}
	if req.Email == "" {
		return response, false, utils.ErrBadRequest
	}

	queryStmt := "Select Username, Bio, Active, PhotoUrl, BackgroundPhotoUrl, Name, Verified, CreatedAt, FollowersCount, ProductViewCount, Address, Email, Phone, BrandHandle, MerchantName from Locklly.Merchant where Email='" + strings.ToLower(req.Email) + "'"
	results, err := dbConn.Query(queryStmt)
	defer results.Close()

	if err != nil {
		log.Errorf("Error while query for Merchant in DB, %v", err)
		return response, false, err
	} else if results.Next() {
		createdAt := ""
		err := results.Scan(&response.UserName, &response.Bio, &response.IsActive, &response.PhotoUrl, &response.BackgroundPhotoUrl, &response.BrandName, &response.IsVerified, &createdAt, &response.FollowersCount, &response.ProductViewCount, &response.Address, &response.Email, &response.Phone, &response.BrandHandle, &response.MerchantDisplayName)
		response.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			log.Errorf("Error while fetching UserName from DB, err %v", err)
			return response, true, err
		}
		return response, true, nil
	}

	return response, false, nil
}

func GenerateRandomUserName(dbConn *sql.DB) (string, error) {

	username := "user" + "_" + GenerateRandomString(8)

	for {
		queryStmt := "Select * from Locklly.Merchant where Username='" + username + "'"
		result, err := dbConn.Query(queryStmt)
		if err != nil {
			log.Errorf("Error while checkout for merchant handle in Db %v", err)
		}
		if result.Next() {
			username = "user" + "_" + GenerateRandomString(8)
		} else {
			break
		}
	}
	return username, nil
}


func GenerateBranchHandleFromEmailId(email string, dbConn *sql.DB) string {

	brandHandle := ""
	brandHandlePrefix := "apparel@"
	regexFilter, _ := regexp.Compile("[^a-zA-Z0-9]+")
	possibleBrandHandle := regexFilter.Split(strings.ToLower(email), -1)[0]
	brandHandleToCheck := []string{possibleBrandHandle, possibleBrandHandle + GenerateRandomNumberString(3), possibleBrandHandle + GenerateRandomNumberString(5)}
	brandHandleToCheck = utils.RemoveAll(brandHandleToCheck, []string{""})

	brandHandle = CheckHandleAvailability(brandHandlePrefix, brandHandleToCheck, dbConn)

	if brandHandle == brandHandlePrefix {
		brandHandle = brandHandlePrefix + GenerateRandomString(5)
	}

	return brandHandle
}

func CheckHandleAvailability(brandHandlePrefix string, brandHandleToCheck []string, dbConn *sql.DB) string{

	brandHandle := brandHandlePrefix
	for _, possibleHandle := range brandHandleToCheck {
		queryStmt := "Select * from Locklly.Merchant where BrandHandle='" + brandHandle+possibleHandle + "'"
		result, err := dbConn.Query(queryStmt)
		if err != nil {
			log.Errorf("Error while checkout for merchant handle in Db %v", err)
		}
		if result.Next() {
			continue
		} else {
			brandHandle = brandHandle + possibleHandle
			return brandHandle
		}
	}
	return brandHandle
}

func (req *RegisterMerchantRequest) CreateFirebaseUser() (string, error) {
	if req.Email == "" {
		return "", utils.ErrBadRequest
	}

	_, app := db.GetFirebaseClient()
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Errorf("Error creating Firebase Client")
		return "", utils.ErrServerError
	}
	var params *auth.UserToCreate
	if req.IsEmailVerified {
		params = (&auth.UserToCreate{}).
			Email(req.Email).
			EmailVerified(true)
	} else {
		params = (&auth.UserToCreate{}).
			Email(req.Email).
			EmailVerified(false)
	}
	u, err := client.CreateUser(context.Background(), params)
	if err != nil {
		log.Errorf("Error creating Firebase User: %v\n", err)
		return "", err
	}
	log.Infof("Successfully created Firebase user: %v\n", u.UID)
	return u.UID, err
}

func (req *RegisterMerchantRequest) CreateDBUser(userName, brandHandle string, dbConn *sql.DB) error {
	log.Infof("Inside CreateDBUser method")

	loc, _ := time.LoadLocation("Asia/Kolkata")

	queryStmt := "Insert into Locklly.Merchant (Username, PhotoUrl, MerchantName, BrandHandle,  Email, CreatedAt, UpdatedAt) values(?,?,?,?,?,?,?)"

	insert1, err := dbConn.Query(queryStmt, userName, req.PhotoUrl, req.DisplayName, brandHandle, req.Email, time.Now().In(loc), time.Now().In(loc))
	defer insert1.Close()
	if err != nil {
		log.Errorf("Error while executing db query for creating db merchant user %v", err)
	}
	return err
}

func (req RegisterMerchantRequest) IsUserRegistrationPossible() (bool, error) {

	if req.Email == "" {
		return false, errors.New("Request Data is Invalid")
	}
	if req.AuthToken == "" {
		return false, utils.ErrUnauthorized
	}
	return true, nil
}
