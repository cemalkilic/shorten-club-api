package middlewares

import (
    "errors"
    "fmt"
    "github.com/cemalkilic/shorten-backend/service"
    "github.com/dgrijalva/jwt-go"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

const BearerSchema = "Bearer"

// AuthorizeJWT Authorizes user by validating the token in Authorization header
// if no authorization header value is given, just passes
func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc{
    return func(c *gin.Context) {

        authToken := c.GetHeader("Authorization")
        if authToken == "" {
            c.Next()
            return
        }

        err := requireJWTToken(c.GetHeader("Authorization"))
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": err.Error(),
            })
            return
        }

        tokenString := authToken[len(BearerSchema):]
        tokenString = strings.TrimSpace(tokenString)

        token, err := jwtService.ValidateToken(tokenString)
        if err != nil {
            fmt.Printf("%v", err)
            c.AbortWithStatusJSON(400, gin.H{
                "error": "Given token is invalid!",
            })
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            c.Set("username", claims["username"])
            c.Next()
        } else {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Given JWT is invalid!",
            })
            return
        }
    }
}

func RequireJWTToken(jwtService service.JWTService) gin.HandlerFunc {
    return func(c *gin.Context) {

        err := requireJWTToken(c.GetHeader("Authorization"))
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": err.Error(),
            })
            return
        }
    }
}

func requireJWTToken(tokenString string) error {

    if tokenString == "" {
        return errors.New("authorization token is not given")
    }

    fmt.Printf("\nAuth token :: %s\n", tokenString)

    if !strings.HasPrefix(tokenString, BearerSchema) {
         return errors.New("invalid authorization token type")
    }

    return nil
}
