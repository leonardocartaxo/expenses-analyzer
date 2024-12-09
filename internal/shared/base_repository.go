package shared

import (
	"gorm.io/gorm"
	"log/slog"
)

type IBaseRepository[Model any] interface {
	GetAll(populateFields []string) ([]Model, error)
	GetById(id string, populateFields []string) (*Model, error)
	Create(model *Model) error
	UpdateById(id string, updates map[string]interface{}) error
	SoftDeleteById(id string) error
	DeleteById(id string) error
	GetNewTransaction() *gorm.DB
	Find(
		conditions []BaseFindCondition,
		populateFields []string,
		orderBy string,
		limit int,
		offset int,
	) ([]Model, error)
}

type BaseRepository[Model any] struct {
	DB     *gorm.DB
	Logger *slog.Logger
}

func (r *BaseRepository[Model]) GetNewTransaction() *gorm.DB {
	return r.DB.Begin()
}

// GetAll retrieves all models with optional relations.
func (r *BaseRepository[Model]) GetAll(relations []string) ([]Model, error) {
	var models []Model
	query := r.DB.Begin()
	for _, relation := range relations {
		query = query.Preload(relation)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

// GetByID retrieves an model by ID with optional relations.
func (r *BaseRepository[Model]) GetById(id string, relations []string) (*Model, error) {
	var model Model
	query := r.DB
	for _, relation := range relations {
		query = query.Preload(relation)
	}
	if err := query.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

// Save creates or updates an model.
func (r *BaseRepository[Model]) Create(model *Model) error {
	return r.DB.Save(model).Error
}

// Update updates a record by id.
func (r *BaseRepository[Model]) UpdateById(id string, partialEntity map[string]interface{}) error {
	return r.DB.Model(new(Model)).Where("id = ?", id).Updates(partialEntity).Error
}

// SoftDelete marks an model as deleted.
func (r *BaseRepository[Model]) SoftDeleteById(id string) error {
	return r.DB.Where("id = ?", id).Delete(new(Model)).Error
}

// Remove deletes an model completely.
func (r *BaseRepository[Model]) Remove(model *Model) error {
	return r.DB.Unscoped().Delete(model).Error
}

// Delete removes models matching id.
func (r *BaseRepository[Model]) DeleteById(id string) error {
	return r.DB.Unscoped().Where("id = ?", id).Delete(new(Model)).Error
}

// Find is a generic method to fetch records with custom conditions.
func (r *BaseRepository[Model]) Find(
	conditions []BaseFindCondition,
	populateFields []string,
	orderBy string,
	limit int,
	offset int,
) ([]Model, error) {
	var models []Model

	// Start query
	query := r.DB

	// Apply conditions
	for _, condition := range conditions {
		query = query.Where(condition.Field+" "+condition.Comparator+" ?", condition.Value)
	}

	// Apply relationships to preload
	for _, field := range populateFields {
		query = query.Preload(field)
	}

	// Apply ordering, if specified
	if orderBy != "" {
		query = query.Order(orderBy)
	}

	// Apply pagination, if specified
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	// Execute query
	err := query.Find(&models).Error
	if err != nil {
		return nil, err
	}

	return models, nil
}
