package services

import (
	"github.com/Casagrande-Lucas/dnd/internal/domain/race/models"
	"github.com/google/uuid"
)

type RaceService interface {
	ListRaces() ([]*models.Race, error)
	GetRaceDetails(id uuid.UUID) (*models.Race, error)
	RegisterRace(race *models.Race) error
	UpdateRaceInfo(id uuid.UUID, race *models.Race) error
	RemoveRace(id uuid.UUID) error
	AddSubraceToRace(raceID uuid.UUID, subrace *models.Subrace) error
	DetachSubraceFromRace(raceID uuid.UUID, subraceID uuid.UUID) error
	AssignTraitToRace(raceID uuid.UUID, traitID uuid.UUID) error
	UnassignTraitFromRace(raceID uuid.UUID, traitID uuid.UUID) error
	FindRaces(criteria map[string]string) ([]models.Race, error)
}
