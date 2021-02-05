package controllers

import (
    "encoding/json"
    "fmt"
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

func (cec *ShortenController) AddRecord(c *gin.Context) {
    var addEndpointRequest service.AddRecordParams
    _ = c.ShouldBindJSON(&addEndpointRequest)

    ctxUsername, exists := c.Get("username")
    if exists {
        addEndpointRequest.Username = fmt.Sprintf("%s", ctxUsername)
    }

    srv := service.NewService(cec.dataStore, cec.validator)
    response, err := srv.AddRecord(addEndpointRequest)
    if err != nil {
        internalError(c, err)
        return
    }

    if e, ok := response.Err.(error); ok && e != nil {
        internalError(c, e)
        return
    }

    //fullEndpointURL := utils.GetFullHTTPUrl(c.Request.Host, response.Endpoint, c.Request.TLS != nil)
    c.JSON(200, response.Record)
}

func (cec *ShortenController) GetContent(c *gin.Context) {
    url := c.Request.URL.Path

    username := c.Query("username")

    srv := service.NewService(cec.dataStore, cec.validator)
    response, err := srv.GetContentBySlug(service.GetContentParams{Slug: url, Username: username})
    if err != nil {
        internalError(c, err)
        return
    }

    if e, ok := response.Err.(error); ok && e != nil {
        internalError(c, e)
        return
    }

    //c.DataFromReader(http.StatusOK,
    //    int64(len(response.Username)),
    //    gin.MIMEJSON,
    //    strings.NewReader(response.Username), nil)


    // Check if requester is owner
    c.Get("userToken")


    c.JSON(http.StatusOK, response)
    //c.DataFromReader(http.StatusOK,
    //   int64(len(response.Username)),
    //   gin.MIMEJSON,
    //   strings.NewReader(response.Username), nil)
}

func (cec *ShortenController) UpdateRecord(c *gin.Context) {
    var updateRecordRequest service.UpdateRecordParams
    _ = c.ShouldBindJSON(&updateRecordRequest)

    ctxUsername, exists := c.Get("username")
    if exists {
        updateRecordRequest.Username = fmt.Sprintf("%s", ctxUsername)
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

    //fullEndpointURL := utils.GetFullHTTPUrl(c.Request.Host, response.Endpoint, c.Request.TLS != nil)
    c.JSON(200, response.Record)
}

func (cec *ShortenController) InitialRecord(c *gin.Context) {

    ctxUsername, _ := c.Get("username")
    if ctxUsername.(string) == "" {
        c.AbortWithStatusJSON(400, gin.H{
            "error": "User not found!",
        })
    }

    srv := service.NewService(cec.dataStore, cec.validator)

    randomSlug := srv.GetRandomSlug()

    type M map[string]interface{}
    content, _ := json.Marshal(M{})

    params := service.AddRecordParams{
        Username: ctxUsername.(string),
        Slug:     randomSlug,
        Content:  string(content),
    }
    record, err := srv.AddRecord(params)
    if err != nil {
        panic(err)
    }


    token, exists := c.Get("userToken")
    if exists != true {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "test",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "record": record.Record,
        "token": token,
    })
}
