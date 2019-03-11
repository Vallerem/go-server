package shared

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Opening a database and save the reference to `Database` struct.
func Init() *gorm.DB {

	envErr := godotenv.Load()
	if envErr != nil {
		log.Printf("Error loading .env file")
	}

	//open a db connection
	dbUrl := os.Getenv("DATABASE_URL")
	db, err := gorm.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println("db err: ", err)
	}

	db.DB().SetMaxIdleConns(10)
	db.LogMode(true)
	DB = db
	return DB
}

// This function will create a temporarily database for running testing cases
// func TestDBInit() *gorm.DB {
// test_db, err := gorm.Open("sqlite3", "./../gorm_test.db")
// if err != nil {
// 	fmt.Println("db err: ", err)
// }
// test_db.DB().SetMaxIdleConns(3)
// test_db.LogMode(true)
// DB = test_db
// return DB
// }

// Delete the database after running testing cases.
// func TestDBFree(test_db *gorm.DB) error {
// 	test_db.Close()
// 	err := os.Remove("./../gorm_test.db")
// 	return err
// }

// Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	return DB
}
