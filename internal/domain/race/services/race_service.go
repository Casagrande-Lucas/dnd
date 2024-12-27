package services

import "github.com/Casagrande-Lucas/dnd/internal/domain/race/models"

type RaceService interface {
	ListRaces() ([]*models.Race, error)
	GetRaceDetails(id uint) (*models.Race, error)
	RegisterRace(race *models.Race) error
	UpdateRaceInfo(id uint, race *models.Race) error
	RemoveRace(id uint) error
	AddSubraceToRace(raceID uint, subrace *models.Subrace) error
	DetachSubraceFromRace(raceID uint, subraceID uint) error
	AssignTraitToRace(raceID uint, traitID uint) error
	UnassignTraitFromRace(raceID uint, traitID uint) error
	FindRaces(criteria map[string]string) ([]models.Race, error)
}
