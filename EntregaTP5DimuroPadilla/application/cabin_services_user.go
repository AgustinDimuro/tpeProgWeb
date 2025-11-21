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
	// 1. Buscamos la cabaña en la BD (usando la capa de infraestructura)
	cabin, err := service.userRepository.GetCabinByID(id)
	if err != nil {
		// Si hay error (ej: no existe ID), devolvemos error genérico de seguridad
		return nil, errors.New("credenciales inválidas")
	}

	// 2. Validamos la contraseña
	// NOTA: Aquí es donde luego pondremos bcrypt.CompareHashAndPassword
	if cabin.Password != password {
		return nil, errors.New("credenciales inválidas")
	}

	// 3. Si todo ok, devolvemos la cabaña (que incluye el campo Role)
	return cabin, nil
}
