package application

import "tpeProgWeb/domain"

type CabinServicesUser struct {
	userRepository domain.CabinUserRepository
}

func NewCabinServiceUser(userRepo domain.CabinUserRepository) *CabinServicesUser {
	return &CabinServicesUser{userRepository: userRepo}
}

func (service *CabinServicesUser) GetCabinByID(id int64) (*domain.Cabin, error) {
	return service.userRepository.GetCabinByID(id)
}

func (service *CabinServicesUser) GetAllCabins() ([]*domain.Cabin, error) {
	return service.userRepository.GetAllCabins()
}

func (service *CabinServicesUser) UpdateCabin(cabin *domain.Cabin) error {
	return service.userRepository.UpdateCabin(cabin)
}
