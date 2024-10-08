package user

import (
	"gorm.io/gorm"
	"log/slog"
)

type Repository struct {
	db *gorm.DB
	l  *slog.Logger
}

func NewRepository(db *gorm.DB, l *slog.Logger) *Repository {
	return &Repository{db: db, l: l}
}

func (r *Repository) Save(user *CreateDTO) (*Model, error) {
	model := &Model{Name: user.Name}
	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (r *Repository) UpdateOne(id string, user *UpdateDTO) error {
	err := r.db.Model(&Model{}).Where("ID = ?", id).Updates(user).Error

	return err
}

func (r *Repository) FindOne(id string) (*Model, error) {
	model := &Model{}
	if err := r.db.Where("ID = ?", id).First(&model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (r *Repository) DeleteOne(id string) error {
	model := &Model{}
	err := r.db.Where("ID = ?", id).Delete(&model).Error

	return err
}

func (r *Repository) All() ([]*Model, error) {
	var models []*Model
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}
