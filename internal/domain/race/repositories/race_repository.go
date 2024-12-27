package repositories

import (
	"github.com/Casagrande-Lucas/dnd/internal/domain/race/models"
)

type RaceRepository interface {
	GetAllRaces() ([]*models.Race, error)
	GetRaceByID(id uint) (*models.Race, error)
	GetRaceByName(name string) (*models.Race, error)
	CreateRace(race *models.Race) error
	UpdateRace(id uint, race *models.Race) error
	DeleteRace(id uint) error
	AddSubrace(raceID uint, subrace *models.Subrace) error
	RemoveSubrace(raceID uint, subraceID uint) error
	AddTrait(raceID uint, traitID uint) error
	RemoveTrait(raceID uint, traitID uint) error
	SearchRaces(criteria map[string]string) ([]models.Race, error)
}
