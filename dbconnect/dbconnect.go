package dbconnect

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Connectdatabase() *sql.DB{

	db,err:=sql.Open("mysql","root:admin@tcp(localhost:3306)/bloodbank")
	if(err!=nil){
		log.Print(err);
	}

	return db
}