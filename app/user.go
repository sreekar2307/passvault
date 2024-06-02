package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"passVault/dtos"
	"passVault/interfaces"
)

func CreateUser(userService interfaces.UserService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			requestBody dtos.CreateUserParams
		)
		if err := ginCtx.BindJSON(&requestBody); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		authToken, err := userService.CreateUser(ginCtx.Request.Context(), requestBody)
		if err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ginCtx.JSON(http.StatusCreated, gin.H{"token": authToken})
	}
}

func LoginUser(userService interfaces.UserService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			requestBody dtos.LoginParams
		)
		if err := ginCtx.BindJSON(&requestBody); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		authToken, err := userService.Login(ginCtx.Request.Context(), requestBody)
		if err != nil {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ginCtx.JSON(http.StatusOK, gin.H{"token": authToken})
	}
}
