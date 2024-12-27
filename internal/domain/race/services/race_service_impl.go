package services

import (
	"errors"
	"fmt"

	"github.com/Casagrande-Lucas/dnd/internal/domain/race/models"
	"github.com/Casagrande-Lucas/dnd/internal/domain/race/repositories"
	"github.com/Casagrande-Lucas/dnd/pkg/failure"
	"github.com/google/uuid"
)

// raceServiceImpl is the concrete implementation of RaceService.
type raceServiceImpl struct {
	repo repositories.RaceRepository
}

// NewRaceService creates a new instance of raceServiceImpl.
func NewRaceService(repo repositories.RaceRepository) RaceService {
	return &raceServiceImpl{
		repo: repo,
	}
}

func (s *raceServiceImpl) ListRaces() ([]*models.Race, error) {
	races, err := s.repo.GetAllRaces()
	if err != nil {
		return nil, failure.NewError(failure.ErrorInternalServer, fmt.Errorf("failed to get list races: %w", err))
	}
	return races, nil
}

func (s *raceServiceImpl) GetRaceDetails(id uuid.UUID) (*models.Race, error) {
	if id == uuid.Nil {
		return nil, failure.NewError(failure.ErrorBadRequest, fmt.Errorf("invalid race ID: %s", id.String()))
	}

	race, err := s.repo.GetRaceByID(id)
	if err != nil {
		return nil, failure.NewError(failure.ErrorInternalServer, fmt.Errorf("failed to get race details by ID: %w", err))
	}
	return race, nil
}

func (s *raceServiceImpl) RegisterRace(race *models.Race) error {
	if err := validateRace(race); err != nil {
		return failure.NewError(failure.ErrorBadRequest, fmt.Errorf("invalid race data: %w", err))
	}

	existingRace, _ := s.repo.GetRaceByName(race.Name)
	if existingRace != nil {
		return failure.NewError(failure.ErrorBadRequest, fmt.Errorf("race with name '%s' already exists", race.Name))
	}

	if err := s.repo.CreateRace(race); err != nil {
		return failure.NewError(failure.ErrorInternalServer, fmt.Errorf("failed to register race: %w", err))
	}
	return nil
}

func (s *raceServiceImpl) UpdateRaceInfo(id uuid.UUID, race *models.Race) error {
	if id == uuid.Nil {
		return failure.NewError(failure.ErrorBadRequest, fmt.Errorf("invalid race ID: %s", id.String()))
	}

	if err := validateRace(race); err != nil {
		return failure.NewError(failure.ErrorBadRequest, fmt.Errorf("invalid race data: %w", err))
	}

	existingRace, err := s.repo.GetRaceByID(id)
	if err != nil {
		return failure.NewError(failure.ErrorBadRequest, fmt.Errorf("race with ID %d does not exist: %w", id, err))
	}

	if existingRace.Name != race.Name {
		duplicateRace, _ := s.repo.GetRaceByName(race.Name)
		if duplicateRace != nil {
			return failure.NewError(failure.ErrorBadRequest, fmt.Errorf("race with name '%s' already exists", race.Name))
		}
	}

	if err := s.repo.UpdateRace(id, race); err != nil {
		return failure.NewError(failure.ErrorInternalServer, fmt.Errorf("failed to update race info: %w", err))
	}
	return nil
}

func (s *raceServiceImpl) RemoveRace(id uuid.UUID) error {
	if id == uuid.Nil {
		return failure.NewError(failure.ErrorBadRequest, fmt.Errorf("invalid race ID: %s", id.String()))
	}

	_, err := s.repo.GetRaceByID(id)
	if err != nil {
		return failure.NewError(failure.ErrorBadRequest, fmt.Errorf("race with ID %d does not exist: %w", id, err))
	}

	if err := s.repo.DeleteRace(id); err != nil {
		return failure.NewError(failure.ErrorInternalServer, fmt.Errorf("failed to remove race: %w", err))
	}
	return nil
}

func (s *raceServiceImpl) AddSubraceToRace(raceID uuid.UUID, subrace *models.Subrace) error {
	if raceID == uuid.Nil {
		return failure.NewError(failure.ErrorBadRequest, fmt.Errorf("invalid race ID: %s", raceID.String()))
	}
	if subrace == nil {
		return failure.NewError(failure.ErrorBadRequest, errors.New("subrace cannot be nil"))
	}

	if err := validateSubrace(subrace); err != nil {
		return failure.NewError(failure.ErrorBadRequest, fmt.Errorf("invalid subrace data: %w", err))
	}

	if err := s.repo.AddSubrace(raceID, subrace); err != nil {
		return failure.NewError(failure.ErrorInternalServer, fmt.Errorf("failed to add subrace to race: %w", err))
	}
	return nil
}

func (s *raceServiceImpl) DetachSubraceFromRace(raceID uuid.UUID, subraceID uuid.UUID) error {
	if raceID == uuid.Nil || subraceID == uuid.Nil {
		return failure.NewError(failure.ErrorBadRequest, fmt.Errorf("invalid race ID or subrace ID: raceID=%d, subraceID=%d", raceID.String(), subraceID.String()))
	}

	if err := s.repo.RemoveSubrace(raceID, subraceID); err != nil {
		return failure.NewError(failure.ErrorInternalServer, fmt.Errorf("failed to detach subrace from race: %w", err))
	}
	return nil
}

func (s *raceServiceImpl) AssignTraitToRace(raceID uuid.UUID, traitID uuid.UUID) error {
	if raceID == uuid.Nil || traitID == uuid.Nil {
		return failure.NewError(failure.ErrorBadRequest, fmt.Errorf("invalid race ID or trait ID: raceID=%d, traitID=%d", raceID.String(), traitID.String()))
	}

	if err := s.repo.AddTrait(raceID, traitID); err != nil {
		return failure.NewError(failure.ErrorInternalServer, fmt.Errorf("failed to assign trait to race: %w", err))
	}
	return nil
}

func (s *raceServiceImpl) UnassignTraitFromRace(raceID uuid.UUID, traitID uuid.UUID) error {
	if raceID == uuid.Nil || traitID == uuid.Nil {
		return failure.NewError(failure.ErrorBadRequest, fmt.Errorf("invalid race ID or trait ID: raceID=%d, traitID=%d", raceID.String(), traitID.String()))
	}

	if err := s.repo.RemoveTrait(raceID, traitID); err != nil {
		return failure.NewError(failure.ErrorInternalServer, fmt.Errorf("failed to unassign trait from race: %w", err))
	}
	return nil
}

func (s *raceServiceImpl) FindRaces(criteria map[string]string) ([]models.Race, error) {
	if len(criteria) == 0 {
		return nil, failure.NewError(failure.ErrorBadRequest, errors.New("no search criteria provided"))
	}

	races, err := s.repo.SearchRaces(criteria)
	if err != nil {
		return nil, failure.NewError(failure.ErrorInternalServer, fmt.Errorf("failed to find races: %w", err))
	}
	return races, nil
}

func validateRace(race *models.Race) error {
	if race.Name == "" {
		return failure.NewError(failure.ErrorBadRequest, errors.New("race name cannot be empty"))
	}
	validSizes := map[string]bool{
		"Small":  true,
		"Medium": true,
		"Large":  true,
	}
	if !validSizes[race.Size] {
		return failure.NewError(failure.ErrorBadRequest, fmt.Errorf("invalid size: %s", race.Size))
	}
	if race.Speed <= 0 {
		return failure.NewError(failure.ErrorBadRequest, fmt.Errorf("invalid speed: %d", race.Speed))
	}

	if err := validateAbilityBonus(&race.AbilityScoreBonuses); err != nil {
		return failure.NewError(failure.ErrorBadRequest, err)
	}

	if err := validateAge(&race.Age); err != nil {
		return failure.NewError(failure.ErrorBadRequest, err)
	}

	for _, prof := range race.Proficiencies {
		if err := validateProficiency(&prof); err != nil {
			return failure.NewError(failure.ErrorBadRequest, err)
		}
	}

	for _, lang := range race.LanguagesKnown {
		if err := validateLanguage(&lang); err != nil {
			return failure.NewError(failure.ErrorBadRequest, err)
		}
	}

	for _, trait := range race.Traits {
		if err := validateTrait(&trait); err != nil {
			return failure.NewError(failure.ErrorBadRequest, err)
		}
	}

	for _, subrace := range race.Subraces {
		if err := validateSubrace(&subrace); err != nil {
			return failure.NewError(failure.ErrorBadRequest, err)
		}
	}

	return nil
}

func validateAbilityBonus(bonus *models.AbilityScoreBonuses) error {
	var errMsgs []string

	if bonus.Strength < 0 {
		errMsgs = append(errMsgs, "Strength bonus cannot be negative")
	}
	if bonus.Dexterity < 0 {
		errMsgs = append(errMsgs, "Dexterity bonus cannot be negative")
	}
	if bonus.Constitution < 0 {
		errMsgs = append(errMsgs, "Constitution bonus cannot be negative")
	}
	if bonus.Intelligence < 0 {
		errMsgs = append(errMsgs, "Intelligence bonus cannot be negative")
	}
	if bonus.Wisdom < 0 {
		errMsgs = append(errMsgs, "Wisdom bonus cannot be negative")
	}
	if bonus.Charisma < 0 {
		errMsgs = append(errMsgs, "Charisma bonus cannot be negative")
	}

	if len(errMsgs) > 0 {
		return fmt.Errorf("negative ability bonuses: %v", errMsgs)
	}

	return nil
}

func validateAge(age *models.Age) error {
	if age.AverageLifespan == "" {
		return errors.New("average lifespan cannot be empty")
	}
	if age.MinimumAge < 0 {
		return fmt.Errorf("minimum age cannot be negative: %d", age.MinimumAge)
	}
	if age.MaximumAge < age.MinimumAge {
		return fmt.Errorf("maximum age (%d) cannot be less than minimum age (%d)", age.MaximumAge, age.MinimumAge)
	}
	return nil
}

func validateProficiency(prof *models.Proficiency) error {
	if prof.Name == "" {
		return errors.New("proficiency name cannot be empty")
	}
	return nil
}

func validateLanguage(lang *models.Language) error {
	if lang.Name == "" {
		return errors.New("language name cannot be empty")
	}
	return nil
}

func validateTrait(trait *models.Trait) error {
	if trait.Name == "" {
		return errors.New("trait name cannot be empty")
	}
	return nil
}

func validateSubrace(subrace *models.Subrace) error {
	if subrace.Name == "" {
		return errors.New("subrace name cannot be empty")
	}
	return nil
}
