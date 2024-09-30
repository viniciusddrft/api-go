package routes

import (
	controllers "api/api/controllers"
	"api/api/db"

	"github.com/gin-gonic/gin"
)

func AppRoutes(router *gin.Engine) *gin.RouterGroup {
	tweetController := controllers.NewTweetController(db.DB)
	v1 := router.Group("/v1")
	v1.GET("/tweets", tweetController.FindAll)
	v1.POST("/tweet", tweetController.Create)
	v1.PUT("/tweet/:id", tweetController.Update)
	v1.DELETE("/tweet/:id", tweetController.Delete)
	return v1
}
