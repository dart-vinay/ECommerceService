package routes

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"main/db"
	"main/model"
	"main/utils"
	"net/http"
)

func CreateProductListing(c echo.Context) error {
	header := model.FetchHeaderInfo(c)
	req := model.CreateProductListingRequest{}
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "JSON values unreadable")
	}
	req.HeaderInfo = header

	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
		return err
	}
	if len(req.PhotoUrls) == 0 {
		return c.JSON(http.StatusBadRequest, "Photos Missing")
	}

	if err := req.CreateProductListing(); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		return c.JSON(http.StatusOK, model.EmptyResponse{})
	}

}

func FetchProductByID(c echo.Context) error {
	productId := c.Param("id")
	header := model.FetchHeaderInfo(c)
	req := model.FetchProductDetailsRequest{}
	req.HeaderInfo = header

	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
		return err
	}

	if response, err := req.FetchProductDetailById(productId); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		return c.JSON(http.StatusOK, response)
	}
}

//func FetchAllProducts(c echo.Context) error {
//	header := model.FetchHeaderInfo(c)
//	req := model.FetchProductDetailsRequest{}
//	req.HeaderInfo = header
//
//	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
//		return err
//	}
//
//	if response, err := req.FetchAllProducts(); err != nil {
//		return handleErrLogic(err, c)
//	} else {
//		return c.JSON(http.StatusOK, response)
//	}
//}

func UpdateProductInfo(c echo.Context) error {
	productId := c.Param("id")
	header := model.FetchHeaderInfo(c)
	req := model.UpdateProductInfoRequest{}
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "JSON values unreadable")
	}
	req.HeaderInfo = header

	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
		return err
	}
	if err := req.UpdateProductInfo(productId); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		return c.JSON(http.StatusOK, model.EmptyResponse{})
	}
}

func DeleteProductByID(c echo.Context) error {
	productId := c.Param("id")
	header := model.FetchHeaderInfo(c)
	req := model.UpdateProductInfoRequest{}
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "JSON values unreadable")
	}
	req.HeaderInfo = header

	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
		return err
	}
	if err := req.Delete(productId); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		return c.JSON(http.StatusOK, model.EmptyResponse{})
	}
}

func TestMethod(c echo.Context) error {

	dbConn := db.DBConn()
	//defer dbConn.Close()

	_, err := dbConn.Query("Insert into Locklly.TestTable values('1', '2')")
	if err != nil {
		log.Info("Error")
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		log.Info("Success")
		return c.HTML(http.StatusOK, "Successful execution!")
	}
}

//func LikeProduct(c echo.Context) error{
//	header := model.FetchHeaderInfo(c)
//	req := model.ProductLikeRequest{}
//	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
//		return c.JSON(http.StatusBadRequest, "JSON values unreadable")
//	}
//	req.HeaderInfo = header
//	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
//		return err
//	}
//	if err := req.LikeProduct(); err != nil {
//		return c.JSON(http.StatusInternalServerError, err)
//	} else {
//		return c.JSON(http.StatusOK, model.EmptyResponse{})
//	}
//
//}
//
//func BookmarkProduct(c echo.Context) error{
//	header := model.FetchHeaderInfo(c)
//	req := model.ProductBookmarkRequest{}
//	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
//		return c.JSON(http.StatusBadRequest, "JSON values unreadable")
//	}
//	req.HeaderInfo = header
//	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
//		return err
//	}
//	if err := req.BookmarkProduct(); err != nil {
//		return c.JSON(http.StatusInternalServerError, err)
//	} else {
//		return c.JSON(http.StatusOK, model.EmptyResponse{})
//	}
//}
//
//func FetchProductLikedByUser(c echo.Context) error{
//	header := model.FetchHeaderInfo(c)
//	req := model.FetchProductLikedByUserRequest{}
//	req.HeaderInfo = header
//	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
//		return err
//	}
//	if response,err := req.FetchProductsLikedByUser(); err != nil {
//		return c.JSON(http.StatusInternalServerError, err)
//	} else {
//		return c.JSON(http.StatusOK, response)
//	}
//}
//
//func FetchProductBookmarkedByUser(c echo.Context) error{
//	header := model.FetchHeaderInfo(c)
//	req := model.FetchProductBookmarkedByUser{}
//	req.HeaderInfo = header
//	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
//		return err
//	}
//	if response,err := req.FetchProductsBookmarkedByUser(); err != nil {
//		return c.JSON(http.StatusInternalServerError, err)
//	} else {
//		return c.JSON(http.StatusOK, response)
//	}
//}
//
//func UpdateLikesInDB(c echo.Context) error{
//	go model.UpdateProductLikesInDB()
//	return c.JSON(http.StatusOK, model.EmptyResponse{})
//}
