package domain

type Cabin struct {
	ID           int64
	EmailContact string
	PhoneContact string
	Password     string
	Role         string
}

type CabinADMRepository interface {
	CreateCabin(cabin *Cabin) error
	GetCabinByID(id int64) (*Cabin, error)
	GetAllCabins() ([]*Cabin, error)
	UpdateCabin(cabin *Cabin) error
}

type CabinUserRepository interface {
	GetCabinByID(id int64) (*Cabin, error)
	GetAllCabins() ([]*Cabin, error)
	UpdateCabin(cabin *Cabin) error
}
