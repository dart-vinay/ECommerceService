package routes

import (
	"encoding/json"
	"main/utils"
	"github.com/labstack/echo"
	"main/model"
	"net/http"
)

//Register Merchnant

func RegisterUser(c echo.Context) error {

	req := model.RegisterMerchantRequest{}
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "JSON values unreadable")
	}
	header := model.FetchHeaderInfo(c)
	req.AuthToken = header.AuthToken
	isPossible, errResponse := req.IsUserRegistrationPossible()
	if !isPossible {
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	if response, err := req.RegisterUser(); err!=nil{
		return handleErrLogic(err, c)
	}else{
		return c.JSON(http.StatusOK, response)
	}
}

func handleErrLogic(err error, c echo.Context) error {
	if err != nil {
		if err == utils.ErrBadRequest {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		} else if err == utils.ErrServerError {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else if err == utils.ErrUnauthorized {
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		} else if err == utils.ErrUserNotRegistered{
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		} else if err == utils.ErrBrandHandleAlreadyExist {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}