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
	//authHandler := NewAuthHandler(ctx, cfg, mongoUsersClient)

	// router.POST("/signin", authHandler.SignIn)
	router.GET("/interviews", interviewsHandler.ListInterviews)
	router.GET("/interviews/:id", interviewsHandler.GetInterview)

	// API private endpoints
	authorised := router.Group("/")
	//authorised.Use(m.AuthCookieMiddleware())
	// authorised.POST("/signout", authHandler.SignOut)
	authorised.POST("/interviews", interviewsHandler.NewInterview)
	authorised.PUT("/interviews/:id", interviewsHandler.UpdateInterview)
	authorised.DELETE("/interviews/:id", interviewsHandler.DeleteInterview)
	authorised.POST("/comment/:id", interviewsHandler.AddComment)
	authorised.PUT("/comment/:idInterview/:idComment", interviewsHandler.UpdateComment)
	authorised.DELETE("/comment/:idInterview/:idComment", interviewsHandler.DeleteComment)

	return router
}
