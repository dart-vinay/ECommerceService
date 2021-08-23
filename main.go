package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"main/db"
	"main/routes"
	"main/utils"
	"net/http"
)

func main() {

	fmt.Println("Starting of the Customer Service")
	e := echo.New()
	//e.HTTPErrorHandler = customErrorHandler
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.DisableHTTP2 = true

	//Initialize Redis Connection Pool
	db.InitRedisPool()
	utils.InitConstants()
	db.InitDBConnection()

	//CheckRegexFunc()
	// Registration API
	e.POST("/registration/signup", routes.RegisterUser)

	// Merchant API
	e.GET("/merchant/:id", routes.FetchMerchantInfoByID, validateUserInfo)
	e.PUT("/merchant/:id", routes.UpdateMerchantInfo, validateUserInfo)
	e.POST("/merchant/verification", routes.UploadVerificationDocuments, validateUserInfo)
	e.GET("/merchantProducts/:id", routes.FetchProductForMerchantId, validateUserInfo)
	e.GET("/merchantOrders/:id", routes.FetchOrderForMerchantId, validateUserInfo)


	//Product API
	e.POST("/product", routes.CreateProductListing, validateUserInfo)
	e.GET("/product/:id", routes.FetchProductByID, validateUserInfo)
	e.PUT("/product/:id", routes.UpdateProductInfo, validateUserInfo)
	e.DELETE("/product/:id", routes.DeleteProductByID, validateUserInfo)
	//e.GET("/products/all", routes.FetchAllProducts, validateUserInfo)

	// Category API
	e.POST("/category", routes.CreateCategory, validateUserAuthInfo)
	e.GET("/category/all", routes.FetchAllCategory, validateUserAuthInfo)
	e.GET("/category/:id", routes.FetchCategoryByCategoryId, validateUserAuthInfo)
	e.PUT("/category/:id", routes.UpdateCategoryInfo, validateUserAuthInfo)
	e.GET("/categoryProducts/:id", routes.FetchProductsForCategoryId, validateUserInfo)


	//Orders API
	e.POST("/order", routes.CreateProductOrder, validateUserInfo)
	e.GET("/order/:id", routes.FetchOrderByID, validateUserInfo)
	e.PUT("/order/:id", routes.UpdateOrderInfo, validateUserInfo)
	// Update Order

	e.GET("/testMethod", routes.TestMethod)
	e.GET("/test", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "The application is up!")
	})
	e.Logger.Fatal(e.Start(":2000"))
}

func customErrorHandler(err error, ctx echo.Context) {
	log.Errorf("ECHO ERROR: Inside Custom Error Handle Method with error %v", err)
}

func validateUserInfo(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("UserName") == "" || c.Request().Header.Get("AuthToken") == "" {
			return c.JSON(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)
	}
}
func validateUserAuthInfo(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("UserName") != "Internal" {
			return c.JSON(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)
	}
}

//func CheckRegexFunc(){
//	stringToVerify := []string{"vka0797@gmail.com", "vinay.agarwal@gmail.com", "1223@gmail.com"}
//
//	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
//	for _, name := range stringToVerify{
//		checkUserNames := reg.Split(strings.ToLower(name), -1)
//		for _, val := range checkUserNames{
//			log.Infof("%v %v", name, val)
//		}
//	}
//}

// tar -C /usr/local -xzf go1.15.8.linux-amd64.tar.gz
// export PATH=$PATH:/usr/local/go/bin