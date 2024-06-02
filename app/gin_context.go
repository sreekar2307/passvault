package app

import (
	"github.com/gin-gonic/gin"
	"passVault/models"
)

var ginContextKeys = struct {
	User string
}{
	User: "user",
}

func userFromGinContext(ginCtx *gin.Context) (models.User, bool) {
	userAsAny, ok := ginCtx.Get(ginContextKeys.User)
	if !ok {
		return models.User{}, false
	}
	user, ok := userAsAny.(models.User)
	return user, ok
}

func setUserInGinContext(ginCtx *gin.Context, user models.User) {
	ginCtx.Set(ginContextKeys.User, user)
}
