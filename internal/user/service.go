package user

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Save(createDTO *CreateDTO) (*DTO, error) {
	dto, err := s.repository.Save(createDTO)
	if err != nil {
		return nil, err
	}

	return dto.ToDTO(), nil
}

func (s *Service) FindOne(id string) (*DTO, error) {
	dto, err := s.repository.FindOne(id)
	if err != nil {
		return nil, err
	}

	return dto.ToDTO(), nil
}

func (s *Service) UpdateOne(id string, updateDTO *UpdateDTO) (*DTO, error) {
	err := s.repository.UpdateOne(id, updateDTO)
	if err != nil {
		return nil, err
	}

	dto, err := s.FindOne(id)
	if err != nil {
		return nil, err
	}

	return dto, nil
}

func (s *Service) DeleteOne(id string) error {
	err := s.repository.DeleteOne(id)

	return err
}

func (s *Service) All() ([]*Model, error) {
	all, err := s.repository.All()
	if err != nil {
		return nil, err
	}

	return all, nil
}
