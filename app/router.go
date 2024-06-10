package app

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"passVault/dependency"
	"passVault/dtos"
	"passVault/interfaces"
	"passVault/models"
	"passVault/resources"
	"strings"
	"time"
)

func SetupRouter(ctx context.Context) *gin.Engine {
	var (
		config = resources.Config()
	)

	if config.Get(dtos.ConfigKeys.Env) != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	// New Gin Engine.
	router := gin.New()

	router.Use(AccessLogMiddleware())

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://www.passvault.fun",
			"http://passvault.fun",
			"https://passvault.fun",
			"http://www.passvault.fun",
		},
		AllowHeaders:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	router.Use(func(ginCtx *gin.Context) {
		allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization," +
			" X-CSRF-Token,  Cache-Control, request-id, version, api-key"
		ginCtx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ginCtx.Writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	})

	router.Use()

	if config.Get(dtos.ConfigKeys.Env) == "dev" {
		pprof.Register(router)
	}

	var (
		api = router.Group("/api")
		v1  = api.Group("/v1")
	)

	SetUpPasswordsEndpoints(v1)
	SetUpUsersEndpoints(v1)

	return router
}

func AccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		logger := resources.Logger(c.Request.Context())
		fields := []any{
			slog.Int("status", c.Writer.Status()),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("ip", c.ClientIP()),
			slog.String("user-agent", c.Request.UserAgent()),
			slog.Duration("latency", latency),
			slog.Time("time", end),
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e, fields...)
			}
		} else {
			logger.Info("", fields...)
		}
	}
}

func SetUpPasswordsEndpoints(router *gin.RouterGroup) {
	dependencies := dependency.Dependencies()

	router.POST("/import/passwords",
		SetUserMiddleware(dependencies.UserService),
		ImportPasswords(dependencies.PasswordService),
	)

	router.POST("/generate/passwords",
		SetUserMiddleware(dependencies.UserService),
		GeneratePasswords(dependencies.PasswordService),
	)

	passwords := router.Group("/passwords", SetUserMiddleware(dependencies.UserService))

	{
		passwords.POST("", StorePassword(dependencies.PasswordService))

		passwords.GET("", GetPasswords(dependencies.PasswordService))

		passwords.GET("/:id", GetPassword(dependencies.PasswordService))

		passwords.PUT("/:id", UpdatePassword(dependencies.PasswordService))

		passwords.DELETE("/:id", DeletePassword(dependencies.PasswordService))
	}

}

func SetUpUsersEndpoints(router *gin.RouterGroup) {

	var (
		dependencies = dependency.Dependencies()
	)

	router.POST("/login/users", LoginUser(dependencies.UserService))

	users := router.Group("/users")

	{
		users.POST("", CreateUser(dependencies.UserService))
	}

}

func SetUserMiddleware(
	userService interfaces.UserService,
) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			user models.User
		)

		authHeader := ginCtx.GetHeader("Authorization")
		if authHeader == "" {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			ginCtx.Abort()
			return
		}
		bearerTokenParts := strings.Split(authHeader, " ")

		if len(bearerTokenParts) != 2 {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
			ginCtx.Abort()
			return
		}

		if bearerTokenParts[0] != "Bearer" {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
			ginCtx.Abort()
			return
		}

		if err := userService.ValidateToken(ginCtx.Request.Context(), bearerTokenParts[1], &user); err != nil {
			ginCtx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ginCtx.Abort()
			return
		}
		setUserInGinContext(ginCtx, user)
		ginCtx.Next()
	}
}
