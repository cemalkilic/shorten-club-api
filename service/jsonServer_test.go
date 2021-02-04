package service_test

import (
    "fmt"
    "github.com/cemalkilic/shorten-backend/models"
    "github.com/cemalkilic/shorten-backend/service"
    "github.com/cemalkilic/shorten-backend/utils/validator"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "net/http"
    "strings"
    "testing"
)

// Proxy the real datastore calls by the mock
type MockDB struct {
    mock.Mock
    InsertFunc          func(endpoint models.CustomEndpoint) error
    SelectFunc          func(username string, uri string) (models.CustomEndpoint, error)
    SelectByIDFunc      func(id int) (models.CustomEndpoint, error)
    SelectAllByUserFunc func(username string) ([]models.CustomEndpoint, error)
    DeleteFunc          func(id int) error
}

func (db *MockDB) Insert(endpoint models.CustomEndpoint) error {
    db.Called(endpoint)
    return db.InsertFunc(endpoint)
}

func (db *MockDB) Select(username string, uri string) (models.CustomEndpoint, error) {
    db.Called(username, uri)
    return db.SelectFunc(username, uri)
}

func (db *MockDB) SelectByID(id int) (models.CustomEndpoint, error) {
    db.Called(id)
    return db.SelectByIDFunc(id)
}

func (db *MockDB) SelectAllByUser(username string) ([]models.CustomEndpoint, error) {
    db.Called(username)
    return db.SelectAllByUserFunc(username)
}

func (db *MockDB) Delete(id int) error {
    db.Called(id)
    return db.DeleteFunc(id)
}

func TestParseUsername(t *testing.T) {

    username := "testing"
    endpoint := "my/endpoint"
    expectedContent := `{"testing":"jsonContent"}`
    userEndpoint := fmt.Sprintf("/%s/%s", username, endpoint)

    dbMock := &MockDB{
        SelectFunc: func(username string, uri string) (models.CustomEndpoint, error) {
            return models.CustomEndpoint{
                ID:         1,
                Username:   username,
                URI:        uri,
                Content:    expectedContent,
                StatusCode: 200,
            }, nil
        },
    }

    dbMock.On("Select", username, endpoint).Once()

    validate := validator.NewValidator()
    jsonService := service.NewService(dbMock, validate)

    response, err := jsonService.GetCustomEndpoint(service.GetContentParams{Endpoint: userEndpoint})
    if err != nil {
        t.Fatalf("Error returned: %v", err)
    }

    dbMock.AssertExpectations(t)

    if response.StatusCode != http.StatusOK {
        t.Fatalf("handler returned wrong status code: got %v want %v", response.StatusCode, http.StatusOK)
    }

    assert.Equal(t, expectedContent, response.Content)

    if response.Err != nil {
        t.Fatalf("handler returned wrong status code: got %v want %v", response.Err, nil)
    }
}

func TestOnlyUsernameShouldFail(t *testing.T) {

    username := "onlyUsername"
    userEndpoint := fmt.Sprintf("/%s", username)

    dbMock := MockDB{}
    validate := validator.NewValidator()

    jsonService := service.NewService(&dbMock, validate)

    _, err := jsonService.GetCustomEndpoint(service.GetContentParams{Endpoint: userEndpoint})
    if err == nil {
        t.Fatalf("Error not returned for non-uri")
    }
}

func TestEmptyUsernameShouldFail(t *testing.T) {

    username := ""
    userEndpoint := fmt.Sprintf("/%s", username)

    dbMock := MockDB{}
    validate := validator.NewValidator()

    jsonService := service.NewService(&dbMock, validate)

    _, err := jsonService.GetCustomEndpoint(service.GetContentParams{Endpoint: userEndpoint})
    if err == nil {
        t.Fatalf("Error not returned for empty username")
    }
}

func TestJsonService_AddEndpoint_WithEmptyUsername(t *testing.T) {
    expectedUsername := "guest"

    username := ""
    endpoint := "myEndpoint"
    content := `{"json":"content"}`
    statusCode := 200

    dbMock := &MockDB{
        SelectFunc: func(username string, uri string) (models.CustomEndpoint, error) {
            return models.CustomEndpoint{}, nil
        },
        InsertFunc: func(endpoint models.CustomEndpoint) error {
            return nil
        },
    }

    dbMock.On("Insert", models.CustomEndpoint{
        Username:   expectedUsername,
        URI:        endpoint,
        Content:    content,
        StatusCode: 200,
    }).Once()

    dbMock.On("Select", expectedUsername, endpoint).Once()
    validate := validator.NewValidator()

    jsonService := service.NewService(dbMock, validate)

    response, err := jsonService.AddRecord(service.AddRecordParams{
        Username:   username,
        Endpoint:   endpoint,
        Content:    content,
        StatusCode: statusCode,
    })

    if err != nil {
        t.Fatalf("Error returned for empty username, it should create a random username")
    }

    dbMock.AssertExpectations(t)

    // Returned endpoint must have a random username prepended
    firstSlashIndex := strings.Index(response.Endpoint, "/")
    if firstSlashIndex == -1 {
        t.Fatalf("No slash found in the created endpoint!")
    }

    if response.Err != nil {
        t.Fatalf("Error should not be returned on response!")
    }

    usernameResp := response.Endpoint[:firstSlashIndex]
    endpointResp := response.Endpoint[firstSlashIndex +1 :]

    assert.Equal(t, expectedUsername, usernameResp, "Created username must be the same")
    assert.Equal(t, endpoint, endpointResp, "Created endpoint must be the same with the given")
}

func TestJsonService_AddEndpoint_WithValidUsername(t *testing.T) {
    username := "myUser"
    endpoint := "my/test/endpoint"
    content := `{"json":"content"}`
    statusCode := 200

    dbMock := MockDB{
        SelectFunc: func(username string, uri string) (models.CustomEndpoint, error) {
            return models.CustomEndpoint{
                ID:         0,
                Username:   username,
                URI:        uri,
                Content:    content,
                StatusCode: 200,
            }, nil
        },
        InsertFunc: func(endpoint models.CustomEndpoint) error {
            return nil
        },
    }

    dbMock.On("Insert", models.CustomEndpoint{
        Username:   username,
        URI:        endpoint,
        Content:    content,
        StatusCode: 200,
    }).Once()

    dbMock.On("Select", username, endpoint).Once()
    validate := validator.NewValidator()

    jsonService := service.NewService(&dbMock, validate)

    response, err := jsonService.AddRecord(service.AddRecordParams{
        Username:   username,
        Endpoint:   endpoint,
        Content:    content,
        StatusCode: statusCode,
    })

    if err != nil {
        t.Fatalf("Error returned for empty username, it should create a random username")
    }

    if response.Err != nil {
        t.Fatalf("Error should not be returned on response!")
    }

    dbMock.AssertExpectations(t)

    expectedCreatedEndpoint := fmt.Sprintf("%s/%s", username, endpoint)
    assert.Equal(t, expectedCreatedEndpoint, response.Endpoint, "Created endpoint must be equal to the given")
}
