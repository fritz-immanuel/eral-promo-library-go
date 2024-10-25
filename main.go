package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fritz-immanuel/eral-promo-library-go/configs"
	"github.com/fritz-immanuel/eral-promo-library-go/databases"
	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/src/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"
)

var loc *time.Location

type stackTracer interface {
	StackTrace() errors.StackTrace
}

var addr = flag.String("addr", ":8080", "http service address")

// Init function for initialize config
func init() {

}

// Main function for start entry golang
func main() {
	gin.SetMode(gin.ReleaseMode)
	os.Setenv("TZ", "Asia/Jakarta")

	config, err := configs.GetConfiguration()
	if err != nil {
		log.Fatalln("failed to get configuration: ", err)
	}

	configs.AppConfig = config

	db, err := sqlx.Open("mysql", config.DBConnectionString)
	if err != nil {
		log.Fatalln("failed to open database x: ", err)
	}
	defer db.Close()

	dataManager := data.NewManager(
		db,
	)

	databases.MigrateUp()

	fmt.Println("Server is running...")
	routes.RegisterRoutes(db, config, dataManager)
}
