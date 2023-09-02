package main

import (
	"example/unit-test-hello-world/config"
	mahasiswa "example/unit-test-hello-world/routes"
	ws "example/unit-test-hello-world/websocket"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "example/unit-test-hello-world/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func Homepage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Service is up and running.",
	})
}

// @title           Crud Mahasiswa
// @version         4.0
// @description     This is where we do crud for mahasiswas.
// @termsOfService  http://youtube.com

// @contact.name   holy wow
// @contact.url    http://comehere.dev
// @contact.email  holyWow@macaroni.dev

// @license.name  MIT
// @license.url   http://mit.dev

// @host      localhost:8080
// @BasePath  /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// load env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Could not load .env")
		os.Exit(1)
	}

	// init server & db
	config.InitDB("PostgreSQL")
	app := gin.Default()
	app.Use(cors.Default())

	// default routes
	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.GET("/ws", ws.ReadIncomingMessage)
	app.GET("/", Homepage)

	v1 := app.Group("/api/v1")
	mahasiswa.Crud(v1.Group("/mahasiswa"))

	DEV_PORT := os.Getenv("DEV_PORT")
	app.Run(fmt.Sprintf(":%v", DEV_PORT))
}
