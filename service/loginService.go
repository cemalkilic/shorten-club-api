package service

import (
    "errors"
    "github.com/cemalkilic/shorten-backend/database"
    "github.com/cemalkilic/shorten-backend/models"
    "github.com/cemalkilic/shorten-backend/utils/validator"
    "golang.org/x/crypto/bcrypt"
    "log"
    "time"
)

type LoginService interface {
    IsValidCredentials(user models.User) bool
    Signup(user models.User) error
}

type loginInformation struct {
    db database.UserStore
    v *validator.CustomValidator
}

func DBLoginService(db database.UserStore, v *validator.CustomValidator) LoginService {
    return &loginInformation{
        db: db,
        v: v,
    }
}

func (info *loginInformation) IsValidCredentials(credentials models.User) bool {
    if err := info.v.ValidateStruct(credentials); err != nil {
        return false
    }

    user, err := info.db.SelectByUsername(credentials.Username)
    if err != nil {
        return false
    }

    return comparePasswords(user.Password, credentials.Password)
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
    hashPwdByte := []byte(hashedPwd)
    plainPwdByte := []byte(plainPwd)

    err := bcrypt.CompareHashAndPassword(hashPwdByte, plainPwdByte)
    if err != nil {
        log.Println(err)
        return false
    }

    return true
}

func (info *loginInformation) Signup(credentials models.User) error {

    if err := info.v.ValidateStruct(credentials); err != nil {
        return err
    }

    // check user exists
    response, err := info.db.SelectByUsername(credentials.Username)
    if err != nil {
        return err
    }

    if response.Username != "" {
        return errors.New("user already exists")
    }

    // Hash & store the password
    credentials.Password = hashAndSaltPassword(credentials.Password)
    credentials.CreatedAt = time.Now()

    err = info.db.InsertUser(credentials)
    if err != nil {
        return err
    }

    return nil
}

func hashAndSaltPassword(password string) string {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
    if err != nil {
        log.Println(err)
    }

    return string(hash)
}
