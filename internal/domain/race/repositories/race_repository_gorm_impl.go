package repositories

import (
	"errors"
	"fmt"

	"github.com/Casagrande-Lucas/dnd/internal/domain/race/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// raceRepositoryGormImpl is a concrete implementation of the RaceRepository interface using GORM.
type raceRepositoryGormImpl struct {
	db *gorm.DB
}

// NewGormRaceRepository creates a new instance of raceRepositoryGormImpl.
func NewGormRaceRepository(db *gorm.DB) RaceRepository {
	return &raceRepositoryGormImpl{
		db: db,
	}
}

// GetAllRaces retrieves all races from the database, including their related entities.
func (r *raceRepositoryGormImpl) GetAllRaces() ([]*models.Race, error) {
	var races []*models.Race
	if err := r.db.Preload("Proficiencies").
		Preload("LanguagesKnown").
		Preload("Traits").
		Preload("Subraces").
		Preload("Age").
		Find(&races).Error; err != nil {
		return nil, err
	}
	return races, nil
}

// GetRaceByID retrieves a race by its ID, including its related entities.
func (r *raceRepositoryGormImpl) GetRaceByID(id uuid.UUID) (*models.Race, error) {
	var race models.Race
	if err := r.db.Preload("Proficiencies").
		Preload("LanguagesKnown").
		Preload("Traits").
		Preload("Subraces").
		Preload("Age").
		First(&race, "id = ?", id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("race with ID %s not found", id.String())
		}
		return nil, err
	}
	return &race, nil
}

// GetRaceByName retrieves a race by its name, including its related entities.
func (r *raceRepositoryGormImpl) GetRaceByName(name string) (*models.Race, error) {
	var race models.Race
	if err := r.db.Preload("Proficiencies").
		Preload("LanguagesKnown").
		Preload("Traits").
		Preload("Subraces").
		Preload("Age").
		Where("name = ?", name).
		First(&race).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("race with name '%s' not found", name)
		}
		return nil, err
	}
	return &race, nil
}

// CreateRace adds a new race to the database along with its related entities.
func (r *raceRepositoryGormImpl) CreateRace(race *models.Race) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(race).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// UpdateRace updates an existing race's details in the database.
func (r *raceRepositoryGormImpl) UpdateRace(id uuid.UUID, race *models.Race) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingRace models.Race
	if err := tx.Preload("Proficiencies").
		Preload("LanguagesKnown").
		Preload("Traits").
		Preload("Subraces").
		Preload("Age").
		First(&existingRace, "id = ?", id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("race with ID %s not found", id.String())
		}
		return err
	}

	existingRace.Name = race.Name
	existingRace.Description = race.Description
	existingRace.Size = race.Size
	existingRace.Speed = race.Speed
	existingRace.Alignment = race.Alignment
	existingRace.AbilityScoreBonuses = race.AbilityScoreBonuses

	if err := tx.Model(&existingRace).Association("Age").Replace(&race.Age); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&existingRace).Association("Proficiencies").Replace(race.Proficiencies); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&existingRace).Association("LanguagesKnown").Replace(race.LanguagesKnown); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&existingRace).Association("Traits").Replace(race.Traits); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&existingRace).Association("Subraces").Replace(race.Subraces); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(&existingRace).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// DeleteRace removes a race from the database.
func (r *raceRepositoryGormImpl) DeleteRace(id uuid.UUID) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var race models.Race
	if err := tx.First(&race, "id = ?", id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("race with ID %d not found", id)
		}
		return err
	}

	if err := tx.Delete(&race).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// AddSubrace adds a subrace to a specific race.
func (r *raceRepositoryGormImpl) AddSubrace(raceID uuid.UUID, subrace *models.Subrace) error {
	var race models.Race
	if err := r.db.First(&race, raceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("race with ID %d not found", raceID)
		}
		return err
	}

	subrace.RaceID = raceID
	if err := r.db.Create(subrace).Error; err != nil {
		return err
	}

	return nil
}

// RemoveSubrace removes a subrace from a specific race.
func (r *raceRepositoryGormImpl) RemoveSubrace(raceID uuid.UUID, subraceID uuid.UUID) error {
	var subrace models.Subrace
	if err := r.db.Where("id = ? AND race_id = ?", subraceID, raceID).First(&subrace).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("subrace with ID %d not found for race ID %d", subraceID, raceID)
		}
		return err
	}

	if err := r.db.Delete(&subrace).Error; err != nil {
		return err
	}

	return nil
}

// AddTrait associates a trait with a specific race.
func (r *raceRepositoryGormImpl) AddTrait(raceID uuid.UUID, traitID uuid.UUID) error {
	var race models.Race
	if err := r.db.Preload("Traits").First(&race, raceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("race with ID %d not found", raceID)
		}
		return err
	}

	var trait models.Trait
	if err := r.db.First(&trait, traitID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("trait with ID %d not found", traitID)
		}
		return err
	}

	return r.db.Model(&race).Association("Traits").Append(&trait)
}

// RemoveTrait dissociates a trait from a specific race.
func (r *raceRepositoryGormImpl) RemoveTrait(raceID uuid.UUID, traitID uuid.UUID) error {
	var race models.Race
	if err := r.db.Preload("Traits").First(&race, raceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("race with ID %d not found", raceID)
		}
		return err
	}

	var trait models.Trait
	if err := r.db.First(&trait, traitID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("trait with ID %d not found", traitID)
		}
		return err
	}

	return r.db.Model(&race).Association("Traits").Delete(&trait)
}

// SearchRaces allows searching for races based on specific criteria.
func (r *raceRepositoryGormImpl) SearchRaces(criteria map[string]string) ([]models.Race, error) {
	var races []models.Race
	query := r.db.Preload("Proficiencies").
		Preload("LanguagesKnown").
		Preload("Traits").
		Preload("Subraces").
		Preload("Age")

	for key, value := range criteria {
		switch key {
		case "size":
			query = query.Where("size = ?", value)
		case "speed":
			query = query.Where("speed = ?", value)
		case "alignment":
			query = query.Where("alignment = ?", value)
		default:
			return nil, fmt.Errorf("unknown search criteria: %s", key)
		}
	}

	if err := query.Find(&races).Error; err != nil {
		return nil, err
	}

	return races, nil
}
