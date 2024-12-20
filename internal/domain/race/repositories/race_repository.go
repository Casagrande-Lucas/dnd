package repositories

import (
	"github.com/Casagrande-Lucas/dnd/internal/domain/race/entities"
)

type RaceRepository interface {
	GetAllRaces() ([]entities.Race, error)
	GetRaceByID(id uint) (*entities.Race, error)
	GetRaceByName(name string) (*entities.Race, error)
	CreateRace(race *entities.Race) error
	UpdateRace(id uint, race *entities.Race) error
	DeleteRace(id uint) error
	AddSubrace(raceID uint, subrace *entities.Subrace) error
	RemoveSubrace(raceID uint, subraceID uint) error
	AddTrait(raceID uint, traitID uint) error
	RemoveTrait(raceID uint, traitID uint) error
	SearchRaces(criteria map[string]string) ([]entities.Race, error)
}
