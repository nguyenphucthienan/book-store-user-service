package application

import (
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func Start() {
	mapUrls()

	if err := router.Run("localhost:8082"); err != nil {
		panic(err)
	}
}
