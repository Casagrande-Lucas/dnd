package repositories

import (
	"github.com/Casagrande-Lucas/dnd/internal/domain/race/models"
	"github.com/google/uuid"
)

type RaceRepository interface {
	GetAllRaces() ([]*models.Race, error)
	GetRaceByID(id uuid.UUID) (*models.Race, error)
	GetRaceByName(name string) (*models.Race, error)
	CreateRace(race *models.Race) error
	UpdateRace(id uuid.UUID, race *models.Race) error
	DeleteRace(id uuid.UUID) error
	AddSubrace(raceID uuid.UUID, subrace *models.Subrace) error
	RemoveSubrace(raceID uuid.UUID, subraceID uuid.UUID) error
	AddTrait(raceID uuid.UUID, traitID uuid.UUID) error
	RemoveTrait(raceID uuid.UUID, traitID uuid.UUID) error
	SearchRaces(criteria map[string]string) ([]models.Race, error)
}
