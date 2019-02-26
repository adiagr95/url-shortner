package database

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2"
)

func Inject(db sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}


func InjectMongo(mongo *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("mongo", mongo)
		c.Next()
	}
}
