package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2"
	"os"
)

func Database() (*sql.DB, error) {
	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", os.Getenv("MYSQL_USERNAME"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_DB")))
	err := db.Ping()
	return db, err
}


func MongoDatabase() (*mgo.Database, error) {
	mongo, err := mgo.Dial(os.Getenv("MONGO_URL"))
	if err != nil {
		return nil, err
	}
	mongodb := mongo.DB(os.Getenv("MONGO_DB"))
	return mongodb, err
}
