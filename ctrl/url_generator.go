package ctrl

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"os"
)


type RequestModel struct {
	Url string `form:"url" json:"url" binding:"required"`
}

func CreateUrl(c *gin.Context)  {
	db := c.MustGet("db").(sql.DB)
	requestData := RequestModel{}
	err := c.ShouldBindWith(&requestData, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" :string(err.Error())})
		return
	}
	url := requestData.Url

	var code, currentTable string

	rows, err := db.Query(`SELECT value FROM config where name='CURRENT_CODE';`)
	for rows.Next() {
		err = rows.Scan(&code)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" :string(err.Error())})
		return
	}


	rows, err = db.Query(`SELECT value FROM config where name='CURRENT_TABLE';`)
	for rows.Next() {
		err = rows.Scan(&currentTable)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" :string(err.Error())})
		return
	}

	code = GetNextString(code)

	stmt, _ := db.Prepare(fmt.Sprintf(`INSERT INTO mapping_%s VALUES ("%s", "%s");`, currentTable, code, url))
	_, err = stmt.Exec()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" :string(err.Error())})
		return
	}

	err = updateConfigTable(db, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" :string(err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url" : url,
		"short_url": os.Getenv("SHORT_DOMAIN") + code,
	})
}

func updateConfigTable(db sql.DB, code string) error {
	nextCodeNumber := GetNumberFromCode(GetNextString(code))
	var tableSuffix string

	rows, _ := db.Query(fmt.Sprintf(`SELECT table_suffix FROM mapping_master where start<=%f and end>=%f;`, nextCodeNumber, nextCodeNumber))
	for rows.Next() {
		err := rows.Scan(&tableSuffix)
		if err != nil {
			return err
		}
	}

	stmt, err := db.Prepare(fmt.Sprintf(`UPDATE config SET value="%s" WHERE name="CURRENT_CODE";`, code))
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	stmt, err = db.Prepare(fmt.Sprintf(`UPDATE config SET value="%s" WHERE name="CURRENT_TABLE";`, tableSuffix))
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
