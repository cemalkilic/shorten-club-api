package service

import (
    "errors"
    "fmt"
    "github.com/cemalkilic/shorten-backend/database"
    "github.com/cemalkilic/shorten-backend/models"
    "github.com/cemalkilic/shorten-backend/utils/validator"
    "math/rand"
    "strings"
    "time"
)

type jsonService struct {
    db database.DataStore
    validate *validator.CustomValidator
}

func NewService(db database.DataStore, v *validator.CustomValidator) *jsonService {
    return &jsonService{
        db: db,
        validate: v,
    }
}

func (srv *jsonService) GetContentBySlug(params GetContentParams) (GetResponse, error) {

    // Terminate the request if the input is not valid
    if err := srv.validate.ValidateStruct(params); err != nil {
        return GetResponse{}, err
    }

    slug := strings.Trim(params.Slug, "/")
    if len(slug) == 0 {
        return GetResponse{}, errors.New("empty URI given")
    }

    //urlParts := strings.Split(slug, "/")
    //username := urlParts[0]
    //if len(urlParts[1:]) < 1 {
    //  return GetResponse{}, errors.New("URI given without custom endpoint")
    //}

    //userEndpoint := strings.Join(urlParts[1:], "/")

    record, err := srv.db.SelectBySlug(slug)
    if err != nil {
        return GetResponse{}, err
    }

    fmt.Printf("\n%#v\n\n", record)

    if record.ID == 0 {
        // not found the custom endpoint
        return GetResponse{}, errors.New("404: Not Found")
    }

    // at its simplest
    permissions := map[string]bool{"readContent": true, "updateContent": false}
    if params.Username != "" {
        permissions["updateContent"] = true
    }

    return GetResponse{
        Username: record.Username,
        Slug:     record.Slug,
        Content:  record.Content,
        Permissions: permissions,
        Err:      nil,
    }, nil
}

func (srv *jsonService) AddRecord(params AddRecordParams) (AddRecordResponse, error) {
    // Prepend with a slash to behave it like a uri
    //if !strings.HasPrefix(params.Endpoint, "/") {
    //    params.Endpoint = "/" + params.Endpoint
    //}

    // Terminate the request if the input is not valid
    if err := srv.validate.ValidateStruct(params); err != nil {
       return AddRecordResponse{}, err
    }

    // Trim the slashes after validation :/ That's way easier than custom validation
    params.Slug = strings.Trim(params.Slug, "/")

    // Use the default username if not exists in the params
    username := params.Username
    if username == "" {
        username = "guest"
    }

    // Make sure the same endpoint does not already exist
    response, err := srv.db.SelectBySlug(params.Slug)
    if err != nil {
        return AddRecordResponse{}, err
    }

    if response.ID != 0 {
        return AddRecordResponse{}, errors.New("endpoint already exists")
    }

    recordObj := models.Record{
        Username:   username,
        Slug:       params.Slug, // TODO
        Content:    params.Content,
    }

    fmt.Printf("\n%#v\n\n", recordObj)

    err = srv.db.Insert(recordObj)
    if err != nil {
        return AddRecordResponse{}, err
    }


    return AddRecordResponse{
        Record:   recordObj,
        Err:      nil,
    }, nil
}

func (srv *jsonService) UpdateRecord(params UpdateRecordParams) (UpdateRecordResponse, error) {
    // Prepend with a slash to behave it like a uri
    //if !strings.HasPrefix(params.Endpoint, "/") {
    //    params.Endpoint = "/" + params.Endpoint
    //}

    // Terminate the request if the input is not valid
    if err := srv.validate.ValidateStruct(params); err != nil {
        return UpdateRecordResponse{}, err
    }

    // Trim the slashes after validation :/ That's way easier than custom validation
    params.Slug = strings.Trim(params.Slug, "/")

    // Use the default username if not exists in the params
    username := params.Username
    if username == "" {
        username = "guest"
    }

    // Make sure record to update exists
    response, err := srv.db.SelectBySlug(params.Slug)
    if err != nil {
        return UpdateRecordResponse{}, err
    }

    if response.ID == 0 {
        return UpdateRecordResponse{}, errors.New("endpoint not exists")
    }

    recordObj := models.Record{
        Username:   username,
        Slug:       params.Slug, // TODO
        Content:    params.Content,
    }

    fmt.Printf("\nupdate: %#v\n\n", recordObj)

    err = srv.db.UpdateBySlug(recordObj)
    if err != nil {
        return UpdateRecordResponse{}, err
    }

    return UpdateRecordResponse{
        Record:   recordObj,
        Err:      nil,
    }, nil
}


func (srv *jsonService) GetRandomSlug() string {
    var randSlug string
    for {
        randSlug = randomSlug()
        record, _ := srv.db.SelectBySlug(randSlug)
        if record.ID == 0 {
            break
        }
    }
    return randSlug
}


func randomSlug() string {
    return StringWithCharset(6, "0123456789")
}

func StringWithCharset(length int, charset string) string {
    b := make([]byte, length)
    var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := range b {
        b[i] = charset[seededRand.Intn(len(charset))]
    }
    return string(b)
}
