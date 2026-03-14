package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

func InitDb() (*sql.DB, error){
    cfg := mysql.NewConfig()
    cfg.User = os.Getenv("DBUSER")
    cfg.Passwd = os.Getenv("DBPASS")
    cfg.Net = "tcp"
    cfg.Addr = "127.0.0.1:3306"
    cfg.DBName = "gosql"

    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatal(err)
		return nil, err
    }

	pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
		return nil, err
    }
    fmt.Println("Connected!")

	return db, nil
}