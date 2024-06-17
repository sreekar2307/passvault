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

func BeginLoginUser(userService interfaces.UserService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			requestBody dtos.BeginLoginParams
		)
		if err := ginCtx.BindJSON(&requestBody); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := userService.BeginWebAuthnLogin(ginCtx.Request.Context(), requestBody)
		if err != nil {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ginCtx.JSON(http.StatusOK, result)
	}
}

func FinishLoginUser(userService interfaces.UserService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		sessionID := ginCtx.Query("session_id")
		if sessionID == "" {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
			return
		}
		defer ginCtx.Request.Body.Close()
		authToken, err := userService.FinishWebAuthnLogin(ginCtx.Request.Context(), sessionID, ginCtx.Request.Body)
		if err != nil {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ginCtx.JSON(http.StatusOK, gin.H{"token": authToken})
	}
}

func BeginRegisterUser(userService interfaces.UserService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			requestBody dtos.RegisterWebAuthnUserParams
		)
		if err := ginCtx.BindJSON(&requestBody); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		credOptions, err := userService.BeginWebAuthnRegister(ginCtx.Request.Context(), requestBody)
		if err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ginCtx.JSON(http.StatusCreated, credOptions)
	}
}

func FinishRegisterUser(userService interfaces.UserService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		sessionID := ginCtx.Query("session_id")
		if sessionID == "" {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
			return
		}
		defer ginCtx.Request.Body.Close()
		authToken, err := userService.FinishWebAuthnRegister(ginCtx.Request.Context(), sessionID, ginCtx.Request.Body)
		if err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ginCtx.JSON(http.StatusCreated, gin.H{"token": authToken})
	}
}
