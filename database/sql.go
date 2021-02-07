package database

import (
    "database/sql"
    "encoding/json"
    "errors"
    "github.com/cemalkilic/shorten-backend/models"
    _ "github.com/go-sql-driver/mysql"
    "log"
    "time"
)

type sqlDatabase struct {
    db *sql.DB
}

func GetSQLDataStore(db *sql.DB) DataStore {
    return &sqlDatabase{
        db: db,
    }
}

func GetSQLUserStore(db *sql.DB) UserStore {
    return &sqlDatabase{
        db: db,
    }
}

func (s *sqlDatabase) Insert(record models.Record) error {
    insertRecordSQL := `INSERT INTO records(username, slug, type, content) VALUES (?, ?, ?, ?)`
    statement, err := s.db.Prepare(insertRecordSQL)
    if err != nil {
        log.Fatal(err)
        return err
    }

    content, _ := json.Marshal(record.Content)
    _, err = statement.Exec(record.Username, record.Slug, record.Type, content)
    if err != nil {
        return err
    }

    return nil
}

func (s *sqlDatabase) Select(username string, slug string) (models.Record, error) {
    rows, err := s.db.Query("SELECT * FROM records WHERE `username` = ? AND `slug` = ?", username, slug)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var idResp int64
    var usernameResp string
    var slugResp string
    var typeResp string
    var contentResp interface{}

    for rows.Next() { // Iterate and fetch the records from result cursor
        rows.Scan(&idResp, &usernameResp, &slugResp, &typeResp, &contentResp)
    }
    if err := rows.Err(); err != nil {
        return models.Record{}, err
    }

    return models.Record{
        ID:         idResp,
        Username:   usernameResp,
        Slug:       slugResp,
        Type:       typeResp,
        Content:    contentResp,
    }, nil
}

func (s *sqlDatabase) SelectByID(id int) (models.Record, error) {
    // TODO
    return models.Record{}, errors.New("in SelectByID")
}

func (s *sqlDatabase) SelectBySlug(slug string) (models.Record, error) {
    rows, err := s.db.Query("SELECT * FROM records WHERE `slug` = ?", slug)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var idResp int64
    var usernameResp string
    var slugResp string
    var typeResp string
    var contentResp string
    var expire_at time.Time
    var created_at time.Time
    var updated_at time.Time

    for rows.Next() { // Iterate and fetch the records from result cursor
        rows.Scan(&idResp, &usernameResp, &slugResp, &typeResp, &contentResp, &expire_at, &created_at, &updated_at)
    }
    if err := rows.Err(); err != nil {
        return models.Record{}, err
    }

    var result interface{}
    _ = json.Unmarshal([]byte(contentResp), &result)

    return models.Record{
        ID:         idResp,
        Username:   usernameResp,
        Slug:       slugResp,
        Type:       typeResp,
        Content:    result,
    }, nil
}

func (s *sqlDatabase) UpdateBySlug(record models.Record) error {
    content, _ := json.Marshal(record.Content)
    slug := record.Slug

    updateRecordSQL := `UPDATE records SET content = ? WHERE slug = ?`
    statement, err := s.db.Prepare(updateRecordSQL)
    if err != nil {
        log.Fatal(err)
        return err
    }

    _, err = statement.Exec(content, slug)
    if err != nil {
        return err
    }

    return nil
}

func (s *sqlDatabase) SelectAllByUser(username string) ([]models.Record, error) {
    // TODO
    return []models.Record{}, errors.New("in SelectAllByUser")
}

func (s *sqlDatabase) Delete(id int) error {
    // TODO
    return errors.New("in Delete")
}

func (s *sqlDatabase) InsertUser(user models.User) error {
    insertUserSql := `INSERT INTO users(username, password, createdAt) VALUES (?, ?, ?)`
    statement, err := s.db.Prepare(insertUserSql)
    if err != nil {
        log.Fatal(err)
        return err
    }

    _, err = statement.Exec(user.Username, user.Password, user.CreatedAt)
    if err != nil {
        return err
    }

    return nil
}

func (s *sqlDatabase) SelectByUsername(uname string) (models.User, error) {
    rows, err := s.db.Query(`SELECT * FROM users WHERE username = ? LIMIT 1`, uname)
    if err != nil {
        return models.User{}, err
    }

    defer rows.Close()

    var id string
    var username string
    var password string
    var createdAt time.Time

    for rows.Next() {
        _ = rows.Scan(&id, &username, &password, &createdAt)

    }

    if err := rows.Err(); err != nil {
        return models.User{}, err
    }

    return models.User{
        Username: username,
        Password: password,
        CreatedAt: createdAt,
    }, nil
}
