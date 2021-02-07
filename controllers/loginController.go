package controllers

import (
    "fmt"
    "github.com/cemalkilic/shorten-backend/models"
    "github.com/cemalkilic/shorten-backend/service"
    "github.com/gin-gonic/gin"
    "crypto/rand"
    "github.com/golang/glog"
    "net/http"
)

type LoginController struct {
    loginService service.LoginService
    jWtService   service.JWTService
}

func NewLoginController(loginService service.LoginService, jwtService service.JWTService) *LoginController {
    return &LoginController{
        loginService: loginService,
        jWtService:   jwtService,
    }
}

func getRandomUsername(prefix string) string {
    n := 5
    b := make([]byte, n)
    if _, err := rand.Read(b); err != nil {
        panic(err)
    }
    randomUsername := fmt.Sprintf("%X", b)
    return prefix + randomUsername
}

func (controller *LoginController) Login(c *gin.Context) {
    var credential models.User
    err := c.ShouldBindJSON(&credential)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Login data must be a valid JSON!",
        })
        return
    }

    var token string
    isUserAuthenticated := controller.loginService.IsValidCredentials(credential)
    if isUserAuthenticated {
        token, err = controller.jWtService.GenerateToken(credential.Username)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
        }
    } else {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Given credentials did not match!",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": token,
    })
}

func (controller *LoginController) Auth(c *gin.Context) {
    // Hacky way to pass user creation if authorization header exists
    authToken := c.GetHeader("Authorization")
    if authToken != "" {
        glog.Info("Request with authorization header. Skipping create user!")
        c.AbortWithStatusJSON(http.StatusOK, gin.H{
            "token": authToken,
        })
        return
    }

    var credential models.User

    credential.Username = getRandomUsername("random")
    credential.Password = "strongPasswordToHash"

    err := controller.loginService.Signup(credential)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    // Return the token to user, to let her login right after signup
    token, err := controller.jWtService.GenerateToken(credential.Username)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.AbortWithStatusJSON(http.StatusOK, gin.H{
        "token": token,
    })
    return
}
