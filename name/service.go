package name

import (
	"github.com/nahidhasan98/remind-name/logger"
)

type NameService struct {
	repo *repository
}

func NewNameService() *NameService {
	return &NameService{
		repo: newRepository(),
	}
}

func (service *NameService) GetName(id int) (*Name, error) {
	name, err := service.repo.getName(id)
	if err != nil {
		logger.Error("GetName failed for id %d: %v", id, err)
		return nil, err
	}

	logger.Info("GetName succeeded for id %d", id)
	return name, nil
}
