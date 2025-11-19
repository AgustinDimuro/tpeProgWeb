package application

import (
	"errors"
	"tpeProgWeb/domain"
)

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

func (service *CabinServicesUser) Authenticate(id int64, password string) (*domain.Cabin, error) {
	cabin, err := service.userRepository.GetCabinByID(id)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Paso 2: Aplicar regla de negocio (comparación)
	// Aquí en el futuro cambiaremos "==" por bcrypt.CompareHashAndPassword
	if cabin.Password != password {
		return nil, errors.New("credenciales inválidas")
	}

	return cabin, nil
}
