package db

import (
	"context"
	"database/sql"
	db "tpeProgWeb/db/sqlc"
	"tpeProgWeb/domain"
)

var ctx = context.Background()
var dbQueries *db.Queries

type DBCabinRepository struct {
	db *sql.DB
}

func NewDBCabinRepository(db *sql.DB) *DBCabinRepository {
	return &DBCabinRepository{db: db}
}

// En db/db_cabin_repository.go

func (repo *DBCabinRepository) CreateCabin(cabin *domain.Cabin) error {
	dbQueries = db.New(repo.db)

	// Si no viene rol, forzamos 'user' por defecto desde el código Go
	roleToSave := cabin.Role
	if roleToSave == "" {
		roleToSave = "user"
	}

	newCabinDB, err := dbQueries.CreateCabin(ctx,
		db.CreateCabinParams{
			EmailContact: cabin.EmailContact,
			PhoneContact: cabin.PhoneContact,
			Password:     cabin.Password,
			Role:         roleToSave, // <--- Mapeamos el Rol
		},
	)

	// Actualizamos el ID y el Rol en el objeto de dominio original
	if err == nil {
		cabin.ID = int64(newCabinDB.ID)
		cabin.Role = newCabinDB.Role
	}

	return err
}

func (repo *DBCabinRepository) GetCabinByID(id int64) (*domain.Cabin, error) {
	dbQueries = db.New(repo.db)
	cabinDB, err := dbQueries.GetCabin(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	cabin := &domain.Cabin{
		ID:           int64(cabinDB.ID),
		EmailContact: cabinDB.EmailContact,
		PhoneContact: cabinDB.PhoneContact,
		Password:     cabinDB.Password,
		Role:         cabinDB.Role, // <--- Mapeamos el Rol
	}
	return cabin, nil
}

func (repo *DBCabinRepository) GetAllCabins() ([]*domain.Cabin, error) {
	dbQueries = db.New(repo.db)
	cabinsDB, err := dbQueries.ListCabins(ctx)
	if err != nil {
		return nil, err
	}
	var cabins []*domain.Cabin
	for _, cabinDB := range cabinsDB {
		cabin := &domain.Cabin{
			ID:           int64(cabinDB.ID),
			EmailContact: cabinDB.EmailContact,
			PhoneContact: cabinDB.PhoneContact,
			Password:     cabinDB.Password,
			Role:         cabinDB.Role, // <--- Mapeamos el Rol
		}
		cabins = append(cabins, cabin)
	}
	return cabins, nil
}

func (repo *DBCabinRepository) UpdateCabin(cabin *domain.Cabin) error {
	dbQueries = db.New(repo.db)
	_, err := dbQueries.UpdateCabin(ctx,
		db.UpdateCabinParams{
			ID:           int32(cabin.ID),
			EmailContact: cabin.EmailContact,
			PhoneContact: cabin.PhoneContact,
			Password:     cabin.Password,
			Role:         cabin.Role, // <--- Mapeamos el Rol
		},
	)
	return err
}
