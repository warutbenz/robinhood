package handlers

import (
	"context"
	"net/http"
	"strconv"

	"robinhood/model"

	"robinhood/datastore"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type InterviewsHandler struct {
	ctx          context.Context
	cfg          *viper.Viper
	mongoDBStore datastore.Interviews
}

func NewInterviewsHandler(ctx context.Context, cfg *viper.Viper, mongoDBStore datastore.Interviews) *InterviewsHandler {
	return &InterviewsHandler{
		ctx:          ctx,
		cfg:          cfg,
		mongoDBStore: mongoDBStore,
	}
}

func (handler *InterviewsHandler) NewInterview(ctx *gin.Context) {
	var interview *model.Interview
	if err := ctx.ShouldBindJSON(&interview); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"invalid input": err.Error(),
		})
		return
	}
	if err := handler.mongoDBStore.AddInterview(handler.ctx, interview); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, interview)
}

func (handler *InterviewsHandler) ListInterviews(ctx *gin.Context) {
	interviews, err := handler.mongoDBStore.ListInterviews(handler.ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, interviews)
}

func (handler *InterviewsHandler) GetInterview(ctx *gin.Context) {
	id := ctx.Param("id")
	interview, err := handler.mongoDBStore.GetInterview(handler.ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if interview.Subject == "" {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, interview)
}

func (handler *InterviewsHandler) UpdateInterview(ctx *gin.Context) {
	id := ctx.Param("id")
	var interview model.Interview
	if err := ctx.ShouldBindJSON(&interview); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"invalid input": err.Error(),
		})
		return
	}
	validStatuses := []string{"To do", "In Progress", "Done"}
	found := false
	for _, s := range validStatuses {
		if interview.Status == s {
			found = true
			break
		}
	}
	if !found {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Status is not in the list To do,In Progress, Done",
		})
		return
	}

	modifiedCount, err := handler.mongoDBStore.UpdateInterview(handler.ctx, id, interview)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if modifiedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "interview not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Interview has been Updated",
	})
}

func (handler *InterviewsHandler) DeleteInterview(ctx *gin.Context) {
	id := ctx.Param("id")

	deletedCount, err := handler.mongoDBStore.DeleteInterview(handler.ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if deletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Interview not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Interview has been deleted",
	})
}

func (handler *InterviewsHandler) AddComment(ctx *gin.Context) {
	id := ctx.Param("id")
	var comment model.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"invalid input": err.Error(),
		})
		return
	}
	interview, err := handler.mongoDBStore.GetInterview(handler.ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if interview.Subject == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	if len(interview.Comments) == 0 {
		comment.ID = 1
	} else {
		comment.ID = interview.Comments[len(interview.Comments)-1].ID + 1
	}

	modifiedCount, err := handler.mongoDBStore.AddComment(handler.ctx, id, comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if modifiedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "interview not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Comment has been Create",
	})
}
func (handler *InterviewsHandler) UpdateComment(ctx *gin.Context) {
	id := ctx.Param("idInterview")
	// Convert the string to an integer
	cid, err := strconv.Atoi(ctx.Param("idComment"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	var comment model.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"invalid input": err.Error(),
		})
		return
	}
	modifiedCount, err := handler.mongoDBStore.UpdateComment(handler.ctx, id, cid, comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if modifiedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "comment not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "comment has been Updated",
	})
}
func (handler *InterviewsHandler) DeleteComment(ctx *gin.Context) {
	id := ctx.Param("idInterview")

	// Convert the string to an integer
	cid, err := strconv.Atoi(ctx.Param("idComment"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	deletedCount, err := handler.mongoDBStore.DeleteComment(handler.ctx, id, cid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if deletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "comment not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "comment has been deleted",
	})
}
