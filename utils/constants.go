package utils

//
//const (
//	ADMIN    = "ADMIN"
//	MERCHANT = "MERCHANT"
//	CONSUMER = "CONSUMER"
//)

const (
	MALE   = "MALE"
	FEMALE = "FEMALE"
)

const (
	CATEGORY                    = "CATEGORY:"
	EVENT                       = "EVENT:"
	TAG                         = "TAG:"
	PRODUCT_LIKES               = "PRODUCT_LIKE:"
	PRODUCTS_LIKED_BY_USER      = "PRODUCTS_LIKED_BY_USER:"
	PRODUCTS_BOOKMARKED_BY_USER = "PRODUCTS_BOOKMARKED_BY_USER:"
	ALL_TAGS                    = "ALL_TAGS"
	ALL_EVENTS                  = "ALL_EVENTS"
	ORDERS_FOR_MERCHANT         = "orders_for_merchant:"
	ORDER_COUNT                 = "order_count"
)

type SIZE string

const (
	XS  = SIZE("XS")
	S   = SIZE("S")
	M   = SIZE("M")
	L   = SIZE("L")
	XL  = SIZE("XL")
	XXL = SIZE("XXL")
)

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

var categoryList []string
var sizeList []SIZE
var statusList []string

func InitConstants() {
	categoryList = []string{"MEN", "WOMEN", "KIDS", "TROUSER"}
	sizeList = []SIZE{XS, S, M, L, XL, XXL}
	statusList = []string{PENDING_PAYMENT, PLACED, ACCEPTED, REJECTED_BY_SELLER, PROCESSED, SHIPPED, DELIVERED, USER_CANCELLED}
}
