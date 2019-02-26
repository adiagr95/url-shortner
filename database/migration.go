package database

import (
	"database/sql"
	"fmt"
	"../ctrl"
)

func Migrate(db *sql.DB)  {
	stmt, err := db.Prepare("CREATE TABLE config (name VARCHAR(255), value VARCHAR(255), UNIQUE KEY `name_UNIQUE` (`name`));")
	_, err = stmt.Exec()
	check(err)

	if err == nil {
		fmt.Println("CREATED TABLE Config")
		stmt, _ := db.Prepare(`INSERT INTO config VALUES ("CURRENT_CREATED_TABLE", "0");`)
		stmt.Exec()

		stmt, _ = db.Prepare(`INSERT INTO config VALUES ("CURRENT_TABLE", "1" );`)
		stmt.Exec()

		stmt, _ = db.Prepare(`INSERT INTO config VALUES ("CURRENT_CODE", "");`)
		stmt.Exec()
	}

	stmt, err = db.Prepare("CREATE TABLE `mapping_master` ( `start` bigint(255) NOT NULL, `end` bigint(255) NOT NULL, " +
		"`table_suffix` varchar(45) NOT NULL, UNIQUE KEY `table_suffix_UNIQUE` (`table_suffix`))")
	_, err = stmt.Exec()
	check(err)

	if err == nil {
		fmt.Println("CREATED TABLE Mapping Master")
		CreateMappingTables(db, 15)
	}
}

func CreateMappingTables(db *sql.DB, num int) {

	rows, err := db.Query(`SELECT value FROM config where name='CURRENT_CREATED_TABLE';`)
	var createdTablesCount int
	for rows.Next() {
		err = rows.Scan(&createdTablesCount)
		check(err)

	}

	for i :=1 ; i <= num; i++ {
		query := fmt.Sprintf("CREATE TABLE `mapping_%d` ( `code` varchar(255) NOT NULL, `url` longtext NOT NULL, PRIMARY KEY (`code`) )",createdTablesCount + i)
		stmt, err := db.Prepare(query)
		_, err = stmt.Exec()
		check(err)

		if err == nil {
			fmt.Println("CREATED TABLE Mapping Table")
		}

		start := ctrl.GetNumberFromCode(ctrl.GetNthCode(((createdTablesCount + i - 1) * 10) + 1))
		end := ctrl.GetNumberFromCode(ctrl.GetNthCode((createdTablesCount + i) * 10))

		query = fmt.Sprintf("INSERT INTO mapping_master (`start`, `end`, `table_suffix`) VALUES (%f, %f, %d)", start, end, createdTablesCount + i)
		fmt.Println(query)
		stmt, _ = db.Prepare(query)
		stmt.Exec()

	}


	stmt, err := db.Prepare(fmt.Sprintf(`UPDATE config SET value="%d" WHERE name="CURRENT_CREATED_TABLE";`,createdTablesCount + num))
	stmt.Exec()
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
