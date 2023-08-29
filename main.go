package main

import (
	"example/unit-test-hello-world/config"
	mahasiswa "example/unit-test-hello-world/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func homepage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Service is up and running.",
	})
}

func main() {
	// load env
	if envErr := godotenv.Load(".env"); envErr != nil {
		log.Println("Could not load .env")
		os.Exit(1)
	}

	// init server & db
	app := gin.Default()
	config.InitDB("PostgreSQL")

	app.GET("/", homepage)
	v1 := app.Group("/api/v1")
	mahasiswa.CrudMahasiswa(v1.Group("/mahasiswa"))

	DEV_HOST := os.Getenv("DEV_HOST")
	DEV_PORT := os.Getenv("DEV_PORT")
	app.Run(fmt.Sprintf("%v:%v", DEV_HOST, DEV_PORT))
}
