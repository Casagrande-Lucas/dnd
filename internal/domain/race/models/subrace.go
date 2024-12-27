package models

import "github.com/google/uuid"

type Subrace struct {
	ID                  uuid.UUID           `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey" swaggertype:"string" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	RaceID              uuid.UUID           `json:"race_id" gorm:"type:uuid;foreignKey:RaceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" swaggertype:"string" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name                string              `json:"name" gorm:"not null"`
	Description         string              `json:"description"`
	AbilityScoreBonuses AbilityScoreBonuses `json:"ability_score_bonuses" gorm:"embedded"`
}
