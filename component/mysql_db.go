package component

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

var db *sqlx.DB

func init() {
	var err error
	db,err = sqlx.Open("mysql", "hackway:0663@tcp(127.0.0.1:3306)/fairy_cms?charset=utf8")
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}