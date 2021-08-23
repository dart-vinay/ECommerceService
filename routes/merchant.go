package routes

import (
	"encoding/json"
	"github.com/labstack/echo"
	"main/db"
	"main/model"
	"main/utils"
	"net/http"
)

//Update Merchant Primary Credentials


// Update Merchant Other Details

func UpdateMerchantInfo(c echo.Context) error{
	merchantId := c.Param("id")
	header := model.FetchHeaderInfo(c)

	req := model.UpdateMerchantInfoRequest{}
	req.HeaderInfo = header
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "JSON values unreadable")
	}
	if req.UserName == "" {
		return c.JSON(http.StatusBadRequest, "AuthToken is Empty")
	}

	dbConn := db.DBConn()
	//defer dbConn.Close()
	if err := req.UpdateMerchantInfo(merchantId, dbConn); err != nil {
		return handleErrLogic(err, c)
	} else {
		return c.JSON(http.StatusOK, model.EmptyResponse{})
	}

}


func FetchMerchantInfoByID(c echo.Context) error{
	merchantId := c.Param("id")
	header := model.FetchHeaderInfo(c)
	req := &model.FetchMerchantInfoRequest{}
	req.HeaderInfo = header
	if merchantId == "" {
		return c.JSON(http.StatusBadRequest, "Empty MerchantID")
	}

	// open and defer close db connection
	dbConn := db.DBConn()
	//defer dbConn.Close()

	if response, err := req.FetchMerchantDetailsByID(merchantId, dbConn); err != nil {
		return handleErrLogic(err, c)
	} else {
		return c.JSON(http.StatusOK, response)
	}

}


//
//func VerifyMerchant(c echo.Context) error{
//	merchantId := c.Param("id")
//	if merchantId==""{
//		return c.JSON(http.StatusBadRequest, "Empty MerchantID")
//	}
//	header := model.FetchHeaderInfo(c)
//	req := &model.VerifyMerchantRequest{}
//	req.HeaderInfo = header
//
//	if err:= req.VerifyMerchant(merchantId); err!=nil{
//		return c.JSON(http.StatusInternalServerError, nil)
//	}else{
//		return c.JSON(http.StatusOK, model.EmptyResponse{})
//	}
//}
//


func FetchProductForMerchantId(c echo.Context) error{
	merchantId := c.Param("id")
	if merchantId==""{
		return c.JSON(http.StatusBadRequest, "Empty MerchantID")
	}
	header := model.FetchHeaderInfo(c)
	req := &model.FetchProductForMerchantRequest{}
	req.HeaderInfo = header

	// open and defer close db connection
	dbConn := db.DBConn()
	//defer dbConn.Close()

	if response, err:= req.FetchProductForMerchantId(dbConn, merchantId); err!=nil{
		return handleErrLogic(err, c)
	}else{
		return c.JSON(http.StatusOK, response)
	}
}

func FetchOrderForMerchantId(c echo.Context) error{
	merchantId := c.Param("id")
	if merchantId==""{
		return c.JSON(http.StatusBadRequest, "Empty MerchantID")
	}
	header := model.FetchHeaderInfo(c)
	req := &model.FetchOrderForMerchantRequest{}
	req.HeaderInfo = header

	// open and defer close db connection
	dbConn := db.DBConn()
	//defer dbConn.Close()

	if response, err:= req.FetchOrderForMerchantId(dbConn, merchantId); err!=nil{
		return handleErrLogic(err, c)
	}else{
		return c.JSON(http.StatusOK, response)
	}
}

func UploadVerificationDocuments(c echo.Context) error{
	header := model.FetchHeaderInfo(c)
	req := model.DocumentUploadRequest{}
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "JSON values unreadable")
	}
	req.HeaderInfo = header

	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
		return err
	}


	// Open and defer close DB connection
	dbConn := db.DBConn()
	//defer dbConn.Close()

	if err := req.UploadVerificationDocuments(dbConn); err != nil {
		return handleErrLogic(err, c)
	} else {
		return c.JSON(http.StatusOK, model.EmptyResponse{})
	}

}