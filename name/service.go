package name

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
		return nil, err
	}

	return name, nil
}
