package controllers

import (
    "encoding/json"
    "github.com/cemalkilic/shorten-backend/database"
    "github.com/cemalkilic/shorten-backend/service"
    "github.com/cemalkilic/shorten-backend/utils/validator"
    "github.com/gin-gonic/gin"
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

func (cec *ShortenController) SetDB(dataStore database.DataStore) {
    cec.dataStore = dataStore
}

func (cec *ShortenController) GetContent(c *gin.Context) {
    url := c.Request.URL.Path

    srv := service.NewService(cec.dataStore, cec.validator)
    response, err := srv.GetContentBySlug(service.GetContentParams{Slug: url})
    if err != nil {
        internalError(c, err)
        return
    }

    if e, ok := response.Err.(error); ok && e != nil {
        internalError(c, e)
        return
    }

    record := response.Record
    permissions, _ := cec.getPermissions(c.GetString("username"), record.Slug)

    c.JSON(http.StatusOK, gin.H{
        "record":      record,
        "permissions": permissions,
    })
}

func (cec *ShortenController) UpdateRecord(c *gin.Context) {
    var updateRecordRequest service.UpdateRecordParams
    _ = c.ShouldBindJSON(&updateRecordRequest)

    username := c.GetString("username")
    if username == "" {
        c.AbortWithStatusJSON(400, gin.H{
            "error": "User not found!",
        })
        return
    }

    permissions, _ := cec.getPermissions(username, updateRecordRequest.Slug)
    if permissions["updateContent"] == false {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "error": "Not allowed to update!",
        })
        return
    }

    srv := service.NewService(cec.dataStore, cec.validator)
    response, err := srv.UpdateRecord(updateRecordRequest)
    if err != nil {
        internalError(c, err)
        return
    }

    if e, ok := response.Err.(error); ok && e != nil {
        internalError(c, e)
        return
    }

    record := response.Record

    c.JSON(http.StatusOK, gin.H{
        "record": record,
        "permissions": permissions,
    })
}

func (cec *ShortenController) InitialRecord(c *gin.Context) {

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

    srv := service.NewService(cec.dataStore, cec.validator)

    randomSlug := srv.GetRandomSlug()

    type M map[string]interface{}
    content, _ := json.Marshal(M{})

    params := service.AddRecordParams{
        Username: ctxUsername,
        Slug:     randomSlug,
        Type:     recordType,
        Content:  string(content),
    }
    response, err := srv.AddRecord(params)
    if err != nil {
        panic(err)
    }

    record := response.Record
    permissions, _ := cec.getPermissions(c.GetString("username"), record.Slug)

    c.JSON(http.StatusOK, gin.H{
        "record": record,
        "permissions": permissions,
    })
}

func (cec *ShortenController) getPermissions(username string, slug string) (map[string]bool, error) {
    permissionService := service.NewPermissionService(cec.dataStore)
    return permissionService.GetPermissionsBySlug(username, slug)
}
