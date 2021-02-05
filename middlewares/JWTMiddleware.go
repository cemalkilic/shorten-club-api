package middlewares

import (
    "fmt"
    "github.com/cemalkilic/shorten-backend/service"
    "github.com/dgrijalva/jwt-go"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc{
    return func(c *gin.Context) {
        const BearerSchema = "Bearer"

        authToken := c.GetHeader("Authorization")
        if authToken == "" {
            ctxToken, exists := c.Get("userTokenAsBearer")
            fmt.Printf("\n\nctxToken: %s\n\n", ctxToken)
            if exists == false {
                c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                    "error": "Authorization token is not given!",
                })
                return
            }

            // TODO :: Check what happens when ctxToken is not string
            authToken = ctxToken.(string)

            /*
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                "error": "Authorization token is not given!",
            })
            return
            */
        }


        fmt.Printf("\nauth token:: %s\n", authToken)

        if !strings.HasPrefix(authToken, BearerSchema) {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                "error": "Authorization token must be type of Bearer!",
            })
            return
        }

        tokenString := authToken[len(BearerSchema):]
        tokenString = strings.TrimSpace(tokenString)


        c.Set("userToken", tokenString)
        fmt.Printf("\n\ntokenstring: %s\n\n", tokenString)

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
            //c.Next()
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Given JWT is invalid!",
            })
            return
        }
    }
}
