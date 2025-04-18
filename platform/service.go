package platform

type platformService struct {
	repo *repository
}

func NewPlatformService() *platformService {
	return &platformService{
		repo: newRepository(),
	}
}

func (service *platformService) GetAllPlatforms() ([]Platform, error) {
	platforms, err := service.repo.getAllPlatforms()
	if err != nil {
		return nil, err
	}

	return platforms, nil
}

func (service *platformService) GetPlatformDetailsByName(name string) (*Platform, error) {
	platform, err := service.repo.getPlatformDetailsByName(name)
	if err != nil {
		return nil, err
	}

	return platform, nil
}
