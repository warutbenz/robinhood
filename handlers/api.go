package handlers

import (
	"context"
	"net/http"

	"robinhood/datastore"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Setup(ctx context.Context, cfg *viper.Viper, router *gin.Engine, mongoInterviewsClient datastore.Interviews,
	mongoUsersClient datastore.Users) *gin.Engine {

	router.GET("/healthCheck", func(c *gin.Context) {
		c.String(http.StatusOK, "service start")
	})
	interviewsHandler := NewInterviewsHandler(ctx, cfg, mongoInterviewsClient)
	router.GET("/interviews", interviewsHandler.ListInterviews)
	router.GET("/interviews/:id", interviewsHandler.GetInterview)
	router.POST("/interviews", interviewsHandler.NewInterview)
	router.PUT("/interviews/:id", interviewsHandler.UpdateInterview)
	router.DELETE("/interviews/:id", interviewsHandler.DeleteInterview)
	router.POST("/comment/:id", interviewsHandler.AddComment)
	router.PUT("/comment/:idInterview/:idComment", interviewsHandler.UpdateComment)
	router.DELETE("/comment/:idInterview/:idComment", interviewsHandler.DeleteComment)

	return router
}
