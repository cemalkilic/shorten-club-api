package main

import (
    "github.com/cemalkilic/shorten-backend/config"
    "github.com/cemalkilic/shorten-backend/controllers"
    "github.com/cemalkilic/shorten-backend/database"
    "github.com/cemalkilic/shorten-backend/middlewares"
    "github.com/cemalkilic/shorten-backend/service"
    "github.com/cemalkilic/shorten-backend/utils/validator"
    "github.com/gin-gonic/gin"
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
    router := gin.Default()
    router.Use(CORSMiddleware())


    cfg, _ := config.LoadConfig(".")

    mysqlHandler := database.NewMySQLDBHandler(cfg)
    dataStore := database.GetSQLDataStore(mysqlHandler)
    userStore := database.GetSQLUserStore(mysqlHandler)

    v := validator.NewValidator()

    shortenController := controllers.NewShortenController(dataStore, v)
    shortenController.SetDB(dataStore)

    loginService := service.DBLoginService(userStore, v)
    jwtService := service.JWTAuthService(cfg)
    loginController := controllers.NewLoginController(loginService, jwtService)

    router.POST("/login", loginController.Login)
    router.POST("/signup", loginController.Signup)
    router.GET("/user/me", middlewares.AuthorizeJWT(jwtService), func(context *gin.Context) {
        context.JSON(200, gin.H{
            "success": true,
        })
    })

    router.GET("/initial", loginController.Signup, middlewares.AuthorizeJWT(jwtService), shortenController.InitialRecord)

    // Default handler to handle user routes
    router.NoRoute(shortenController.GetContent)
    router.POST("/addRecord", middlewares.AuthorizeJWT(jwtService), shortenController.AddRecord)
    router.POST("/updateRecord", middlewares.AuthorizeJWT(jwtService), shortenController.UpdateRecord)

    router.Run(cfg.ServerAddress)
}
