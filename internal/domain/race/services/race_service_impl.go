package services

import (
	"errors"
	"fmt"

	"github.com/Casagrande-Lucas/dnd/internal/domain/race/models"
	"github.com/Casagrande-Lucas/dnd/internal/domain/race/repositories"
	"github.com/Casagrande-Lucas/dnd/pkg/failure"
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
		return nil, failure.NewError(failure.ErrorInternalServer, err)
	}
	return races, nil
}

func (s *raceServiceImpl) GetRaceDetails(id uint) (*models.Race, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid race ID: %d", id)
	}

	race, err := s.repo.GetRaceByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get race details by ID: %w", err)
	}
	return race, nil
}

func (s *raceServiceImpl) RegisterRace(race *models.Race) error {
	if err := validateRace(race); err != nil {
		return fmt.Errorf("invalid race data: %w", err)
	}

	existingRace, _ := s.repo.GetRaceByName(race.Name)
	if existingRace != nil {
		return fmt.Errorf("race with name '%s' already exists", race.Name)
	}

	if err := s.repo.CreateRace(race); err != nil {
		return fmt.Errorf("failed to register race: %w", err)
	}
	return nil
}

func (s *raceServiceImpl) UpdateRaceInfo(id uint, race *models.Race) error {
	if id == 0 {
		return fmt.Errorf("invalid race ID: %d", id)
	}

	if err := validateRace(race); err != nil {
		return fmt.Errorf("invalid race data: %w", err)
	}

	existingRace, err := s.repo.GetRaceByID(id)
	if err != nil {
		return fmt.Errorf("race with ID %d does not exist: %w", id, err)
	}

	if existingRace.Name != race.Name {
		duplicateRace, _ := s.repo.GetRaceByName(race.Name)
		if duplicateRace != nil {
			return fmt.Errorf("race with name '%s' already exists", race.Name)
		}
	}

	if err := s.repo.UpdateRace(id, race); err != nil {
		return fmt.Errorf("failed to update race info: %w", err)
	}
	return nil
}

func (s *raceServiceImpl) RemoveRace(id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid race ID: %d", id)
	}

	_, err := s.repo.GetRaceByID(id)
	if err != nil {
		return fmt.Errorf("race with ID %d does not exist: %w", id, err)
	}

	if err := s.repo.DeleteRace(id); err != nil {
		return fmt.Errorf("failed to remove race: %w", err)
	}
	return nil
}

func (s *raceServiceImpl) AddSubraceToRace(raceID uint, subrace *models.Subrace) error {
	if raceID == 0 {
		return fmt.Errorf("invalid race ID: %d", raceID)
	}
	if subrace == nil {
		return errors.New("subrace cannot be nil")
	}

	if err := validateSubrace(subrace); err != nil {
		return fmt.Errorf("invalid subrace data: %w", err)
	}

	if err := s.repo.AddSubrace(raceID, subrace); err != nil {
		return fmt.Errorf("failed to add subrace to race: %w", err)
	}
	return nil
}

func (s *raceServiceImpl) DetachSubraceFromRace(raceID uint, subraceID uint) error {
	if raceID == 0 || subraceID == 0 {
		return fmt.Errorf("invalid race ID or subrace ID: raceID=%d, subraceID=%d", raceID, subraceID)
	}

	if err := s.repo.RemoveSubrace(raceID, subraceID); err != nil {
		return fmt.Errorf("failed to detach subrace from race: %w", err)
	}
	return nil
}

func (s *raceServiceImpl) AssignTraitToRace(raceID uint, traitID uint) error {
	if raceID == 0 || traitID == 0 {
		return fmt.Errorf("invalid race ID or trait ID: raceID=%d, traitID=%d", raceID, traitID)
	}

	if err := s.repo.AddTrait(raceID, traitID); err != nil {
		return fmt.Errorf("failed to assign trait to race: %w", err)
	}
	return nil
}

func (s *raceServiceImpl) UnassignTraitFromRace(raceID uint, traitID uint) error {
	if raceID == 0 || traitID == 0 {
		return fmt.Errorf("invalid race ID or trait ID: raceID=%d, traitID=%d", raceID, traitID)
	}

	if err := s.repo.RemoveTrait(raceID, traitID); err != nil {
		return fmt.Errorf("failed to unassign trait from race: %w", err)
	}
	return nil
}

func (s *raceServiceImpl) FindRaces(criteria map[string]string) ([]models.Race, error) {
	if len(criteria) == 0 {
		return nil, errors.New("no search criteria provided")
	}

	races, err := s.repo.SearchRaces(criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to find races: %w", err)
	}
	return races, nil
}

func validateRace(race *models.Race) error {
	if race.Name == "" {
		return errors.New("race name cannot be empty")
	}
	validSizes := map[string]bool{
		"Small":  true,
		"Medium": true,
		"Large":  true,
	}
	if !validSizes[race.Size] {
		return fmt.Errorf("invalid size: %s", race.Size)
	}
	if race.Speed <= 0 {
		return fmt.Errorf("invalid speed: %d", race.Speed)
	}

	if err := validateAbilityBonus(&race.AbilityScoreBonuses); err != nil {
		return err
	}

	if err := validateAge(&race.Age); err != nil {
		return err
	}

	for _, prof := range race.Proficiencies {
		if err := validateProficiency(&prof); err != nil {
			return err
		}
	}

	for _, lang := range race.LanguagesKnown {
		if err := validateLanguage(&lang); err != nil {
			return err
		}
	}

	for _, trait := range race.Traits {
		if err := validateTrait(&trait); err != nil {
			return err
		}
	}

	for _, subrace := range race.Subraces {
		if err := validateSubrace(&subrace); err != nil {
			return err
		}
	}

	for _, lang := range race.Languages {
		if err := validateLanguage(&lang); err != nil {
			return err
		}
	}

	return nil
}

func validateAbilityBonus(bonus *models.AbilityBonus) error {
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
