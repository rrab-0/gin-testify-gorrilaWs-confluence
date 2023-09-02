package routes

import (
	"example/unit-test-hello-world/controllers"

	"github.com/gin-gonic/gin"
)

func Mahasiswa(router *gin.RouterGroup) {
	mahasiswa := &controllers.MahasiswaController{}

	router.POST("/", mahasiswa.Create)
	router.GET("/", mahasiswa.Reads)
	router.GET("/:id", mahasiswa.Read)
	router.PATCH("/:id", mahasiswa.Update)
	router.DELETE("/:id", mahasiswa.Destroy)
}
