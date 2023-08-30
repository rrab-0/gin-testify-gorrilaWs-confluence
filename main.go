package main

import (
	"example/unit-test-hello-world/config"
	mahasiswa "example/unit-test-hello-world/routes"
	"example/unit-test-hello-world/ws"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "example/unit-test-hello-world/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func homepage(c *gin.Context) {
	go ws.BroadcastMessage(c.Writer, c.Request, "Hello WS")
	c.JSON(http.StatusOK, gin.H{
		"message": "Service is up and running.",
	})
}

// @title           Crud Mahasiswa
// @version         6.9
// @description     This is where we do crud for mahasiswas.
// @termsOfService  http://youtube.com

// @contact.name   holy wow
// @contact.url    http://comehere.dev
// @contact.email  holyWow@macaroni.dev

// @license.name  MIT
// @license.url   http://mit.dev

// @host      localhost:8081
// @BasePath  /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// load env
	if envErr := godotenv.Load(".env"); envErr != nil {
		log.Println("Could not load .env")
		os.Exit(1)
	}

	// init server & db
	app := gin.Default()
	config.InitDB("PostgreSQL")

	// docs route
	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.GET("/ws", func(c *gin.Context) {
		ws.WsHandler(c.Writer, c.Request)
	})
	app.GET("/", homepage)

	v1 := app.Group("/api/v1")
	mahasiswa.CrudMahasiswa(v1.Group("/mahasiswa"))

	DEV_HOST := os.Getenv("DEV_HOST")
	DEV_PORT := os.Getenv("DEV_PORT")
	app.Run(fmt.Sprintf("%v:%v", DEV_HOST, DEV_PORT))
}
