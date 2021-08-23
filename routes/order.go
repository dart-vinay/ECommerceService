package routes

import (
	"encoding/json"
	"github.com/labstack/echo"
	"main/db"
	"main/model"
	"main/utils"
	"net/http"
)

func CreateProductOrder(c echo.Context) error {
	header := model.FetchHeaderInfo(c)
	req := model.CreateProductOrderRequest{}
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "JSON values unreadable")
	}
	req.HeaderInfo = header

	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
		return err
	}
	if len(req.OrderBundle) == 0 {
		return c.JSON(http.StatusBadRequest, "Empty Product Order")
	}

	dbConn := db.DBConn()
	//defer dbConn.Close()

	if err := req.CreateProductOrder(dbConn); err != nil {
		return handleErrLogic(err, c)
	} else {
		return c.JSON(http.StatusOK, model.EmptyResponse{})
	}

}

func FetchOrderByID(c echo.Context) error {
	orderId := c.Param("id")
	header := model.FetchHeaderInfo(c)
	req := model.FetchOrderDetailsRequest{}
	req.HeaderInfo = header

	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
		return err
	}

	dbConn := db.DBConn()
	//defer dbConn.Close()
	if response, err := req.FetchOrderDetailById(dbConn, orderId); err != nil {
		return handleErrLogic(err, c)
	} else {
		return c.JSON(http.StatusOK, response)
	}
}

func UpdateOrderInfo(c echo.Context) error {
	orderId := c.Param("id")
	header := model.FetchHeaderInfo(c)
	req := model.UpdateOrderInfoRequest{}
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "JSON values unreadable")
	}
	req.HeaderInfo = header

	// open and defer close db connections
	dbConn := db.DBConn()
	//defer dbConn.Close()

	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
		return err
	}

	if err := req.UpdateOrderInfo(dbConn, orderId); err != nil {
		return handleErrLogic(err, c)
	} else {
		return c.JSON(http.StatusOK, model.EmptyResponse{})
	}
}
