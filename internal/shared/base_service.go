package shared

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"log/slog"
)

type BaseService[Model any, Dto any, CreateDto any, UpdateDto any] struct {
	Repo           IBaseRepository[Model]
	EntityToDto    func(Model) Dto
	PopulateFields []string
	DtoFactory     func() Dto // Factory function to create Dto instances
	Logger         *slog.Logger
}

// Create creates a new model and returns its DTO.
func (s *BaseService[Model, Dto, CreateDto, UpdateDto]) Create(dto *CreateDto) (Dto, error) {
	var model Model
	if err := mapStruct(dto, &model); err != nil {
		return s.DtoFactory(), fmt.Errorf("failed to map CreateDto to Model: %w", err)
	}

	if err := s.Repo.Create(&model); err != nil {
		return s.DtoFactory(), fmt.Errorf("failed to save model: %w", err)
	}

	return s.EntityToDto(model), nil
}

// GetEntity retrieves a model by ID, optionally with related fields.
func (s *BaseService[Model, Dto, CreateDto, UpdateDto]) GetEntity(id string) (*Model, error) {
	model, err := s.Repo.GetById(id, s.PopulateFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("model not found: %w", err)
		}
		return nil, err
	}
	return model, nil
}

// Get retrieves a model by ID and converts it to a DTO.
func (s *BaseService[Model, Dto, CreateDto, UpdateDto]) GetById(id string) (Dto, error) {
	model, err := s.GetEntity(id)
	if err != nil {
		return s.DtoFactory(), err
	}
	return s.EntityToDto(*model), nil
}

// GetAll retrieves all models and converts them to DTOs.
func (s *BaseService[Model, Dto, CreateDto, UpdateDto]) GetAll() ([]Dto, error) {
	models, err := s.Repo.GetAll(s.PopulateFields)
	if err != nil {
		return nil, err
	}

	var dtos []Dto
	for _, model := range models {
		dtos = append(dtos, s.EntityToDto(model))
	}
	return dtos, nil
}

// Update updates a dto by ID and returns its updated DTO.
func (s *BaseService[Model, Dto, CreateDto, UpdateDto]) Update(id string, dto *UpdateDto) (Dto, error) {
	updates := map[string]interface{}{}
	if err := mapStruct(dto, &updates); err != nil {
		return s.DtoFactory(), fmt.Errorf("failed to map UpdateDto to map: %w", err)
	}

	if err := s.Repo.UpdateById(id, updates); err != nil {
		return s.DtoFactory(), fmt.Errorf("failed to update model: %w", err)
	}

	return s.GetById(id)
}

// SoftDelete soft-deletes an model by ID.
func (s *BaseService[Model, Dto, CreateDto, UpdateDto]) SoftDelete(id string) (bool, error) {
	err := s.Repo.SoftDeleteById(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Delete hard-deletes an model by ID.
func (s *BaseService[Model, Dto, CreateDto, UpdateDto]) Delete(id string) (bool, error) {
	err := s.Repo.DeleteById(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Find uses the repository's Find function to retrieve records
func (s *BaseService[Model, Dto, CreateDto, UpdateDto]) Find(
	conditions []BaseFindCondition,
	populateFields []string,
	orderBy string,
	limit int,
	offset int,
) ([]Dto, error) {
	// Call the repository's Find function
	models, err := s.Repo.Find(conditions, populateFields, orderBy, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert models to DTOs
	dtos := make([]Dto, len(models))
	for i, entity := range models {
		dtos[i] = s.EntityToDto(entity)
	}

	return dtos, nil
}

func mapStruct(input any, output any) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		TagName:  "json",
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}
