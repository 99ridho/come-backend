package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/99ridho/come-backend/env"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v2"
)

var Dbm *gorp.DbMap

// initialize dbMap instance
func init() {

	dbName := env.Getenv("DB_NAME", "come")
	dbHost := env.Getenv("DB_HOST", "127.0.0.1")
	dbUsername := env.Getenv("DB_USERNAME", "root")
	dbPort := env.Getenv("DB_PORT", "3306")
	dbPassword := env.Getenv("DB_PASSWORD", "")
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbPort, dbName)
	log.Println(dbUrl)

	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		panic(err)
	}

	Dbm = &gorp.DbMap{
		Db: db,
		Dialect: gorp.MySQLDialect{
			Engine:   "InnoDB",
			Encoding: "UTF8",
		},
	}

	Dbm.TraceOn("[gorm]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
	Dbm.AddTableWithName(User{}, "users").SetKeys(true, "ID").AddIndex("EmailIndex", "Btree", []string{"email"}).SetUnique(true)
	Dbm.TraceOff()

}

// Create tables
func CreateTables() error {
	if err := Dbm.CreateTablesIfNotExists(); err != nil {
		return err
	}
	if err := Dbm.CreateIndex(); err != nil {
		return err
	}

	return nil
}
