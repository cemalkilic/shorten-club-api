package main

import (
    "flag"
    "github.com/cemalkilic/shorten-backend/config"
    "github.com/cemalkilic/shorten-backend/controllers"
    "github.com/cemalkilic/shorten-backend/database"
    "github.com/cemalkilic/shorten-backend/middlewares"
    "github.com/cemalkilic/shorten-backend/service"
    "github.com/cemalkilic/shorten-backend/utils/validator"
    "github.com/gin-gonic/gin"
    ginglog "github.com/szuecs/gin-glog"
    "time"
)

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func main() {
    flag.Parse()
    cfg, _ := config.LoadConfig(".")
    if cfg.GinMode == "release" {
        gin.SetMode(cfg.GinMode)
    }

    router := gin.New()

    router.Use(ginglog.Logger(3 * time.Second))
    router.Use(gin.Recovery())
    router.Use(CORSMiddleware())

    mysqlHandler := database.NewMySQLDBHandler(cfg)
    dataStore := database.GetSQLDataStore(mysqlHandler)
    userStore := database.GetSQLUserStore(mysqlHandler)

    v := validator.NewValidator()

    shortenController := controllers.NewShortenController(dataStore, v)
    shortenController.SetDB(dataStore)

    loginService := service.DBLoginService(userStore, v)
    jwtService := service.JWTAuthService(cfg)
    loginController := controllers.NewLoginController(loginService, jwtService)

    healthCheckController := controllers.NewHealthCheckController(mysqlHandler)
    router.GET("/healthcheckz", healthCheckController.HealthCheck)

    router.Use(middlewares.AuthorizeJWT(jwtService))

    authorized := router.Group("/")
    // Routes that are not available without a valid JWT token
    authorized.Use(middlewares.RequireJWTToken(jwtService))
    {
        authorized.GET("/initial", shortenController.InitialRecord)
        authorized.POST("/updateRecord", shortenController.UpdateRecord)
    }

    router.GET("/auth", loginController.Auth)

    // Default handler to handle user routes
    router.NoRoute(shortenController.GetContent)

    router.Run(cfg.ServerAddress + ":" + cfg.ServerPort)
}
