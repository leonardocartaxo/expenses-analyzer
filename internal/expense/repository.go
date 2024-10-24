package expense

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
	err := r.db.First(&Model{}, id).Updates(user).Error

	return err
}

func (r *Repository) FindOne(id string) (*Model, error) {
	model := &Model{}
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (r *Repository) DeleteOne(id string) error {
	model := &Model{}
	err := r.db.First(&model, id).Delete(&model).Error

	return err
}

func (r *Repository) All() ([]*Model, error) {
	var models []*Model
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

func (r *Repository) Find(options FindOptions) ([]*Model, error) {
	if options.Name == "" {
		options.Name = "%"
	}
	if options.Start == "" {
		options.Start = "1970-01-01 00:00:00"
	}
	if options.End == "" {
		options.End = "2070-01-01 00:00:00"
	}
	var models []*Model
	if err := r.db.Find(
		&models,
		"name LIKE ? AND created_at > ? AND created_at < ?",
		options.Name,
		options.Start,
		options.End,
	).Error; err != nil {
		return nil, err
	}

	return models, nil
}
