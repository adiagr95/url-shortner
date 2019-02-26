package ctrl

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"net/http"
)

func ResolveUrl(c *gin.Context)  {
	code := c.Param("code")
	db := c.MustGet("db").(sql.DB)
	mongo := c.MustGet("mongo").(*mgo.Database)

	insertUrlFetchInfo(c, mongo, code)
	url, err := GetUrlFromCode(db, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" :string(err.Error())})
		return
	}

	c.HTML(http.StatusOK, "redirect.html" ,gin.H{"url" : url})

}

func GetUrlFromCode(db sql.DB, code string) (string, error) {
	codeNumber := GetNumberFromCode(code)
	var tableSuffix, url string
	rows, err := db.Query(fmt.Sprintf(`SELECT table_suffix FROM mapping_master where start<=%f and end>=%f;`, codeNumber, codeNumber))
	if err != nil {
		return url, err
	}
	for rows.Next() {
		err := rows.Scan(&tableSuffix)
		if err != nil {
			return url, err
		}
	}

	rows, err = db.Query(fmt.Sprintf(`SELECT url FROM mapping_%s where code='%s';`, tableSuffix, code))
	if err != nil {
		return url, err
	}
	for rows.Next() {
		err := rows.Scan(&url)
		if err != nil {
			return url, err
		}
	}

	return url, nil
}

func insertUrlFetchInfo(c *gin.Context, mongo *mgo.Database, code string) {
	data := map[string]interface{}{}
	req := c.Request

	for name, values := range req.Header {
		for _, value := range values {
			data[name] = value
		}
	}

	data["RemoteAddr"] = req.RemoteAddr
	data["RequestURI"] = req.RequestURI
	data["code"] = code
	mongo.C("analytics").Insert(data)

}