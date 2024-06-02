package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"passVault/dtos"
	"passVault/interfaces"
)

func StorePassword(passwordService interfaces.PasswordService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		user, ok := userFromGinContext(ginCtx)
		if !ok {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		var (
			requestBody dtos.StorePasswordParams
		)
		if err := ginCtx.BindJSON(&requestBody); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		stored, err := passwordService.StorePassword(ginCtx.Request.Context(), user, requestBody)
		if err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ginCtx.JSON(http.StatusCreated, gin.H{"stored": stored})
	}
}

func GetPasswords(passwordService interfaces.PasswordService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		user, ok := userFromGinContext(ginCtx)
		if !ok {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		var (
			params dtos.GetPasswordsParams
		)
		if err := ginCtx.BindQuery(&params); err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		passwords, err := passwordService.GetPasswords(ginCtx.Request.Context(), user, params)
		if err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ginCtx.JSON(http.StatusOK, gin.H{"data": passwords})
	}
}

func GetPassword(passwordService interfaces.PasswordService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		user, ok := userFromGinContext(ginCtx)
		if !ok {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		var (
			params struct {
				ID uint `uri:"id" binding:"required,min=1"`
			}
		)
		if err := ginCtx.BindUri(&params); err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		password, err := passwordService.GetPassword(ginCtx.Request.Context(), user, params.ID)
		if err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ginCtx.JSON(http.StatusOK, gin.H{"data": password})
	}
}

func ImportPasswords(passwordService interfaces.PasswordService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		user, ok := userFromGinContext(ginCtx)
		if !ok {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		fileHeader, err := ginCtx.FormFile("passwords")
		if err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		imported, err := passwordService.ImportPasswords(ginCtx.Request.Context(), user, file)
		if err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ginCtx.JSON(http.StatusCreated, gin.H{"imported": imported})
	}
}

func GeneratePasswords(passwordService interfaces.PasswordService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		user, ok := userFromGinContext(ginCtx)
		if !ok {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		var (
			requestBody dtos.GeneratePasswordParams
		)
		if err := ginCtx.BindJSON(&requestBody); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		password, err := passwordService.GeneratePassword(ginCtx.Request.Context(), user, requestBody)
		if err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ginCtx.JSON(http.StatusOK, gin.H{"password": password})
		return
	}
}

func UpdatePassword(passwordService interfaces.PasswordService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		user, ok := userFromGinContext(ginCtx)
		if !ok {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		var (
			params dtos.UpdatePasswordParams
		)
		if err := ginCtx.BindUri(&params); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := ginCtx.Bind(&params); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ok, err := passwordService.UpdatePassword(ginCtx.Request.Context(), user, params)
		if err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ginCtx.JSON(http.StatusOK, gin.H{"updated": ok})
	}
}

func DeletePassword(passwordService interfaces.PasswordService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		user, ok := userFromGinContext(ginCtx)
		if !ok {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		var (
			params struct {
				ID uint `uri:"id" binding:"required,min=1"`
			}
		)
		if err := ginCtx.BindUri(&params); err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ok, err := passwordService.DeletePassword(ginCtx.Request.Context(), user, params.ID)
		if err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ginCtx.JSON(http.StatusOK, gin.H{"deleted": ok})
	}
}
