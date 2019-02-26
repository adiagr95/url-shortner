package main

import (
	"./ctrl"
	"./database"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}


	app := gin.Default()
	app.Use(static.Serve("/assets", static.LocalFile("assets", true)))

	db, err := database.Database()
	if err != nil {
		panic(err)
	}
	mongo, err := database.MongoDatabase()
	if err != nil {
		panic(err)
	}
	app.Use(database.Inject(*db))
	app.Use(database.InjectMongo(mongo))

	database.Migrate(db);
	app.LoadHTMLGlob("templates/*")
	initializeRoutes(app)

	app.Run()
}


func initializeRoutes(app *gin.Engine)  {

	app.GET("/:code", ctrl.ResolveUrl)
	app.POST("/create", ctrl.CreateUrl)

}