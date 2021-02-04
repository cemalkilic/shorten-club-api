package controllers

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func internalError(c *gin.Context, err error) {
    c.JSON(http.StatusInternalServerError, gin.H{
        "error": err.Error(),
    })
}
