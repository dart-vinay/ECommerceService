package model

import (
	"database/sql"
	"github.com/labstack/gommon/log"
	"main/db"
	"main/utils"
	"strconv"
	"strings"
	"time"
)

type ProductOrder struct {
	OrderId         string `json:"order_id"`
	OrderBundle     []OrderUnit
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	MerchantId      string    `json:"merchant_id"`
	DeliveryAddress string    `json:"delivery_address"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	Pincode         string    `json:"pincode"`
	Status          string    `json:"status"`
	PaymentMode     string    `json:"payment_mode"`
	PaymentStatus   string    `json:"payment_status"`
	Description     string    `json:"description"`
}

type OrderUnit struct {
	ProductId   string `json:"product_id"`
	Quantity    int    `json:"quantity"`
	Size        string `json:"size"`
	PhotoUrl    string `json:"photo_url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type OrderHistory struct {
	ID           string    `json:"id"`
	OrderId      string    `json:"order_id"`
	Status       string    `json:"status"`
	Descriptions string    `json:"descriptions"`
	DeliveryDate time.Time `json:"delivery_date"`
	ModifiedBy   string    `json:"modified_by"`
	CreatedAt    time.Time `json:"created_at"`
}

// Order Status
const (
	PENDING_PAYMENT    = "PENDING_PAYMENT"
	PLACED             = "PLACED"
	ACCEPTED           = "ACCEPTED"
	REJECTED_BY_SELLER = "REJECTED_BY_SELLER"
	PROCESSED          = "IN_PROCESS"
	SHIPPED            = "SHIPPED"
	DELIVERED          = "DELIVERED"
	USER_CANCELLED     = "USER_CANCELLED"
)

type CreateProductOrderRequest struct {
	HeaderInfo
	OrderBundle     []OrderUnit `json:"order_bundle"`
	MerchantId      string      `json:"merchant_id"`
	DeliveryAddress string      `json:"delivery_address"`
	City            string      `json:"city"`
	State           string      `json:"state"`
	Status          string      `json:"status"`
	PaymentMode     string      `json:"payment_mode"`
	PaymentStatus   string      `json:"payment_status"`
	Description     string      `json:"description"`
	Pincode         string      `json:"pincode"`
}

type FetchOrderDetailsRequest struct {
	HeaderInfo
}

type FetchOrderDetailsReponse struct {
	ProductOrder
}

type UpdateOrderInfoRequest struct {
	HeaderInfo
	Status            string `json:"status"`
	Description       string `json:"description"`
	EstimatedDelivery string `json:"estimated_delivery"`
}

// TODO : Need to Specify Order Status as Per Payment Mode and status. Need to Verify payments first.
func (req *CreateProductOrderRequest) CreateProductOrder(dbConn *sql.DB) error {

	redisConn := db.GetRedisConnFromPool()
	defer redisConn.Close()

	if len(req.OrderBundle) == 0 || req.MerchantId == "" || req.DeliveryAddress == "" {
		return utils.ErrBadRequest
	}

	orderId := GenerateRandomString(8)

	// TODO: Think of a better struct to store order details
	orderBundleString := ""
	productIds := []string{}
	for _, orderUnit := range req.OrderBundle {
		orderUnitString := ""
		if orderUnit.Quantity == 0 || orderUnit.Size == "" || !utils.ValidateSize(strings.ToUpper(orderUnit.Size)) {
			log.Error("Quantity can't be zero. Size can't be empty ")
			return utils.ErrBadRequest
		}
		productIds = append(productIds, orderUnit.ProductId)
		orderUnitString = orderUnit.ProductId + "#" + strconv.Itoa(orderUnit.Quantity) + "#" + strings.ToUpper(orderUnit.Size)
		orderBundleString += orderUnitString + "##"
	}
	orderBundleString = strings.Trim(orderBundleString, "##")

	// TODO : Think of a better product validation scheme
	validateProductId := ValidateProductIds(productIds, dbConn)
	if !validateProductId {
		return utils.ErrBadRequest
	}

	queryStmt := "Insert into Locklly.Order (OrderId, OrderBundle,  CreatedAt, UpdatedAt, MerchantId, DeliveryAddress, City, State, Status, PaymentMode, PaymentStatus, Description, Pincode) values(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err := dbConn.Query(queryStmt, orderId, orderBundleString, time.Now(), time.Now(), req.MerchantId, req.DeliveryAddress, req.City, req.State, PLACED, req.PaymentMode, req.PaymentStatus, req.Description, req.Pincode)
	if err != nil {
		log.Errorf("Error while inserting order in orders table %v", err)
		return err
	}
	err = (redisConn).Send("ZADD", utils.ORDERS_FOR_MERCHANT+req.MerchantId, time.Now().Unix(), orderId)

	for _, productId := range productIds {
		err := redisConn.Send("ZINCRBY", utils.ORDER_COUNT, 1, productId)
		if err != nil {
			log.Errorf("Exception Caught")
			return err
		}
	}
	err = redisConn.Flush()
	if err != nil {
		log.Errorf("Error while flushing to Redis %v", err)
	}
	return nil
}

func (req FetchOrderDetailsRequest) FetchOrderDetailById(dbConn *sql.DB, orderId string) (FetchOrderDetailsReponse, error) {
	log.Info("Inside FetchOrderDetailById method")
	response := FetchOrderDetailsReponse{}
	if orderInfo, err := FetchOrderById(dbConn, orderId); err != nil {
		log.Errorf("Error while trying to fetch order info for order id : %v", orderId)
		return response, err
	} else {
		response.ProductOrder = orderInfo
	}
	return response, nil

}

func FetchOrderById(dbConn *sql.DB, orderId string) (ProductOrder, error) {
	response, err := FetchOrdersByOrderIds(dbConn, []string{orderId})
	if err != nil {
		log.Errorf("Error while fetching Order info")
		return ProductOrder{}, err
	} else {
		return response[orderId], nil
	}
}

func FetchOrdersByOrderIds(dbConn *sql.DB, orderIds []string) (map[string]ProductOrder, error) {
	log.Info("Inside FetchProductsByProductIds method")

	response := make(map[string]ProductOrder)

	orderQueryString := utils.FetchInQueryStringFromArray(orderIds)

	queryStmt := "Select OrderId, OrderBundle, CreatedAt, UpdatedAt, MerchantId, DeliveryAddress, City, State, Status, PaymentMode, PaymentStatus, Description, Pincode from Locklly.Order where OrderId in (" + orderQueryString + ")"
	results, err := dbConn.Query(queryStmt)
	if err != nil {
		log.Errorf("Error while making query for fetching order for order ids %v", err)
		return response, err
	}
	for results.Next() {
		order := ProductOrder{}
		createdAt := ""
		updatedAt := ""
		orderBundleString := ""

		err = results.Scan(&order.OrderId, &orderBundleString, &createdAt, &updatedAt, &order.MerchantId, &order.DeliveryAddress, &order.City, &order.State, &order.Status, &order.PaymentMode, &order.PaymentStatus, &order.Description, &order.Pincode)
		if err != nil {
			log.Errorf("Error while scanning order data from db %v", err)
		}
		order.CreatedAt, err = utils.ParseStringToTime(createdAt)
		if err != nil {
			log.Errorf("Error parsing time for createdAt or updatedAt, %v", err)
			return response, err
		}
		order.UpdatedAt, err = utils.ParseStringToTime(updatedAt)
		if err != nil {
			log.Errorf("Error parsing time for createdAt or updatedAt, %v", err)
			return response, err
		}
		orderBundle, productIds := ConvertOrderBundleStringToOrderBundle(orderBundleString)
		productBasicInfo, err := FetchProductsByProductIds(productIds, true)
		if err != nil {
			return response, err
		}
		for index, _ := range orderBundle {
			productId := orderBundle[index].ProductId
			orderBundle[index].Title = productBasicInfo[productId].Title
			orderBundle[index].PhotoUrl = productBasicInfo[productId].PrimaryURL
			orderBundle[index].Description = productBasicInfo[productId].Description
		}
		order.OrderBundle = orderBundle
		response[order.OrderId] = order
	}
	return response, nil
}

func (req *UpdateOrderInfoRequest) UpdateOrderInfo(dbConn *sql.DB, orderId string) error {

	redisConn := db.GetRedisConnFromPool()
	defer redisConn.Close()

	if req.Status != "" {
		if !utils.ValidateStatus(strings.ToUpper(req.Status)) {
			log.Infof("Invalid Status")
			return utils.ErrBadRequest
		}
	} else {
		log.Infof("Status Cannot be empty ")
		return utils.ErrBadRequest
	}
	var err error
	deliveryDate := time.Time{}
	if req.EstimatedDelivery != "" {
		deliveryDate, err = time.Parse("2006-01-02", req.EstimatedDelivery)
		if err != nil {
			log.Errorf("Unable to parse delivery date %v", err)
			return err
		}
	} else {
		return utils.ErrBadRequest
	}
	val, err := (redisConn).Do("ZRANK", utils.ORDERS_FOR_MERCHANT+req.UserName, orderId)
	if err != nil {
		log.Errorf("Exception Caught %v", err)
		return err
	}
	if val == nil {
		log.Infof("User %v unauthorized to modify order info for orderid %v", req.UserName, orderId)
		return utils.ErrUnauthorized
	}
	queryStmt := "Insert into Locklly.OrderHistory (OrderId, Status,  Description, CreatedAt, CreatedBy, DeliveryDate) values(?,?,?,?,?,?)"
	insert, err := dbConn.Query(queryStmt, orderId, strings.ToUpper(req.Status), req.Description, time.Now(), req.UserName, deliveryDate)
	defer insert.Close()

	queryStmt = "Update Locklly.Order set Status=? where OrderId=?"
	update, err := dbConn.Query(queryStmt, req.Status, orderId)
	defer update.Close()
	return nil
}

func ConvertOrderBundleStringToOrderBundle(orderBundleString string) ([]OrderUnit, []string) {
	res := []OrderUnit{}
	productIds := []string{}
	orderUnitString := strings.Split(orderBundleString, "##")
	for _, orderUnit := range orderUnitString {
		orderEntityArray := strings.Split(orderUnit, "#")
		quantity, err := strconv.Atoi(orderEntityArray[1])
		if err != nil {
			log.Errorf("Error while parsing Order Bundle String %v", err)
			return []OrderUnit{}, []string{}
		}
		orderUnitObj := OrderUnit{ProductId: orderEntityArray[0], Quantity: quantity, Size: orderEntityArray[2]}
		res = append(res, orderUnitObj)
		productIds = append(productIds, orderEntityArray[0])
	}
	return res, productIds
}
