package collector

import "log/slog"

type Service struct {
	repository *Repository
	l          *slog.Logger
}

func NewService(repository *Repository, l *slog.Logger) *Service {
	return &Service{repository: repository, l: l}
}

func (s *Service) Save(createDTO *CreateDTO) (*DTO, error) {
	model, err := s.repository.Save(createDTO)
	if err != nil {
		return nil, err
	}

	return model.ToDTO(), nil
}

func (s *Service) FindOne(id string) (*DTO, error) {
	model, err := s.repository.FindOne(id)
	if err != nil {
		return nil, err
	}

	return model.ToDTO(), nil
}

func (s *Service) UpdateOne(id string, updateDTO *UpdateDTO) (*DTO, error) {
	err := s.repository.UpdateOne(id, updateDTO)
	if err != nil {
		return nil, err
	}

	model, err := s.FindOne(id)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (s *Service) DeleteOne(id string) error {
	err := s.repository.DeleteOne(id)

	return err
}

func (s *Service) All() ([]*DTO, error) {
	var models Models
	models, err := s.repository.All()
	if err != nil {
		return nil, err
	}

	return models.ToDTO(), nil
}

func (s *Service) Find(options FindOptions) ([]*DTO, error) {
	var models Models
	models, err := s.repository.Find(options)
	if err != nil {
		return nil, err
	}

	return models.ToDTO(), nil
}
