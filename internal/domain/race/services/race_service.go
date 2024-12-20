package services

import "github.com/Casagrande-Lucas/dnd/internal/domain/race/entities"

type RaceService interface {
	ListRaces() ([]entities.Race, error)
	GetRaceDetails(id uint) (*entities.Race, error)
	RegisterRace(race *entities.Race) error
	UpdateRaceInfo(id uint, race *entities.Race) error
	RemoveRace(id uint) error
	AddSubraceToRace(raceID uint, subrace *entities.Subrace) error
	DetachSubraceFromRace(raceID uint, subraceID uint) error
	AssignTraitToRace(raceID uint, traitID uint) error
	UnassignTraitFromRace(raceID uint, traitID uint) error
	FindRaces(criteria map[string]string) ([]entities.Race, error)
}
