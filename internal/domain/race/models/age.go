package models

import "github.com/google/uuid"

type Age struct {
	RaceID          uuid.UUID `json:"race_id" gorm:"type:uuid;primaryKey" swaggertype:"string" format:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	AverageLifespan string    `json:"average_lifespan"`
	MinimumAge      int       `json:"minimum_age"`
	MaximumAge      int       `json:"maximum_age"`
}
