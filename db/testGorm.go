package db

import (
	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
	NewField string
	CreatedAt time.Time `gorm:"type:time"`
	//UpdatedAt time.Time `gorm:"type:time"`
	//DeletedAt time.Time `gorm:"type:time"`
	TimeField time.Time `gorm:"type:time,omitempty"`

}

func Test() {
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8&parseTime=True&loc=Local"
	//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: DBConn(),
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	err1 := db.Create(&Product{Code: "D4444", Price: 100, NewField: "Yeay"})

	log.Infof("%v",err1)
	// Read
	var product Product
	db.First(&product, 1)                 // find product with integer primary key
	db.First(&product, "code = ?", "F42") // find product with code D42

	// Update - update product's price to 200
	//db.Model(&product).Update("Price", 200)
	//// Update - update multiple fields
	//db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	//db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	//// Delete - delete product
	//db.Delete(&product, 1)
}
