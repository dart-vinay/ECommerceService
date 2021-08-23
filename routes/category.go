package routes

import (
	"main/model"
	"main/utils"
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
)

func CreateCategory(c echo.Context) error {
	header := model.FetchHeaderInfo(c)
	req := model.CreateCategoryListingRequest{}
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "JSON values unreadable")
	}
	req.HeaderInfo = header

	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
		return err
	}
	if err := req.CreateCategory(); err != nil {
		return handleErrLogic(err, c)
	} else {
		return c.JSON(http.StatusOK, model.EmptyResponse{})
	}
}

func FetchAllCategory(c echo.Context) error {
	header := model.FetchHeaderInfo(c)
	req := model.FetchAllCategoriesRequest{}
	req.HeaderInfo = header

	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
		return handleErrLogic(err,c)
	}
	if response, err := req.FetchAllCategory(); err != nil {
		return handleErrLogic(err, c)
	} else {
		return c.JSON(http.StatusOK, response)
	}
}

func FetchCategoryByCategoryId(c echo.Context) error {
	categoryId := c.Param("id")
	header := model.FetchHeaderInfo(c)
	req := model.FetchCategoryRequest{}
	req.HeaderInfo = header

	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
		return err
	}
	if response, err := req.FetchCategoryByCategoryId(categoryId); err != nil {
		return handleErrLogic(err, c)
	} else {
		return c.JSON(http.StatusOK, response)
	}
}

func UpdateCategoryInfo(c echo.Context) error {
	categoryId := c.Param("id")
	header := model.FetchHeaderInfo(c)
	req := model.UpdateCategoryInfoRequest{}
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "JSON values unreadable")
	}
	req.HeaderInfo = header

	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
		return err
	}
	if err := req.UpdateCategoryInfo(categoryId); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		return c.JSON(http.StatusOK, model.EmptyResponse{})
	}
}

func FetchProductsForCategoryId(c echo.Context) error{
	categoryId := c.Param("id")
	header := model.FetchHeaderInfo(c)
	req := model.FetchProductsForCategoryRequest{}
	req.HeaderInfo = header

	if err := utils.VerifyRequestCredentials(req.UserName, req.AuthToken); err != nil {
		return err
	}
	if response, err := req.FetchProductsForCategory(categoryId); err != nil {
		return handleErrLogic(err, c)
	} else {
		return c.JSON(http.StatusOK, response)
	}
}
