package controllers

import (
    "github.com/cemalkilic/shorten-backend/database"
    "github.com/cemalkilic/shorten-backend/service"
    "github.com/cemalkilic/shorten-backend/utils/validator"
    "github.com/gin-gonic/gin"
    "github.com/golang/glog"
    "net/http"
)

type ShortenController struct {
    dataStore database.DataStore
    validator *validator.CustomValidator
}

func NewShortenController(db database.DataStore, v *validator.CustomValidator) *ShortenController {
    return &ShortenController{
        dataStore: db,
        validator: v,
    }
}

func (sc *ShortenController) SetDB(dataStore database.DataStore) {
    sc.dataStore = dataStore
}

func (sc *ShortenController) GetContent(c *gin.Context) {
    url := c.Request.URL.Path

    srv := service.NewService(sc.dataStore, sc.validator)
    response, err := srv.GetContentBySlug(service.GetContentParams{Slug: url})
    if err != nil {
        glog.Error(err)
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "error": "Internal server occurred",
        })
        return
    }

    if e, ok := response.Err.(error); ok && e != nil {
        c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
            "error": e.Error(),
        })
        return
    }

    record := response.Record
    permissions, _ := sc.getPermissions(c.GetString("username"), record.Slug)

    c.JSON(http.StatusOK, gin.H{
        "record":      record,
        "permissions": permissions,
    })
}

func (sc *ShortenController) UpdateRecord(c *gin.Context) {
    var updateRecordRequest service.UpdateRecordParams
    _ = c.ShouldBindJSON(&updateRecordRequest)

    username := c.GetString("username")
    if username == "" {
        c.AbortWithStatusJSON(400, gin.H{
            "error": "User not found!",
        })
        return
    }

    permissions, _ := sc.getPermissions(username, updateRecordRequest.Slug)
    if permissions["updateContent"] == false {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "error": "Not allowed to update!",
        })
        return
    }

    srv := service.NewService(sc.dataStore, sc.validator)
    response, err := srv.UpdateRecord(updateRecordRequest)
    if err != nil {
        glog.Error(err)
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "error": "Internal server occurred",
        })
        return
    }

    if e, ok := response.Err.(error); ok && e != nil {
        c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
            "error": e.Error(),
        })
        return
    }

    record := response.Record

    c.JSON(http.StatusOK, gin.H{
        "record": record,
        "permissions": permissions,
    })
}

func (sc *ShortenController) InitialRecord(c *gin.Context) {

    ctxUsername := c.GetString("username")
    if ctxUsername == "" {
        c.AbortWithStatusJSON(400, gin.H{
            "error": "User not found!",
        })
        return
    }

    recordType := c.Query("type")
    if recordType == "" {
        recordType = "LINK" // TODO :: must check enums! LINK vs NOTE
    }

    srv := service.NewService(sc.dataStore, sc.validator)

    randomSlug := srv.GetRandomSlug()

    params := service.AddRecordParams{
        Username: ctxUsername,
        Slug:     randomSlug,
        Type:     recordType,
        Content:  make([]string, 0),
    }
    response, err := srv.AddRecord(params)
    if err != nil {
        glog.Error("Could not add record: " + err.Error())
    }

    record := response.Record
    permissions, _ := sc.getPermissions(c.GetString("username"), record.Slug)

    c.JSON(http.StatusOK, gin.H{
        "record": record,
        "permissions": permissions,
    })
}

func (sc *ShortenController) getPermissions(username string, slug string) (map[string]bool, error) {
    permissionService := service.NewPermissionService(sc.dataStore)
    return permissionService.GetPermissionsBySlug(username, slug)
}
