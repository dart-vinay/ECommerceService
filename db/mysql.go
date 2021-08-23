package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbDriver = "mysql"
	dbUser   = "newuser"
	dbPasswd = "root"
	dbName   = "Locklly"
	dbGlobal *sql.DB
)

//type SQLDBConnectionDetails struct{
//	Session *sql.Conn
//}
//func InitDBConn() {
//	//db, err := sql.Open(dbDriver, dbUser+":"+dbPasswd+"@tcp(docker.for.mac.localhost:3306)/"+dbName) // While running on docker use this
//	db, err := sql.Open(dbDriver, dbUser+":"+dbPasswd+"@tcp(localhost:3306)/"+dbName)
//	if err != nil {
//		panic(err)
//	}
//	dbGlobal = db
//}

func InitDBConnection() {
	//db, err := sql.Open(dbDriver, dbUser+":"+dbPasswd+"@tcp(docker.for.mac.localhost:3306)/"+dbName) // While running on docker use this
	db, err := sql.Open(dbDriver, dbUser+":"+dbPasswd+"@/"+dbName)
	db.SetMaxOpenConns(500)
	if err != nil {
		panic(err)
	}
	dbGlobal = db
}
func DBConn() *sql.DB {
	return dbGlobal
}

//func DBConn() *sql.DB {
//	//db, err := sql.Open(dbDriver, dbUser+":"+dbPasswd+"@tcp(docker.for.mac.localhost:3306)/"+dbName) // While running on docker use this
//	db, err := sql.Open(dbDriver, dbUser+":"+dbPasswd+"@/"+dbName)
//	db.SetMaxOpenConns(500)
//	if err != nil {
//		panic(err)
//	}
//	return db
//}

//func (this *SQLDBConnectionDetails) Init() *sql.Conn{
//	this.Session, _ = dbGlobal.Conn(context.Background())
//	return this.Session
//}
