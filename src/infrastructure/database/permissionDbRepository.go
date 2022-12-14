package database

import (
	"context"
	"go-as/src/domain/permission"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PermissionDbRepository struct {
	db *gorm.DB
}

func (repo *PermissionDbRepository) Save(ctx context.Context, permission permission.Permission) error {
	db := repo.db.WithContext(ctx)
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&permission)
	return result.Error
}

func (repo *PermissionDbRepository) FindByNames(ctx context.Context, permissionNames []string) ([]permission.Permission, error) {
	var foundPermissions []permission.Permission
	db := repo.db.WithContext(ctx)
	result := db.Where("name IN ?", permissionNames).Find(&foundPermissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return foundPermissions, nil
}

func NewPermissionDbRepository(db *gorm.DB) *PermissionDbRepository {
	repo := PermissionDbRepository{
		db: db,
	}
	return &repo
}
