package application

import "tpeProgWeb/domain"

type CabinServicesADM struct {
	userRepository domain.CabinADMRepository
}

func NewCabinServiceADM(userRepo domain.CabinADMRepository) *CabinServicesADM {
	return &CabinServicesADM{userRepository: userRepo}
}

func (service *CabinServicesADM) CreateCabin(cabin *domain.Cabin) error {
	return service.userRepository.CreateCabin(cabin)
}

func (service *CabinServicesADM) GetCabinByID(id int64) (*domain.Cabin, error) {
	return service.userRepository.GetCabinByID(id)
}

func (service *CabinServicesADM) GetAllCabins() ([]*domain.Cabin, error) {
	return service.userRepository.GetAllCabins()
}

func (service *CabinServicesADM) UpdateCabin(cabin *domain.Cabin) error {
	return service.userRepository.UpdateCabin(cabin)
}
