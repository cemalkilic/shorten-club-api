package service

import (
    "github.com/cemalkilic/shorten-backend/database"
)

type permissionService struct {
    db database.DataStore
}


func NewPermissionService(db database.DataStore) *permissionService {
    return &permissionService{
        db: db,
    }
}

func (srv *permissionService) GetPermissionsBySlug(username string, slug string) (map[string]bool, error) {
    permissions := map[string]bool{"readContent": true, "updateContent": false} // defaults

    record, err := srv.db.SelectBySlug(slug)
    if err != nil {
        return permissions, err
    }

    if username == record.Username {
        permissions["updateContent"] = true
        return permissions, nil
    }

    return permissions, nil
}
