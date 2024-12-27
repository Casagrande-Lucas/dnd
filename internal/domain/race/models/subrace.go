package models

type Subrace struct {
	ID                  uint                `json:"id" gorm:"primaryKey"`
	RaceID              uint                `json:"race_id"`
	Name                string              `json:"name" gorm:"not null"`
	Description         string              `json:"description"`
	AbilityScoreBonuses AbilityScoreBonuses `json:"ability_score_bonuses" gorm:"embedded"`
}
