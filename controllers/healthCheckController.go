package controllers

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    "net/http"
)

type HealthCheckController struct {
    database *sql.DB
}

func NewHealthCheckController(db *sql.DB) *HealthCheckController {
    return &HealthCheckController{
        database: db,
    }
}

func (hc *HealthCheckController) HealthCheck(c *gin.Context) {
    // Check database
    err := hc.database.Ping()
    if err != nil {
        c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
            "healthy": false,
            "status": "MySQL down",
        })
        return
    }

    c.AbortWithStatusJSON(http.StatusOK, gin.H{
        "healthy": true,
        "status": "OK",
    })
}
